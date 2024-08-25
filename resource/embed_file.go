package resource

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/util/logger"
)

var webTitle string

// 嵌入静态文件
func EmbedResoure(engine *gin.Engine, title string) {
	webTitle = title
	//go template 设置网页标题
	embedTemplate(engine)
	//vue静态文件 web阅读器主体
	EmbedWeb(engine)
	//react静态文件，admin界面
	EmbedAdmin(engine)
}

// 使用go模板，设置网页标题
//
//go:embed web_static/index.html
var templateFile string

func embedTemplate(engine *gin.Engine) {
	//使用自定义的模板引擎，命名为"template-data"，为了与VUE兼容，把左右分隔符自定义为 [[ ]]
	tmpl := template.Must(template.New("template-data").Delims("[[", "]]").Parse(templateFile))
	//使用模板
	engine.SetHTMLTemplate(tmpl)
	//解析模板到HTML
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template-data", gin.H{
			"title": webTitle, //页面标题
		})
	})
}

//go:embed  web_static
var staticFS embed.FS

//go:embed  web_static/assets
var staticAssetFS embed.FS

//go:embed  web_static/images
var staticImageFS embed.FS

func EmbedWeb(engine *gin.Engine) {
	//https://stackoverflow.com/questions/66248258/serve-embedded-filesystem-from-root-path-of-url
	assetsEmbedFS, err := fs.Sub(staticAssetFS, "web_static/assets")
	if err != nil {
		logger.Infof("%s", err)
	}
	engine.StaticFS("/assets/", http.FS(assetsEmbedFS))
	imagesEmbedFS, errStaticImageFS := fs.Sub(staticImageFS, "web_static/images")
	if errStaticImageFS != nil {
		logger.Info(errStaticImageFS)
	}
	engine.StaticFS("/images/", http.FS(imagesEmbedFS))
	//favicon.ico
	engine.GET("/favicon.ico", func(c *gin.Context) {
		file, _ := staticFS.ReadFile("web_static/images/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
}

//go:embed  admin_static
var adminFS embed.FS

func EmbedAdmin(engine *gin.Engine) {
	adminEmbedFS, errAdminFS := fs.Sub(adminFS, "admin_static")
	if errAdminFS != nil {
		logger.Info(errAdminFS)
	}
	engine.StaticFS("/admin", http.FS(adminEmbedFS))
}
