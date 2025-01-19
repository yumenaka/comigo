package settings_page

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/util/logger"
)

// -------------------------
// 使用templ模板响应htmx请求
// -------------------------

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

// -------------------------
// 抽取更新配置的公共逻辑
// -------------------------

// parseSingleHTMXFormPair 提取并返回表单中的“第一对” key/value。
// 如果不是 HTMX 请求或解析失败等情况，返回对应的错误。
func parseSingleHTMXFormPair(c *gin.Context) (string, string, error) {
	// 手动测试错误
	//return "", "", errors.New("test")
	if !htmx.IsHTMX(c.Request) {
		return "", "", errors.New("non-htmx request")
	}
	if err := c.Request.ParseForm(); err != nil {
		return "", "", fmt.Errorf("parseForm error: %v", err)
	}
	formData := c.Request.PostForm
	if len(formData) == 0 {
		return "", "", errors.New("no form data")
	}
	var name, newValue string
	for key, values := range formData {
		name = key
		if len(values) > 0 {
			newValue = values[0]
		}
		break
	}
	return name, newValue, nil
}

// updateConfigGeneric 抽取了更新配置的大部分公共逻辑，返回 name 和 newValue 等。
func updateConfigGeneric(c *gin.Context) (string, string, error) {
	name, newValue, err := parseSingleHTMXFormPair(c)
	if err != nil {
		return "", "", err
	}

	logger.Infof("Update config: %s = %s", name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := state.ServerConfig.SetConfigValue(name, newValue); setErr != nil {
		logger.Errorf("Failed to set config value: %v", setErr)
		return "", "", setErr
	}

	// 写入配置文件
	if writeErr := config.WriteConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return name, newValue, nil
}

// -------------------------
// 各类配置的更新 Handler
// -------------------------

// UpdateStringConfigHandler 处理 String 类型
func UpdateStringConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 渲染对应的模板
	updatedHTML := StringConfig(name, newValue, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// UpdateBoolConfigHandler 处理 Bool 类型
func UpdateBoolConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 将字符串解析为 bool
	boolVal, parseErr := strconv.ParseBool(newValue)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 bool: %v", newValue, parseErr)
		c.String(http.StatusBadRequest, "parse bool error")
		return
	}
	// 渲染对应的模板
	updatedHTML := BoolConfig(name, boolVal, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// UpdateNumberConfigHandler 处理 Number 类型
func UpdateNumberConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 将字符串解析为 int
	intVal, parseErr := strconv.ParseInt(newValue, 10, 64)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 int: %v", newValue, parseErr)
		c.String(http.StatusBadRequest, "parse int error")
		return
	}
	// 渲染对应的模板
	updatedHTML := NumberConfig(name, int(intVal), name+"_Description", 0, 65535)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func AddArrayConfigHandler(c *gin.Context) {
	// 解析htmx的form数据
	// 原始数据类似：configName=SupportMediaType&addValue=.test
	// 解析后的数据类似：ConfigName=SupportMediaType, AddValue=.test
	if !htmx.IsHTMX(c.Request) {
		c.String(http.StatusBadRequest, "non-htmx request")
	}
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "parseForm error")
		return
	}
	formData := c.Request.PostForm
	if len(formData) == 0 {
		c.String(http.StatusBadRequest, "no form data")
		return
	}
	var ConfigName, AddValue string
	for key, values := range formData {
		if key == "configName" {
			ConfigName = values[0]
		}
		if key == "addValue" {
			AddValue = values[0]
		}
	}
	// 打印日志
	logger.Infof("AddArrayConfigHandler: %s = %s\n", ConfigName, AddValue)
	// 根据 ConfigName 找到对应配置，并添加 AddValue
	values, err := doAdd(ConfigName, AddValue)
	if err != nil {
		c.String(http.StatusInternalServerError, "add error")
		return
	}
	// 返回更新后的片段
	updatedHTML := StringArrayConfig(ConfigName, values, ConfigName+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func doAdd(configName, addValue string) ([]string, error) {
	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	values, err := config.GetCfg().AddStringArrayConfig(configName, addValue)
	if err != nil {
		logger.Errorf("Failed to add config value: %v", err)
		return nil, err
	}
	// 写入配置文件
	if writeErr := config.WriteConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	return values, nil
}

// DeleteArrayConfigHandler 处理删除数组元素
func DeleteArrayConfigHandler(c *gin.Context) {
	// 解析htmx的form数据
	// 原始数据类似：configName=SupportMediaType&deleteValue=.test
	// 解析后的数据类似：ConfigName=SupportMediaType, DeleteValue=.test
	if !htmx.IsHTMX(c.Request) {
		c.String(http.StatusBadRequest, "non-htmx request")
	}
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "parseForm error")
		return
	}
	formData := c.Request.PostForm
	if len(formData) == 0 {
		c.String(http.StatusBadRequest, "no form data")
		return
	}
	var ConfigName, DeleteValue string
	for key, values := range formData {
		if key == "configName" {
			ConfigName = values[0]
		}
		if key == "deleteValue" {
			DeleteValue = values[0]
		}
	}
	// 打印日志
	logger.Infof("DeleteArrayConfigHandler: %s = %s\n", ConfigName, DeleteValue)

	// 根据 ConfigName 找到对应配置，
	// 然后根据 Index 与 DeleteValue 删除该元素
	values, err := doDelete(ConfigName, DeleteValue)
	if err != nil {
		c.String(http.StatusInternalServerError, ConfigName+" delete Failed")
		return
	}
	// 最后把更新后的同一段 HTML 片段返回给前端
	updatedHTML := StringArrayConfig(ConfigName, values, ConfigName+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func doDelete(configName string, deleteValue string) ([]string, error) {
	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	values, err := state.ServerConfig.DeleteStringArrayConfig(configName, deleteValue)
	if err != nil {
		return nil, err
	}

	// 写入配置文件
	if writeErr := config.WriteConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	return values, nil
}

// HandleConfigSave 处理 /api/config-save 的 POST 请求
func HandleConfigSave(c *gin.Context) {
	if !htmx.IsHTMX(c.Request) {
		c.String(http.StatusBadRequest, "non-htmx request")
	}
	// 从表单获取选中的目录
	selectedDir := c.PostForm("selectedDir")
	if selectedDir == "" {
		// 如果前端没传，返回错误信息
		c.String(http.StatusBadRequest, "No directory selected")
		return
	}
	if selectedDir != config.WorkingDirectory && selectedDir != config.HomeDirectory && selectedDir != config.ProgramDirectory {
		// 如果不是三个目录之一，就不能保存
		c.String(http.StatusBadRequest, "Invalid directory selected")
		return
	}
	if err := config.SaveConfig(selectedDir); err != nil {
		// 保存失败，返回错误信息
		c.String(http.StatusInternalServerError, "Failed to save config")
		return
	}
	// 最后把更新后的同一段 HTML 片段返回给前端
	updatedHTML := ConfigManager(config.CheckConfigLocation(), config.GetExistingConfigFilePath())
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// HandleConfigDelete 处理 /api/config-delete 的 POST 请求
func HandleConfigDelete(c *gin.Context) {
	if !htmx.IsHTMX(c.Request) {
		c.String(http.StatusBadRequest, "non-htmx request")
	}
	// 从表单获取选中的目录
	selectedDir := c.PostForm("selectedDir")
	if selectedDir == "" {
		// 如果前端没传，返回错误信息
		c.String(http.StatusBadRequest, "No directory selected")
		return
	}
	if selectedDir != config.WorkingDirectory && selectedDir != config.HomeDirectory && selectedDir != config.ProgramDirectory {
		// 如果不是三个目录之一，就不能保存
		c.String(http.StatusBadRequest, "Invalid directory selected")
		return
	}
	if err := config.DeleteConfigIn(selectedDir); err != nil {
		// 删除失败，返回错误信息
		c.String(http.StatusInternalServerError, "Failed to save config")
		return
	}
	// 最后把更新后的同一段 HTML 片段返回给前端
	updatedHTML := ConfigManager(config.CheckConfigLocation(), config.GetExistingConfigFilePath())
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
