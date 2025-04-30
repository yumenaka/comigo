package reverse_proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/util/logger"
)

// ReverseProxyOptions 用于配置反向代理选项
type ReverseProxyOptions struct {
	TargetHost  string
	TargetPort  string
	RewritePath string
}

// ReverseProxyHandle TODO: 使用 echo 中间件实现反向代理
func ReverseProxyHandle(path string, option ReverseProxyOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 如果请求的 URI 以指定 path 开头，则进行代理请求
			if strings.HasPrefix(c.Request().RequestURI, path) {
				client := &http.Client{}

				// 将路径中的 RewritePath 替换移除，以构造目标地址
				requestUrl := strings.Replace(c.Request().RequestURI, option.RewritePath, "", 1)
				url := option.TargetHost + ":" + option.TargetPort + "/" + requestUrl

				req, err := http.NewRequest(c.Request().Method, url, c.Request().Body)
				if err != nil {
					logger.Infof("http.NewRequest Error: %s", err)
					return err
				}

				// 复制原请求的 Header
				req.Header = c.Request().Header.Clone()

				// 发送请求
				resp, err := client.Do(req)
				if err != nil {
					logger.Infof("client.Do Error: %s", err)
					return err
				}
				defer func() {
					if closeErr := resp.Body.Close(); closeErr != nil {
						logger.Infof("Body.Close() Error: %s", closeErr)
					}
				}()

				// 读取响应体
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Infof("io.ReadAll Error: %s", err)
					return err
				}

				// 复制响应头到 echo.Context
				for key, value := range resp.Header {
					if len(value) == 1 {
						c.Response().Header().Add(key, value[0])
					} else {
						for _, val := range value {
							c.Response().Header().Add(key, val)
						}
					}
				}

				// 设置响应状态码
				c.Response().WriteHeader(resp.StatusCode)

				// 写出响应内容
				if _, err = c.Response().Writer.Write(body); err != nil {
					logger.Infof("Response.Write Error: %s", err)
					return err
				}

				return nil
			}
			// 如果不匹配 path，则执行下一个处理
			return next(c)
		}
	}
}
