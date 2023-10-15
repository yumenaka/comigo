package routers

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/util"

	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/types"
)

// TemplateString 模板文件
//
//go:embed static/index.html
var TemplateString string

//go:embed  static
var staticFS embed.FS

//go:embed  static/assets
var staticAssetFS embed.FS

//go:embed  static/images
var staticImageFS embed.FS

//go:embed  admin
var adminFS embed.FS

// 用来防止重复注册的URL表，key是bookID，值是StaticURL
var staticUrlMap = make(map[string]string)

func checkUrlRegistered(bookID string) bool {
	_, ok := staticUrlMap[bookID]
	return ok
}

// SetDownloadLink 设定压缩包下载链接
func SetDownloadLink() {
	if types.GetBooksNumber() >= 1 {
		allBook, err := types.GetAllBookInfoList("name")
		if err != nil {
			fmt.Println("设置文件下载失败")
		} else {
			for _, info := range allBook.BookInfos {
				//下载文件
				if info.Type != types.TypeBooksGroup && info.Type != types.TypeDir {
					//staticUrl := "/raw/" + info.BookID + "/" + url.QueryEscape(info.Name)
					staticUrl := "/raw/" + info.BookID + "/" + info.Name
					if checkUrlRegistered(info.BookID) {
						if config.Config.Debug {
							fmt.Println("路径已注册：", info)
						}
						continue
					} else {
						api.StaticFile(staticUrl, info.FilePath)
						staticUrlMap[info.BookID] = staticUrl
					}
				}
			}
		}
	}
}

// 1、设置web文件
func embedFile(engine *gin.Engine) {
	//使用自定义的模板引擎，命名为"template-data"，为了与VUE兼容，把左右分隔符自定义为 [[ ]]
	tmpl := template.Must(template.New("template-data").Delims("[[", "]]").Parse(TemplateString))
	//使用模板
	engine.SetHTMLTemplate(tmpl)
	if config.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		//gin.DisableConsoleColor()
		// 中间件，输出 log 到文件
		engine.Use(util.LoggerToFile(config.Config.LogFilePath, config.Config.LogFileName))
		//禁止控制台输出
		gin.DefaultWriter = io.Discard
	}

	//https://stackoverflow.com/questions/66248258/serve-embedded-filesystem-from-root-path-of-url
	assetsEmbedFS, err := fs.Sub(staticAssetFS, "static/assets")
	if err != nil {
		fmt.Println(err)
	}
	engine.StaticFS("/assets/", http.FS(assetsEmbedFS))
	imagesEmbedFS, errStaticImageFS := fs.Sub(staticImageFS, "static/images")
	if errStaticImageFS != nil {
		fmt.Println(errStaticImageFS)
	}
	engine.StaticFS("/images/", http.FS(imagesEmbedFS))

	engine.GET("/favicon.ico", func(c *gin.Context) {
		file, _ := staticFS.ReadFile("static/images/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})
	//用react写的后台界面：
	adminEmbedFS, errAdminFS := fs.Sub(adminFS, "admin")
	if errAdminFS != nil {
		fmt.Println(errAdminFS)
	}
	engine.StaticFS("/admin", http.FS(adminEmbedFS))

	//解析模板到HTML
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template-data", gin.H{
			"title": locale.GetString("HTML_TITLE") + config.Version, //页面标题
		})
	})
}
