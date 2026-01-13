package reverse_proxy

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
)

// 只允许代理的前缀（外部访问路径）
const allowedPrefix = "/yumenaka/"

// GitHub 目标（固定）
const githubBase = "https://github.com"

// 预编译的正则表达式，用于替换 URL 中的 "latest" 为实际版本号
var (
	// 匹配 /releases/download/latest/ 路径
	releasePathRegex = regexp.MustCompile(`/releases/download/latest/`)
	// 匹配文件名中的 _latest_ 格式
	fileNameRegex = regexp.MustCompile(`_latest_`)
)

// ProxyHandler 处理反向代理请求
func ProxyHandler(c echo.Context) error {
	req := c.Request()
	cfg := config.GetCfg()
	// 安全检查：仅在调试模式下启用此功能
	if !cfg.Debug {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	// 安全检查：只允许 /yumenaka/*
	if !strings.HasPrefix(req.URL.Path, allowedPrefix) {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	// 两种模式：
	// 1) 普通路径：根据路径前缀映射到不同的 GitHub 服务
	// 2) 跟随重写跳转：携带 ?u=... （base64url） -> 直接抓取 u 指向的 https 资源并回传
	if uParam := req.URL.Query().Get("u"); uParam != "" {
		target, err := decodeURL(uParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad u param")
		}
		return fetchAndStream(c, target)
	}

	// 根据路径前缀确定目标 URL
	path := req.URL.Path
	var target string

	// 检查是否是 raw.githubusercontent.com 的代理路径
	// /yumenaka/raw.githubusercontent.com/... -> https://raw.githubusercontent.com/...
	if strings.HasPrefix(path, "/yumenaka/raw.githubusercontent.com/") {
		// 提取 raw.githubusercontent.com 之后的路径
		rawPath := strings.TrimPrefix(path, "/yumenaka/raw.githubusercontent.com")
		target = "https://raw.githubusercontent.com" + rawPath
	} else if strings.HasPrefix(path, "/yumenaka/api.github.com/") {
		// 检查是否是 api.github.com 的代理路径
		// /yumenaka/api.github.com/... -> https://api.github.com/...
		apiPath := strings.TrimPrefix(path, "/yumenaka/api.github.com")
		target = "https://api.github.com" + apiPath
	} else {
		// 默认映射到 github.com
		// /yumenaka/xxx -> https://github.com/yumenaka/xxx
		target = githubBase + path
	}

	// 替换 URL 中的 "latest" 为当前版本号
	// 支持 releases/download/latest/... 和文件名中的 _latest_ 格式
	target = replaceLatestWithVersion(target)

	if req.URL.RawQuery != "" {
		target += "?" + req.URL.RawQuery
	}

	return fetchAndStream(c, target)
}

// fetchAndStream 获取目标资源并流式传输
func fetchAndStream(c echo.Context, target string) error {
	// 只允许 https
	tu, err := url.Parse(target)
	if err != nil || tu.Scheme != "https" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid target")
	}

	// 检查目标主机是否在允许列表中（避免 SSRF）
	// 允许的主机：github.com, raw.githubusercontent.com, api.github.com 以及 GitHub release 相关的可信域名
	if tu.Host != "github.com" && tu.Host != "raw.githubusercontent.com" && tu.Host != "api.github.com" {
		if !isAllowedRedirectHost(tu.Host) {
			return echo.NewHTTPError(http.StatusForbidden, "redirect host not allowed")
		}
	}

	upReq, err := http.NewRequestWithContext(c.Request().Context(), c.Request().Method, target, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	// 透传一些常用头（可按需增减）
	copyHeader(upReq.Header, c.Request().Header, []string{
		"Range",
		"If-Modified-Since",
		"If-None-Match",
		"User-Agent",
		"Accept",
		"Accept-Encoding", // 注意：我们这里直接透传 body，不做解压缩；对下载一般没问题
	})

	cli := newHTTPClient()

	// 不自动跟随跳转：我们要把 302 Location 改写回自己的域名
	resp, err := cli.Do(upReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()

	// 处理 301/302/307/308：把 Location 重写到 publicBase，并塞进 ?u=
	if isRedirect(resp.StatusCode) {
		loc := resp.Header.Get("Location")
		if loc == "" {
			return echo.NewHTTPError(http.StatusBadGateway, "redirect without location")
		}

		// Location 可能是相对路径，这里补齐为 https://github.com/...
		locURL, err := resolveLocation(target, loc)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, "bad redirect location")
		}

		// 仍然做一次 host 白名单，防止被引去奇怪的地方
		if locURL.Scheme != "https" {
			return echo.NewHTTPError(http.StatusForbidden, "redirect scheme not allowed")
		}
		if locURL.Host != "github.com" && locURL.Host != "raw.githubusercontent.com" && locURL.Host != "api.github.com" {
			if !isAllowedRedirectHost(locURL.Host) {
				return echo.NewHTTPError(http.StatusForbidden, "redirect host not allowed")
			}
		}

		// 获取公共基础 URL（用于重定向）
		publicBase := getPublicBase(c)

		// 让客户端继续请求我们的域名：/yumenaka/... ?u=<base64url(最终下载URL)>
		// 这样用户看到的始终是我们的域名
		next := publicBase + c.Request().URL.Path
		q := c.Request().URL.Query()
		q.Set("u", encodeURL(locURL.String()))
		// 保留原始查询参数（如果你不想保留可以清掉）
		// q.Del("...")

		redir := next + "?" + q.Encode()
		return c.Redirect(resp.StatusCode, redir)
	}

	// 非跳转：开始回传
	// 状态码
	c.Response().WriteHeader(resp.StatusCode)

	// 复制响应头（过滤掉一些 hop-by-hop）
	copyResponseHeader(c.Response().Header(), resp.Header)

	// 直接流式拷贝（支持大文件）
	_, copyErr := io.Copy(c.Response().Writer, resp.Body)
	if copyErr != nil && !errors.Is(copyErr, context.Canceled) {
		// 客户端断开等情况会出现 context canceled，不当成错误
		return copyErr
	}
	return nil
}

// newHTTPClient 创建用于反向代理的 HTTP 客户端
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// 两边都是 https
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0, // 下载可能很大，不设总超时；依赖上游/下游超时和反向代理层
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁止自动跟随
			return http.ErrUseLastResponse
		},
	}
}

