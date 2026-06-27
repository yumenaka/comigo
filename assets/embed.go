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
	wailsBuild, _ := json.Marshal(isWailsBuild())
	debugMode, _ := json.Marshal(config.GetCfg().Debug)
	wailsRuntimeScript := ""
	if isWailsBuild() {
		// Wails3 runtime 负责通知窗口 ready；普通 Web 构建不插入这个脚本。
		wailsRuntimeScript = "<script type=\"module\" src=\"/wails/runtime.js\"></script>\n"
	}
	return `<script>
window.ComiGoBasePath = ` + string(basePath) + `;
window.ComiGoWails = ` + string(wailsBuild) + `;
window.ComiGoDebug = ` + string(debugMode) + `;
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
window.ComiGoElectron = (function() {
  const key = 'ComiGoElectron';
  const params = new URLSearchParams(window.location.search || '');
  const flag = params.get('comigo_electron');
  try {
    if (flag === '1') {
      window.localStorage.setItem(key, '1');
      return true;
    }
    if (flag === '0') {
      window.localStorage.removeItem(key);
      return false;
    }
    return window.localStorage.getItem(key) === '1';
  } catch (_) {
    return flag === '1';
  }
})();
window.ComiGoElectronAction = function(action) {
  if (!window.ComiGoElectron || !action) return false;
  const normalizedAction = String(action).trim();
  if (!normalizedAction) return false;
  // Electron 外壳会拦截这个自定义协议，不会真正离开当前页面，也不会给 WebView 暴露 Node 能力。
  window.location.href = 'comigo-electron://' + encodeURIComponent(normalizedAction);
  return true;
};
window.ComiGoIsWails = function() {
  return window.location.protocol === 'wails:' || !!window.WailsInvoke || !!window._wails || !!window.wails || !!window.go?.main?.App;
};
if (window.ComiGoIsWails()) {
  // Wails dev/debug 默认会放开右键菜单，这里只在桌面壳内统一禁用。
  window.addEventListener('contextmenu', function(event) {
    event.preventDefault();
  }, { capture: true });
}
window.ComiGoShareURL = function(currentURL, publicBaseURL) {
  const target = new URL(currentURL || window.location.href);
  const publicBase = new URL(publicBaseURL || window.location.origin + '/');
  if (!publicBase.pathname.endsWith('/')) publicBase.pathname += '/';
  if (window.ComiGoIsWails() || target.protocol === 'wails:') {
    const relative = window.ComiGoRelativePath(target.pathname).replace(/^\/+/, '') + target.search + target.hash;
    return new URL(relative, publicBase).toString();
  }
  const host = target.hostname;
  // 只有本机地址才替换，反向代理或远程域名保持当前访问来源。
  if (host === 'localhost' || host === '::1' || host === '[::1]' || host.startsWith('127.')) {
    target.protocol = publicBase.protocol;
    target.host = publicBase.host;
  }
  return target.toString();
};
// Wails WebView 内的新窗口链接由宿主打开，普通浏览器保留 window.open 行为。
window.ComiGoOpenExternalURL = async function(targetURL) {
  if (!targetURL) return;
  if (window.ComiGoWails) {
    try {
      const response = await fetch(window.ComiGoPath('/api/wails/open-url'), {
        method: 'POST',
        credentials: 'same-origin',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: targetURL }),
      });
      if (response.ok) return;
    } catch (_) {}
  }
  window.open(targetURL, '_blank', 'noopener');
};
window.ComiGoToggleFullscreen = async function() {
  if (window.ComiGoElectronAction('toggle-fullscreen')) return;
  if (window.ComiGoWails) {
    try {
      const response = await fetch(window.ComiGoPath('/api/wails/toggle-fullscreen'), {
        method: 'POST',
        credentials: 'same-origin',
      });
      if (response.ok) return;
    } catch (_) {}
  }
  if (window.Screenfull && window.Screenfull.isEnabled) {
    window.Screenfull.toggle();
  } else {
    showToast(i18next.t('not_support_fullscreen'));
  }
};
</script>
` + wailsRuntimeScript
}

// GetJavaScript 在页面中插入需要的js代码
func GetJavaScript(oneFileMode bool, insertScript []string) (jsString string) {
	// 插入通用 JS 代码，初始化 Alpine 与页面交互等前端基础能力。
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
	return jsString
}

// GetFileStr 获取字符串形式的脚本
func GetFileStr(filePath string) string {
	// 使用ReadFile从嵌入文件系统中读取文件内容
	data, err := Frontend.ReadFile(filePath)
	if err != nil {
		return "Not Found Script:" + filePath
	}
	// 在把文件内容注入模板之前做一次替换,避免</script>或 </body> 提前把 <script> 标签“截断”
	str := string(data)
	safe := strings.ReplaceAll(str, "</script", "<\\/script")
	safe = strings.ReplaceAll(safe, "</body", "<\\/body")
	// 将内容转换为字符串并返回
	return safe
}

// GetImageSrc 获取Base64编码的图片的src属性
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

// GetData 获取字节切片形式的数据
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
