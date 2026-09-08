package routers

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yumenaka/comigo/tools/traffic"
)

// 验证统计压缩后的正文，并且上传读取和错误输出均经过统计器。
func TestTrafficGzipAndErrors(t *testing.T) {
	count := traffic.New()
	e := echo.New()
	e.Pre(trackTraffic(count))
	e.Use(middleware.Gzip())
	body := strings.Repeat("comigo", 2000)
	e.POST("/", func(c echo.Context) error {
		_, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		return c.String(200, body)
	})
	req := httptest.NewRequest("POST", "/", strings.NewReader("upload"))
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	s := count.Snapshot()
	if s.ReceivedBytes != 6 || s.SentBytes != uint64(rec.Body.Len()) || rec.Body.Len() >= len(body) {
		t.Fatalf("gzip: %+v, bytes %d", s, rec.Body.Len())
	}
	reader, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	raw, err := io.ReadAll(reader)
	if err != nil || string(raw) != body {
		t.Fatal("gzip content changed", err)
	}
	e.GET("/error", func(c echo.Context) error { return echo.ErrForbidden })
	before := count.Snapshot().SentBytes
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest("GET", "/error", nil))
	if rec.Code != 403 || count.Snapshot().SentBytes-before != uint64(rec.Body.Len()) {
		t.Fatal("error response not counted")
	}
}

// SSE 的 Flush 必须可用，且第一块正文写入后立即计数。
func TestTrafficStreamingAndRange(t *testing.T) {
	count := traffic.New()
	e := echo.New()
	e.Pre(trackTraffic(count))
	e.GET("/events", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/event-stream")
		_, err := c.Response().Write([]byte("data: ready\n\n"))
		if err != nil {
			return err
		}
		c.Response().Flush()
		if count.Snapshot().SentBytes != 13 {
			t.Error("stream delayed until request completion")
		}
		return nil
	})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest("GET", "/events", nil))
	if !rec.Flushed {
		t.Fatal("flush lost")
	}
	e.GET("/file", func(c echo.Context) error {
		http.ServeContent(c.Response(), c.Request(), "book", count.Snapshot().StartedAt, strings.NewReader("0123456789"))
		return nil
	})
	req := httptest.NewRequest("GET", "/file", nil)
	req.Header.Set("Range", "bytes=2-5")
	before := count.Snapshot().SentBytes
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != 206 || rec.Body.String() != "2345" || count.Snapshot().SentBytes-before != 4 {
		t.Fatal("range changed")
	}
}
