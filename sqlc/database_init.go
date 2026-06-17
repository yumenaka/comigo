//go:build !js

package sqlc

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"path"
	"strings"

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
	queries bookQueries
}

// DBOptions 描述本次启动要使用的数据库后端。
type DBOptions struct {
	Type      string
	DSN       string
	ConfigDir string
}

// NewDBStore 创建新的BookRepository实例
func NewDBStore(db DBTX) *StoreDatabase {
	return &StoreDatabase{
		queries: New(db),
	}
}

// NewPostgresDBStore 使用 PostgreSQL 生成查询和 adapter 创建统一的数据访问层。
func NewPostgresDBStore(db postgres.DBTX) *StoreDatabase {
	return &StoreDatabase{
		queries: newPostgresAdapter(postgres.New(db)),
	}
}

// OpenDatabase 根据配置选择 SQLite 或 PostgreSQL 后端。
func OpenDatabase(options DBOptions) error {
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

func openSQLiteDatabase(configDir string) error {
	// 文件类型数据库，默认在当前目录下创建 comigo.sqlite 文件
	// 内存数据库的语法是:  dataSourceName := ":memory:"
	dataSourceName := "file:comigo.sqlite?cache=shared"
	// 把数据库文件在configDir文件夹内
	if configDir != "" {
		dataSourceName = "file:" + path.Join(configDir, "comigo.sqlite") + "?cache=shared"
		logger.Infof(locale.GetString("init_database")+"%s", dataSourceName)
	}
	if configDir == "" {
		dataSourceName = ":memory:"
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
		// 即使创建表失败，我们也要尝试创建 DBQueries，因为表可能已经存在
		// 只要数据库连接正常，就应该能正常工作
	}
	if err := migrateDatabase(ctx, client); err != nil {
		logger.Infof("database migration failed: %v", err)
		return err
	}

	// 创建 StoreDatabase 实例
	DbStore = NewDBStore(client)
	logger.Info(locale.GetString("log_database_initialized_successfully"))
	return nil
}

// openPostgresDatabase 初始化 PostgreSQL 连接、schema 和迁移补丁。
func openPostgresDatabase(dsn string) error {
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
	}
	if err := migratePostgresDatabase(ctx, client); err != nil {
		logger.Infof("database migration failed: %v", err)
		return err
	}
	DbStore = NewPostgresDBStore(client)
	logger.Info(locale.GetString("log_database_initialized_successfully"))
	return nil
}

func configureSQLitePragmas(ctx context.Context, db *sql.DB) error {
	// 建库时开启 incremental auto_vacuum，后续大量删除书籍数据后可以逐步回收空闲页。
	if _, err := db.ExecContext(ctx, "PRAGMA auto_vacuum = INCREMENTAL"); err != nil {
		return err
	}
	return nil
}

func CloseDatabase() {
	if client == nil {
		return
	}
	err := client.Close()
	if err != nil {
		logger.Infof("%s", err)
	}
}

func migrateDatabase(ctx context.Context, db *sql.DB) error {
	if err := ensureColumn(ctx, db, "bookmarks", "book_store_id", "TEXT"); err != nil {
		return err
	}
	return nil
}

func migratePostgresDatabase(ctx context.Context, db *sql.DB) error {
	if err := ensurePostgresColumn(ctx, db, "bookmarks", "book_store_id", "TEXT"); err != nil {
		return err
	}
	return nil
}

func ensureColumn(ctx context.Context, db *sql.DB, tableName string, columnName string, columnType string) error {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info("+tableName+")")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == columnName {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "ALTER TABLE "+tableName+" ADD COLUMN "+columnName+" "+columnType)
	return err
}

func ensurePostgresColumn(ctx context.Context, db *sql.DB, tableName string, columnName string, columnType string) error {
	_, err := db.ExecContext(ctx, "ALTER TABLE "+tableName+" ADD COLUMN IF NOT EXISTS "+columnName+" "+columnType)
	return err
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
