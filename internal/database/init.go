//go:build !(windows && 386)

package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/yumenaka/comi/util/locale"
	"github.com/yumenaka/comi/util/logger"
	"path"
	"path/filepath"

	"entgo.io/ent/dialect"
	"github.com/yumenaka/comi/internal/ent"
	"modernc.org/sqlite"
)

// 参考：
// Go製CGOフリーなSQLiteドライバーでentを使う
// https://zenn.dev/nobonobo/articles/e9f17d183c19f6

// 数据库为sqlite3
// 查看工具：SQLiteStudio https://github.com/pawelsalawa/sqlitestudio/releases
// 查看工具： DB Browser for SQLite  https://sqlitebrowser.org/dl/
type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to enable enable foreign keys: %w", err)
	}
	return conn, nil
}

// 注册 sqlite
func init() {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
}

var client *ent.Client

func InitDatabase(configFilePath string) error {
	if client != nil {
		//logger.Infof("database already initialized")
		return nil
	}
	//链接或创建数据库
	var entOptions []ent.Option
	//是否打印log
	//entOptions = append(entOptions, ent.Debug())
	//连接器
	var err error
	dataSourceName := "file:comigo.sqlite?cache=shared"
	//如果有配置文件的话，数据库文件，就在同一文件夹内
	if configFilePath != "" {
		configDir := filepath.Dir(configFilePath) //不能用path.Dir()，因为windows返回 "."
		dataSourceName = "file:" + path.Join(configDir, "comigo.sqlite") + "?cache=shared"
	}
	logger.Infof(locale.GetString("InitDatabase")+"%s", dataSourceName)
	client, err = ent.Open(dialect.SQLite, dataSourceName, entOptions...)
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
		//time.Sleep(3 * time.Second)
		//log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	//defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
		//time.Sleep(3 * time.Second)
		//log.Fatalf("failed creating schema resources: %v", err)
	}
	return nil
}

func CloseDatabase() {
	err := client.Close()
	if err != nil {
		logger.Infof("%s", err)
	}
}
