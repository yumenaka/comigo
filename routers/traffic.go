package routers

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/tools/traffic"
)

type trafficBody struct {
	io.ReadCloser
	counter *traffic.Counter
}

func (r *trafficBody) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	r.counter.Add(n, 0)
	return n, err
}

type trafficWriter struct {
	http.ResponseWriter
	counter *traffic.Counter
}

func (w *trafficWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.counter.Add(0, n)
	return n, err
}
func (w *trafficWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
func (w *trafficWriter) Flush()                      { _ = http.NewResponseController(w.ResponseWriter).Flush() }
func (w *trafficWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

// trackTraffic 位于 gzip 和错误处理之外，长响应每次写出时即入账。
func trackTraffic(counter *traffic.Counter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Body != nil {
				c.Request().Body = &trafficBody{c.Request().Body, counter}
			}
			c.Response().Writer = &trafficWriter{c.Response().Writer, counter}
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}
