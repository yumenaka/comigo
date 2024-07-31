package router

import (
	"embed"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
	"net/http"
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
