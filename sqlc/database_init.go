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

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/sqlc/postgres"
	"github.com/yumenaka/comigo/tools/logger"
	_ "modernc.org/sqlite"
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
		return openSQLiteDatabase(options.ConfigDir)
	case "postgres":
		return openPostgresDatabase(strings.TrimSpace(options.DSN))
	default:
		return fmt.Errorf("unsupported database type: %s", options.Type)
	}
}

func openSQLiteDatabase(configDir string) (openErr error) {
	defer func() {
		if openErr != nil {
			CloseDatabase()
		}
	}()
	dataSourceName := ":memory:"
	if configDir != "" {
		if err := os.MkdirAll(configDir, 0o700); err != nil {
			return err
		}
		dbPath, err := filepath.Abs(filepath.Join(configDir, "comigo.sqlite"))
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

	// create tables - 现在使用 IF NOT EXISTS，所以即使表已存在也不会报错
	if _, err := client.ExecContext(ctx, ddl); err != nil {
		logger.Infof(locale.GetString("log_failed_to_create_tables"), err)
		return err
	}
	DbStore = NewDBStore(client)
	if err := validateSchema(DbStore); err != nil {
		return err
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

// validateSchema 拒绝不兼容的旧结构，不静默退回 JSON，也不自动删除用户数据。
func validateSchema(db *StoreDatabase) error {
	_, err := db.queries.GetBookByID(context.Background(), "")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
