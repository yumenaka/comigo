package settings_page

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/util/logger"
)

// -------------------------
// 使用templ模板响应htmx请求
// -------------------------

func Tab1(c echo.Context) error {
	template := tab1(&state.Global)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, template); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func Tab2(c echo.Context) error {
	template := tab2(&state.Global)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, template); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func Tab3(c echo.Context) error {
	template := tab3(&state.Global)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, template); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// -------------------------
// 抽取更新配置的公共逻辑
// -------------------------

// parseSingleHTMXFormPair 提取并返回表单中的"第一对" key/value
func parseSingleHTMXFormPair(c echo.Context) (string, string, error) {
	if !htmx.IsHTMX(c.Request()) {
		return "", "", errors.New("non-htmx request")
	}
	formData, err := c.FormParams()
	if err != nil {
		return "", "", fmt.Errorf("parseForm error: %v", err)
	}
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
func updateConfigGeneric(c echo.Context) (string, string, error) {
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
func UpdateStringConfigHandler(c echo.Context) error {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	updatedHTML := StringConfig(name, newValue, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateBoolConfigHandler 处理 Bool 类型
func UpdateBoolConfigHandler(c echo.Context) error {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	boolVal, parseErr := strconv.ParseBool(newValue)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 bool: %v", newValue, parseErr)
		return echo.NewHTTPError(http.StatusBadRequest, "parse bool error")
	}
	updatedHTML := BoolConfig(name, boolVal, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateNumberConfigHandler 处理 Number 类型
func UpdateNumberConfigHandler(c echo.Context) error {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 将字符串解析为 int
	intVal, parseErr := strconv.ParseInt(newValue, 10, 64)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 int: %v", newValue, parseErr)
		return echo.NewHTTPError(http.StatusBadRequest, "parse int error")
	}
	// 渲染对应的模板
	updatedHTML := NumberConfig(name, int(intVal), name+"_Description", 0, 65535)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

func AddArrayConfigHandler(c echo.Context) error {
	// 解析htmx的form数据
	// 原始数据类似：configName=SupportMediaType&addValue=.test
	// 解析后的数据类似：ConfigName=SupportMediaType, AddValue=.test
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}

	configName := c.FormValue("configName")
	addValue := c.FormValue("addValue")

	logger.Infof("AddArrayConfigHandler: %s = %s\n", configName, addValue)

	values, err := doAdd(configName, addValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "add error")
	}

	updatedHTML := StringArrayConfig(configName, values, configName+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
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
func DeleteArrayConfigHandler(c echo.Context) error {
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}

	configName := c.FormValue("configName")
	deleteValue := c.FormValue("deleteValue")

	logger.Infof("DeleteArrayConfigHandler: %s = %s\n", configName, deleteValue)

	values, err := doDelete(configName, deleteValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, configName+" delete Failed")
	}

	updatedHTML := StringArrayConfig(configName, values, configName+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
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
func HandleConfigSave(c echo.Context) error {
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}
	// 保存到什么文件夹
	selectedDir := c.FormValue("selectedDir")
	if selectedDir == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No directory selected")
	}
	if selectedDir != config.WorkingDirectory && selectedDir != config.HomeDirectory && selectedDir != config.ProgramDirectory {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid directory selected")
	}
	// 保存失败时
	if err := config.SaveConfig(selectedDir); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save config")
	}
	// 更新后的HTML片段
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// HandleConfigDelete 处理 /api/config-delete 的 POST 请求
func HandleConfigDelete(c echo.Context) error {
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}
	// 保存到什么文件夹
	selectedDir := c.FormValue("selectedDir")
	if selectedDir == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No directory selected")
	}
	if selectedDir != config.WorkingDirectory && selectedDir != config.HomeDirectory && selectedDir != config.ProgramDirectory {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid directory selected")
	}

	if err := config.DeleteConfigIn(selectedDir); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save config")
	}
	// 更新后的HTML片段
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}
