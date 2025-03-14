package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/cmd/image_viewer/ent"
	_ "modernc.org/sqlite" // SQLite 驱动 (纯Go实现，无需CGO)
)

// 全局配置变量
var (
	rootPath   string          // 要扫描的根目录路径
	useSQLite  bool            // 是否使用SQLite数据库存储（默认使用JSON文件）
	maxDepth   int             // 最大扫描深度（-1表示不限制）
	ignoreDirs map[string]bool // 要忽略的目录名
)

func main() {
	// 解析命令行参数：root 路径和是否使用 SQLite
	rootPathFlag := flag.String("root", ".", "要扫描的根目录路径")
	useSQLiteFlag := flag.Bool("sqlite", false, "是否使用SQLite数据库存储（默认使用JSON文件）")
	scanDepthFlag := flag.Int("maxdepth", -1, "最大扫描深度（-1表示不限制）")
	flag.Parse()

	rootPath = *rootPathFlag
	useSQLite = *useSQLiteFlag
	maxDepth = *scanDepthFlag

	// 将要忽略的目录名转为小写存入map，便于快速判断
	ignoreNames := []string{"temp", ".cache"}
	ignoreDirs = make(map[string]bool)
	for _, name := range ignoreNames {
		ignoreDirs[strings.ToLower(name)] = true
	}

	// 初始化 Echo Web 服务器
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 如果使用 SQLite 存储，初始化数据库客户端
	if useSQLite {
		var err error
		// 打开 SQLite 数据库文件（如果不存在将自动创建）
		// 使用modernc.org/sqlite驱动，需要手动创建SQL连接并注册到ent
		db, err := sql.Open("sqlite", "file:comigo.db?_pragma=foreign_keys(1)")
		if err != nil {
			e.Logger.Fatal("无法打开 SQLite 数据库: ", err)
		}

		// 创建ent驱动
		drv := entsql.OpenDB(dialect.SQLite, db)

		// 创建ent客户端
		entClient = ent.NewClient(ent.Driver(drv))

		// 确保关闭数据库
		defer func() {
			entClient.Close()
			db.Close()
		}()

		// 运行自动迁移，创建表结构
		if err := entClient.Schema.Create(context.Background()); err != nil {
			e.Logger.Fatal("数据库迁移失败: ", err)
		}
	}

	// 如果使用 JSON 模式，尝试加载已有的数据文件，避免重复扫描
	if !useSQLite {
		if err := loadJSONStorage(); err != nil {
			fmt.Println("加载JSON数据失败，将进行初始扫描:", err)
		}
	}

	// 应用启动时先进行一次初始扫描
	go performScan()

	// 设置后台定期扫描任务（例如每隔1小时扫描一次，可根据需要调整）
	scanInterval := time.Hour * 1
	go func() {
		ticker := time.NewTicker(scanInterval)
		defer ticker.Stop()
		for range ticker.C {
			performScan()
		}
	}()

	// 注册路由
	e.GET("/", indexHandler)
	e.GET("/api/list", listHandler)
	e.POST("/api/rescan", rescanHandler)
	e.GET("/raw", rawImageHandler)

	// 启动服务器
	e.Logger.Fatal(e.Start(":1323"))
}
