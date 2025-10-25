//go:build !js

package sqlc

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"path"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	_ "modernc.org/sqlite"
)

// 参考：
// https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html#setting-up

//go:embed schema.sql
var ddl string

var (
	client  *sql.DB
	DbStore *StoreDatabase
)

// StoreDatabase 书籍数据访问层
type StoreDatabase struct {
	queries *Queries
}

// NewDBStore 创建新的BookRepository实例
func NewDBStore(db DBTX) *StoreDatabase {
	return &StoreDatabase{
		queries: New(db),
	}
}

func OpenDatabase(configDir string) error {
	// 文件类型数据库，默认在当前目录下创建 comigo.sqlite 文件
	// 内存数据库的语法是:  dataSourceName := ":memory:"
	dataSourceName := "file:comigo.sqlite?cache=shared"
	// 把数据库文件在configDir文件夹内
	if configDir != "" {
		dataSourceName = "file:" + path.Join(configDir, "comigo.sqlite") + "?cache=shared"
		logger.Infof(locale.GetString("init_database")+"%s", dataSourceName)
	}
	if configDir == "" {
		dataSourceName = "\":memory:\""
	}
	ctx := context.Background()
	var err error
	client, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		logger.Infof("Failed to open database: %v", err)
		return err
	}

	// Test database connection
	if err = client.PingContext(ctx); err != nil {
		logger.Infof("Failed to ping database: %v", err)
		return err
	}

	// create tables - 现在使用 IF NOT EXISTS，所以即使表已存在也不会报错
	if _, err := client.ExecContext(ctx, ddl); err != nil {
		logger.Infof("Failed to create tables: %v", err)
		// 即使创建表失败，我们也要尝试创建 DBQueries，因为表可能已经存在
		// 只要数据库连接正常，就应该能正常工作
	}

	// 创建 StoreDatabase 实例
	DbStore = NewDBStore(client)
	logger.Infof("Database initialized successfully")
	return nil
}

func CloseDatabase() {
	err := client.Close()
	if err != nil {
		logger.Infof("%s", err)
	}
}

// CheckDBQueries 检查 queries 是否已初始化
func (db *StoreDatabase) CheckDBQueries() error {
	if db.queries == nil {
		return fmt.Errorf("database not initialized, DBQueries is nil")
	}
	return nil
}