// isRedirect 判断状态码是否为重定向
func isRedirect(code int) bool {
	return code == http.StatusMovedPermanently ||
		code == http.StatusFound ||
		code == http.StatusTemporaryRedirect ||
		code == http.StatusPermanentRedirect
}

// resolveLocation 解析相对路径的 Location 头
func resolveLocation(baseStr, loc string) (*url.URL, error) {
	baseURL, err := url.Parse(baseStr)
	if err != nil {
		return nil, err
	}
	locURL, err := url.Parse(loc)
	if err != nil {
		return nil, err
	}
	return baseURL.ResolveReference(locURL), nil
}

// isAllowedRedirectHost 检查重定向目标主机是否在白名单中
func isAllowedRedirectHost(host string) bool {
	// GitHub releases 常见跳转域名（按需补充）
	switch strings.ToLower(host) {
	case "objects.githubusercontent.com",
		"release-assets.githubusercontent.com",
		"github-releases.githubusercontent.com",
		"githubusercontent.com",
		"raw.githubusercontent.com",
		"api.github.com":
		return true
	default:
		// 也允许 *.githubusercontent.com（有时会出现子域）
		if strings.HasSuffix(strings.ToLower(host), ".githubusercontent.com") {
			return true
		}
		return false
	}
}

// copyHeader 复制指定的请求头
func copyHeader(dst http.Header, src http.Header, keys []string) {
	for _, k := range keys {
		if v := src.Values(k); len(v) > 0 {
			// 覆盖式设置
			dst.Del(k)
			for _, vv := range v {
				dst.Add(k, vv)
			}
		}
	}
}

// copyResponseHeader 复制响应头（过滤 hop-by-hop headers）
func copyResponseHeader(dst http.Header, src http.Header) {
	// hop-by-hop header 必须丢弃
	hopByHop := map[string]bool{
		"Connection":          true,
		"Proxy-Connection":    true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailer":             true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}

	for k, vv := range src {
		if hopByHop[http.CanonicalHeaderKey(k)] {
			continue
		}
		dst.Del(k)
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// encodeURL 使用 base64url 编码 URL（无 padding）
func encodeURL(s string) string {
	// base64url（无 padding）避免 query 太丑/需要额外转义
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// decodeURL 解码 base64url 编码的 URL
func decodeURL(s string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GetComigoScriptHandler 处理 /get.sh 的简化路径
func GetComigoScriptHandler(c echo.Context) error {
	cfg := config.GetCfg()
	// 安全检查：仅在调试模式下启用此功能
	if !cfg.Debug {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	// 直接映射到 get.sh 脚本
	target := "https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh"
	if c.Request().URL.RawQuery != "" {
		target += "?" + c.Request().URL.RawQuery
	}

	return fetchAndStream(c, target)
}

// replaceLatestWithVersion 将 URL 中的 "latest" 替换为当前版本号
// 支持两种模式：
// 1. releases/download/latest/ -> releases/download/v1.2.4/
// 2. 文件名中的 _latest_ -> _v1.2.4_ (如 comi_latest_linux_amd64.tar.gz -> comi_v1.2.4_linux_amd64.tar.gz)
func replaceLatestWithVersion(target string) string {
	version := config.GetVersion() // 从 version.go 获取版本号，如 "v1.2.4"

	// 模式1：替换 releases/download/latest/ 路径中的 latest
	// 匹配 /releases/download/latest/ 并替换为 /releases/download/{version}/
	if releasePathRegex.MatchString(target) {
		target = releasePathRegex.ReplaceAllString(target, "/releases/download/"+version+"/")
	}

	// 模式2：替换文件名中的 _latest_ 格式
	// 例如：comi_latest_linux_amd64.tar.gz -> comi_v1.2.4_linux_amd64.tar.gz
	if fileNameRegex.MatchString(target) {
		target = fileNameRegex.ReplaceAllString(target, "_"+version+"_")
	}

	return target
}

// getPublicBase 获取公共基础 URL（用于重定向）
func getPublicBase(c echo.Context) string {
	cfg := config.GetCfg()
	host := cfg.Host

	// 如果配置了 Host，使用配置的 Host
	if host != "" {
		// 判断是否启用 TLS
		if cfg.EnableTLS || cfg.AutoTLSCertificate {
			return "https://" + host
		}
		return "http://" + host
	}

	// 否则从请求中获取
	scheme := "http"
	if c.IsTLS() {
		scheme = "https"
	}
	return scheme + "://" + c.Request().Host
}
