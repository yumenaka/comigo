package router

import (
	"context"
	"embed"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/htmx/file_server"
	"github.com/yumenaka/comi/htmx/router/handler"
)

//go:embed all:static
var static embed.FS

// TemplRender 实现了 render.Render 接口。
type TemplRender struct {
	Code int
	Data templ.Component
}

// Render 实现了 render.Render 接口。
func (t TemplRender) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}
	return nil
}

// WriteContentType 实现了 render.Render 接口。
func (t TemplRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// Instance 实现了render.Render接口。
func (t *TemplRender) Instance(name string, data interface{}) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &TemplRender{
			Code: http.StatusOK,
			Data: templData,
		}
	}
	return nil
}

// RunServer 运行一个新的 HTTP 服务器。
func RunServer() (err error) {
	// 扫描漫画
	file_server.StartComigoServer()
	// 创建一个新的Gin服务器。
	router := gin.Default()
	// 为模板引擎定义 HTML 渲染器。
	router.HTMLRender = &TemplRender{}
	// 静态文件。
	router.Static("/static", "./router/static")
	// 嵌入静态文件。
	// staticFS, err := fs.Sub(static, "static")
	// if err != nil {
	// 	logger.Infof("%s", err)
	// }
	// router.StaticFS("/static/", http.FS(staticFS))

	// 处理主页视图
	router.GET("/", handler.IndexViewHandler)
	// 处理 API
	router.GET("/api/hello-world", handler.ShowContentAPIHandler)

	// 发消息
	slog.Info("Starting server...", "port", config.Config.Port)

	// 是否对外服务
	webHost := ":"
	if config.Config.DisableLAN {
		webHost = "localhost:"
	}
	// 是否启用TLS
	enableTls := config.Config.CertFile != "" && config.Config.KeyFile != ""
	server := &http.Server{
		Addr:         webHost + strconv.Itoa(config.Config.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router, // gin.Engine本身可以作为一个Handler传递到http包,用于启动服务器
	}

	// 监听并启动服务(TLS)
	if enableTls {
		if err := server.ListenAndServeTLS(config.Config.CertFile, config.Config.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			log.Fatalf("listen: %s\n", err)
		}
	}
	if !enableTls {
		// 监听并启动服务(HTTP)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			time.Sleep(3 * time.Second)
			log.Fatalf("listen: %s\n", err)
		}
	}
	return err
}
