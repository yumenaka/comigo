package assets

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

//go:embed dist/* static/*
var Frontend embed.FS
var FrontendFS fs.FS

//go:embed images/*
var Images embed.FS
var ImagesFS fs.FS

//go:embed epub/*
var Epub embed.FS

//go:embed pwa/*
var Pwa embed.FS

//go:embed robots.txt
var Robots embed.FS

// GetCSS 在页面中插入需要的css代码
func GetCSS(oneFileMode bool) (cssString string) {
	if oneFileMode {
		cssString = "<style>" + GetFileStr("dist/styles.css") + "</style>\n"
	} else {
		cssString = "<link rel=\"stylesheet\" href=\"" + config.PrefixPath("/assets/dist/styles.css") + "\">\n"
	}
	return cssString
}

// GetBasePathScript 暴露前端路径工具，供静态 JS 与模板内联脚本统一处理反向代理基础路径。
func GetBasePathScript() string {
	basePath, _ := json.Marshal(config.GetBasePath())
	return `<script>
window.ComiGoBasePath = ` + string(basePath) + `;
window.ComiGoPath = function(path) {
  const base = window.ComiGoBasePath || '';
  if (!path) return base ? base + '/' : '/';
  if (/^(https?:|data:|blob:|#)/.test(path)) return path;
  const normalized = path.startsWith('/') ? path : '/' + path;
  if (!base) return normalized;
  return normalized === '/' ? base + '/' : base + normalized;
};
window.ComiGoRelativePath = function(pathname) {
  const base = window.ComiGoBasePath || '';
  const path = pathname || window.location.pathname || '/';
  if (!base) return path;
  if (path === base) return '/';
  return path.startsWith(base + '/') ? path.slice(base.length) : path;
};
</script>
`
}

// GetJavaScript 在页面中插入需要的js代码
func GetJavaScript(oneFileMode bool, insertScript []string) (jsString string) {
	jsString = GetBasePathScript()
	// <!-- 通用js代码,初始化htmx、Alpine等第三方库  -->
	if oneFileMode {
		jsString += "<script>" + GetFileStr("dist/main.js") + "</script>\n"
	} else {
		jsString += "<script src=\"" + config.PrefixPath("/assets/dist/main.js") + "\"></script>\n"
	}
	// <!-- 每个页面的特殊js代码片段  -->
	for _, script := range insertScript {
		if oneFileMode {
			jsString += "<script>" + GetFileStr(script) + "</script>\n"
		} else {
			jsString += "<script src=\"" + config.PrefixPath("/assets/"+script) + "\"></script>\n"
		}
	}
	// fmt.Println("jsString:", jsString)
	return jsString
}

// GetFileStr 从Static获取字符串形式的脚本
func GetFileStr(filePath string) string {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Frontend.ReadFile(filePath)
	if err != nil {
		return "Not Found Script:" + filePath
	}
	str := string(data)
	// 在把文件内容注入模板之前做一次替换,避免</script>或 </body> 提前把 <script> 标签“截断”
	safe := strings.ReplaceAll(str, "</script", "<\\/script")
	safe = strings.ReplaceAll(str, "</body", "<\\/body")
	// 将内容转换为字符串并返回
	return safe
}

// GetImageSrc 从Static获取Base64编码的图片的src属性
func GetImageSrc(filePath string) string {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Frontend.ReadFile(filePath)
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_read_embedded_image"), err)
		return "Not Found Image:" + filePath
	}

	// 获取文件扩展名，并猜测MIME类型
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Base64 编码图片数据
	base64Data := base64.StdEncoding.EncodeToString(data)

	// 生成图片的 src 属性
	src := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

	return src
}

// GetData 从Static获取字节切片形式的数据
func GetData(filePath string) []byte {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Frontend.ReadFile(filePath)
	if err != nil {
		// 如果有错误发生，返回空的字节切片，并输出错误信息
		logger.Errorf(locale.GetString("err_failed_to_read_embedded_data"), err)
		return []byte{}
	}
	// 返回文件内容作为字节切片
	return data
}

// GetDataBase64 从 Static 获取 Base64 字符串，便于把 wasm 等二进制资源内联到静态 HTML。
func GetDataBase64(filePath string) string {
	return base64.StdEncoding.EncodeToString(GetData(filePath))
}

// GetImageDataSrc 从 Images embed.FS 获取图片并返回 data URL，便于静态 HTML 内联图片资源。
func GetImageDataSrc(imageName string) string {
	data := GetImageData(imageName)
	if len(data) == 0 {
		return ""
	}
	ext := filepath.Ext(imageName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
}

// GetImageData 从Images embed.FS获取图片字节数据
func GetImageData(imageName string) []byte {
	filePath := "images/" + imageName
	// 使用ReadFile从嵌入文件系统中读取图片内容
	data, err := Images.ReadFile(filePath)
	if err != nil {
		// 如果有错误发生，返回空的字节切片，并输出错误信息
		logger.Errorf(locale.GetString("err_failed_to_read_embedded_image"), err)
		return []byte{}
	}
	// 返回图片内容作为字节切片
	return data
}
