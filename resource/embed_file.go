package resource

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/util/logger"
)

var webTitle string

// 嵌入静态文件
func EmbedResoure(e *echo.Echo, title string) {
	webTitle = title
	// go template 设置网页标题
	embedTemplate(e)
	// vue静态文件 web阅读器主体
	EmbedWeb(e)
	// react静态文件，admin界面
	EmbedAdmin(e)
}

//go:embed web_static/index.html
var templateFile string

// 自定义模板渲染器
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func embedTemplate(e *echo.Echo) {
	// 创建模板并设置分隔符
	t := &Template{
		templates: template.Must(template.New("template-data").Delims("[[", "]]").Parse(templateFile)),
	}
	// 设置渲染器
	e.Renderer = t

	// 路由处理
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "template-data", map[string]interface{}{
			"title": webTitle,
		})
	})
}

//go:embed  web_static
var staticFS embed.FS

//go:embed  web_static/assets
var staticAssetFS embed.FS

//go:embed  web_static/images
var staticImageFS embed.FS

func EmbedWeb(e *echo.Echo) {
	assetsEmbedFS, err := fs.Sub(staticAssetFS, "web_static/assets")
	if err != nil {
		logger.Infof("%s", err)
	}
	assetHandler := http.FileServer(http.FS(assetsEmbedFS))
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))

	imagesEmbedFS, errStaticImageFS := fs.Sub(staticImageFS, "web_static/images")
	if errStaticImageFS != nil {
		logger.Info(errStaticImageFS)
	}
	imageHandler := http.FileServer(http.FS(imagesEmbedFS))
	e.GET("/images/*", echo.WrapHandler(http.StripPrefix("/images/", imageHandler)))

	// favicon.ico
	e.GET("/favicon.ico", func(c echo.Context) error {
		file, _ := staticFS.ReadFile("web_static/images/favicon.ico")
		return c.Blob(http.StatusOK, "image/x-icon", file)
	})
}

//go:embed  admin_static
var adminFS embed.FS

func EmbedAdmin(e *echo.Echo) {
	adminEmbedFS, errAdminFS := fs.Sub(adminFS, "admin_static")
	if errAdminFS != nil {
		logger.Info(errAdminFS)
	}
	adminHandler := http.FileServer(http.FS(adminEmbedFS))
	e.GET("/admin/*", echo.WrapHandler(http.StripPrefix("/admin/", adminHandler)))
}
