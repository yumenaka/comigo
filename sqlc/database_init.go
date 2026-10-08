//go:build !js

package sqlc

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/sqlc/postgres"
	"github.com/yumenaka/comigo/tools/logger"
	"modernc.org/sqlite"
)

// 参考：
// https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html#setting-up

//go:embed schema.sql
var ddl string

//go:embed postgres/schema.sql
var postgresDDL string

var (
	client  *sql.DB
	DbStore *StoreDatabase
)

// StoreDatabase 书籍数据访问层
type StoreDatabase struct {
	queries    bookQueries
	connection *sql.DB
	postgres   bool
	mu         sync.RWMutex
	groupsMu   sync.Mutex // 串行重建拓扑；计算期间仍允许读取旧书组。
}

// DBOptions 描述本次启动要使用的数据库后端。
type DBOptions struct {
	Type      string
	DSN       string
	ConfigDir string
}

// NewDBStore 创建新的BookRepository实例
func NewDBStore(db *sql.DB) *StoreDatabase {
	return &StoreDatabase{
		queries:    New(db),
		connection: db,
	}
}

// NewPostgresDBStore 使用 PostgreSQL 生成查询和 adapter 创建统一的数据访问层。
func NewPostgresDBStore(db *sql.DB) *StoreDatabase {
	return &StoreDatabase{
		queries:    newPostgresAdapter(postgres.New(db)),
		connection: db,
		postgres:   true,
	}
}

// OpenDatabase 根据配置选择 SQLite 或 PostgreSQL 后端。
func OpenDatabase(options DBOptions) error {
	CloseDatabase()
	dbType := strings.ToLower(strings.TrimSpace(options.Type))
	switch dbType {
	case "sqlite":
		return openSQLiteDatabase(options.ConfigDir, true)
	case "postgres":
		return openPostgresDatabase(strings.TrimSpace(options.DSN))
	default:
		return fmt.Errorf("unsupported database type: %s", options.Type)
	}
}

// openSQLiteDatabase 初始化 SQLite；旧结构备份后只重试一次，避免循环重建。
func openSQLiteDatabase(configDir string, rebuild bool) (openErr error) {
	defer func() {
		if openErr != nil {
			CloseDatabase()
		}
	}()
	dataSourceName := ":memory:"
	var dbPath string
	if configDir != "" {
		if err := os.MkdirAll(configDir, 0o700); err != nil {
			return err
		}
		var err error
		dbPath, err = filepath.Abs(filepath.Join(configDir, "comigo.sqlite"))
		if err != nil {
			return err
		}
		dataSourceName = (&url.URL{Scheme: "file", Path: dbPath}).String()
	}

	ctx := context.Background()
	var err error
	client, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_database"), err)
		return err
	}

	// Test database connection
	if err = client.PingContext(ctx); err != nil {
		logger.Infof(locale.GetString("log_failed_to_ping_database"), err)
		return err
	}
	if err := configureSQLitePragmas(ctx, client); err != nil {
		logger.Infof("database pragma configuration failed: %v", err)
		return err
	}

	// 建表或查询缺少列时说明旧结构不兼容；其他错误不能触发重建。
	_, err = client.ExecContext(ctx, ddl)
	if err == nil {
		DbStore = NewDBStore(client)
		err = validateSchema(DbStore)
	}
	if err != nil {
		var sqliteErr *sqlite.Error
		if !rebuild || dbPath == "" || !errors.As(err, &sqliteErr) || sqliteErr.Code() != 1 ||
			(!strings.Contains(err.Error(), "no such column:") && !strings.Contains(err.Error(), "no such table:")) {
			return err
		}
		// 先写回 WAL；仍有事务占用时停止，不能只移动主文件而丢失已提交数据。
		if _, err := client.ExecContext(ctx, "PRAGMA locking_mode = EXCLUSIVE; BEGIN EXCLUSIVE; COMMIT;"); err != nil {
			return err
		}
		var busy, logFrames, checkpointed int
		if err := client.QueryRowContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)").Scan(&busy, &logFrames, &checkpointed); err != nil {
			return err
		}
		if busy != 0 {
			return fmt.Errorf("database is busy; cannot back up incompatible SQLite database")
		}
		if err := client.Close(); err != nil {
			return err
		}
		client, DbStore = nil, nil
		backupPath := dbPath + ".bak-" + time.Now().Format("20060102-150405.000000000")
		if err := os.Rename(dbPath, backupPath); err != nil {
			return err
		}
		logger.Infof(locale.GetString("log_sqlite_schema_rebuilt"), backupPath)
		return openSQLiteDatabase(configDir, false)
	}

	logger.Info(locale.GetString("log_database_initialized_successfully"))
	return nil
}

// openPostgresDatabase 初始化 PostgreSQL 连接并检查 schema。
func openPostgresDatabase(dsn string) (openErr error) {
	defer func() {
		if openErr != nil {
			CloseDatabase()
		}
	}()
	if dsn == "" {
		return fmt.Errorf("postgres database dsn is empty")
	}
	ctx := context.Background()
	var err error
	client, err = sql.Open("pgx", dsn)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_open_database"), err)
		return err
	}
	if err = client.PingContext(ctx); err != nil {
		logger.Infof(locale.GetString("log_failed_to_ping_database"), err)
		return err
	}
	if _, err := client.ExecContext(ctx, postgresDDL); err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_tables"), err)
		return err
	}
	DbStore = NewPostgresDBStore(client)
	if err := validateSchema(DbStore); err != nil {
		return err
	}

	logger.Info(locale.GetString("log_database_initialized_successfully"))
	return nil
}

func configureSQLitePragmas(ctx context.Context, db *sql.DB) error {
	// SQLite 的连接级设置与内存库必须使用同一个连接；写入由事务串行完成。
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, "PRAGMA auto_vacuum = INCREMENTAL; PRAGMA foreign_keys = ON; PRAGMA busy_timeout = 5000; PRAGMA journal_mode = WAL;"); err != nil {
		return err
	}

	return nil
}

func CloseDatabase() {
	if client == nil {
		return
	}
	err := client.Close()
	client = nil
	DbStore = nil
	if err != nil {
		logger.Infof("%s", err)
	}
}

// validateSchema 检查书籍、页面和书签结构，由对应后端决定如何处理不兼容。
func validateSchema(db *StoreDatabase) error {
	ctx := context.Background()
	_, bookErr := db.queries.GetBookByID(ctx, "")
	if errors.Is(bookErr, sql.ErrNoRows) {
		bookErr = nil
	}
	_, pageErr := db.queries.GetPageInfosByBookID(ctx, "")
	_, bookmarkErr := db.queries.ListBookmarksByBookID(ctx, "")
	err := errors.Join(bookErr, pageErr, bookmarkErr)
	if err != nil {
		return fmt.Errorf("incompatible database schema (use a new database): %w", err)
	}
	return nil
}

// CheckDBQueries 检查 queries 是否已初始化
func (db *StoreDatabase) CheckDBQueries() error {
	if db == nil {
		return fmt.Errorf("database not initialized, StoreDatabase is nil")
	}
	if db.queries == nil {
		return fmt.Errorf("database not initialized, DBQueries is nil")
	}
	return nil
}
