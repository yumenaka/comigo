package settings_page

import (
	"errors"
	"fmt"
	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/util/logger"
	"net/http"
	"strconv"
)

// 使用模板中响应htmx请求，页面比较复杂时用
func Tab1(c *gin.Context) {
	template := tab1(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab2(c *gin.Context) {
	template := tab2(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab3(c *gin.Context) {
	template := tab3(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// parseSingleHTMXFormPair 提取并返回表单中的“第一对” key/value。
// 如果不是 HTMX 请求或解析失败等情况，返回对应的错误。
func parseSingleHTMXFormPair(c *gin.Context) (string, string, error) {
	// 1. 仅接收 HTMX 请求
	if !htmx.IsHTMX(c.Request) {
		return "", "", errors.New("non-htmx request")
	}
	// 2. 解析表单
	if err := c.Request.ParseForm(); err != nil {
		return "", "", fmt.Errorf("parseForm error: %v", err)
	}
	// 3. 检查是否有表单数据
	formData := c.Request.PostForm
	if len(formData) == 0 {
		return "", "", errors.New("no form data")
	}

	// 4. 假设只有一对数据 (key=value)，只取第一对就行
	var name, newValue string
	for key, values := range formData {
		name = key
		if len(values) > 0 {
			newValue = values[0] // 一般只有一个
		}
		break
	}

	return name, newValue, nil
}

// UpdateStringConfigHandler 处理 htmx 请求
func UpdateStringConfigHandler(c *gin.Context) {
	name, newValue, err := parseSingleHTMXFormPair(c)
	if err != nil {
		// 这里根据需要决定是 c.String 还是 c.AbortWithError
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof("Update config: %s = %s", name, newValue)

	// 更新配置
	if err := state.ServerConfig.SetConfigValue(name, newValue); err != nil {
		logger.Errorf("Failed to set config value: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 渲染对应的模板
	updatedHTML := StringConfig(name, newValue, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func UpdateBoolConfigHandler(c *gin.Context) {
	name, newValue, err := parseSingleHTMXFormPair(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof("Update config: %s = %s", name, newValue)

	// 更新配置（先保存字符串形式）
	if err := state.ServerConfig.SetConfigValue(name, newValue); err != nil {
		logger.Errorf("Failed to set config value: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 把字符串形式转换为bool
	boolVal, err := strconv.ParseBool(newValue)
	if err != nil {
		logger.Errorf("无法将 '%s' 解析为 bool: %v", newValue, err)
		// 看需求决定返回什么
		c.String(http.StatusBadRequest, "parse bool error")
		return
	}

	// 渲染对应的模板
	updatedHTML := BoolConfig(name, boolVal, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func UpdateNumberConfigHandler(c *gin.Context) {
	name, newValue, err := parseSingleHTMXFormPair(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof("Update config: %s = %s", name, newValue)

	// 更新配置（字符串形式）
	if err := state.ServerConfig.SetConfigValue(name, newValue); err != nil {
		logger.Errorf("Failed to set config value: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 转成数字
	intVal, err := strconv.ParseInt(newValue, 10, 64)
	if err != nil {
		logger.Errorf("无法将 '%s' 解析为 int: %v", newValue, err)
		c.String(http.StatusBadRequest, "parse int error")
		return
	}

	// 渲染对应的模板
	updatedHTML := NumberConfig(name, int(intVal), name+"_Description", 0, 65535)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
