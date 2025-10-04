package settings

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/templ/state"
	"github.com/yumenaka/comigo/tools/logger"
)

// -------------------------
// 使用templ模板响应htmx请求
// -------------------------

//func AllSetting(c echo.Context) error {
//	tsStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
//	if err != nil {
//		return err
//	}
//	template := settings_all(tsStatus)
//	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, template); renderErr != nil {
//		return echo.NewHTTPError(http.StatusInternalServerError)
//	}
//	return nil
//}

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
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return name, newValue, nil
}

// -------------------------
// 各类配置的更新 PageHandler
// -------------------------

// UpdateStringConfigHandler 处理 String 类型
func UpdateStringConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
		saveSuccessHint = true
	}
	updatedHTML := StringConfig(name, newValue, name+"_Description", saveSuccessHint)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// UpdateBoolConfigHandler 处理 Bool 类型
func UpdateBoolConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	boolVal, parseErr := strconv.ParseBool(newValue)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 bool: %v", newValue, parseErr)
		return echo.NewHTTPError(http.StatusBadRequest, "parse bool error")
	}
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" || name == "Debug" {
		saveSuccessHint = true
	}
	updatedHTML := BoolConfig(name, boolVal, name+"_Description", saveSuccessHint)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateNumberConfigHandler 处理 Number 类型
func UpdateNumberConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
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
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
		saveSuccessHint = true
	}
	// 渲染对应的模板
	updatedHTML := NumberConfig(name, int(intVal), name+"_Description", 0, 65535, saveSuccessHint)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// UpdateLoginSettingsHandler 登录相关设置
func UpdateLoginSettingsHandler(c echo.Context) error {
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
	username := c.FormValue("Username")
	currentPassword := c.FormValue("CurrentPassword")
	password := c.FormValue("Password") // 新密码(初次设定) 或 原始密码（已有密码时）
	reEnterPassword := c.FormValue("ReEnterPassword")
	// 除非是调试模式, 密码不明文记录到日志，
	if state.ServerConfig.Debug {
		logger.Infof("Update user info: Username=%s", username)
		logger.Infof("Update user info: CurrentPassword=%s", currentPassword)
		logger.Infof("Update user info: Password=%s", password)
		logger.Infof("Update user info: ReEnterPassword=%s", reEnterPassword) // ReEnterPassword
	}

	// 两次输入的密码不一致
	if password != reEnterPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Password do not match")
	}
	// 用户名或密码为空
	if username == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and Password cannot be empty")
	}
	//// 密码过短
	//if len(password) < 6 {
	//	return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters long")
	//}
	//// 用户名过短
	//if len(username) < 3 {
	//	return echo.NewHTTPError(http.StatusBadRequest, "Username must be at least 3 characters long")
	//}

	if state.ServerConfig.Password != "" {
		// 当前密码不正确
		if state.ServerConfig.Password != currentPassword {
			return echo.NewHTTPError(http.StatusBadRequest, "Current Password is incorrect")
		}
	}

	// 旧配置做个备份（后面需要对比）
	oldConfig := config.CopyCfg()

	// 更新用户名
	if err := state.ServerConfig.SetConfigValue("Username", username); err != nil {
		logger.Errorf("Failed to set Username: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update username")
	}
	// 更新密码
	if err := state.ServerConfig.SetConfigValue("Password", password); err != nil {
		// logger.Errorf("Failed to set Password: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
		// 这里可能需要回滚配置更改或返回错误
		// return echo.NewHTTPError(http.StatusInternalServerError, "Failed to write config file")
	}

	// 根据配置的变化，做相应操作。
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	// fmt.Printf("New config: %+v\n", config.GetCfg())

	// 渲染 UserInfoConfig 模板并返回
	updatedHTML := UserInfoConfig(config.GetCfg().Username, config.GetCfg().Password, true) // 如果UserInfoConfig期望的是加密后的密码，这里需要调整
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		logger.Errorf("Failed to render UserInfoConfig template: %v", renderErr)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return nil
}

// UpdateTailscaleConfigHandler 处理Tailscale配置更新的JSON API
func UpdateTailscaleConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}

	// 解析请求体（JSON 或 x-www-form-urlencoded）
	var request struct {
		EnableTailscale     bool   `json:"EnableTailscale"`
		TailscaleAuthKey    string `json:"TailscaleAuthKey"`
		TailscaleHostname   string `json:"TailscaleHostname"`
		TailscalePort       int    `json:"TailscalePort"`
		TailscaleFunnelMode bool   `json:"TailscaleFunnelMode"`
	}
	contentType := c.Request().Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
		}
	} else {
		// 表单解析（HTMX fetch 使用 x-www-form-urlencoded）
		if v := c.FormValue("EnableTailscale"); v != "" {
			if b, perr := strconv.ParseBool(v); perr == nil {
				request.EnableTailscale = b
			}
		}
		if v := c.FormValue("TailscaleAuthKey"); v != "" {
			request.TailscaleAuthKey = v
		}
		if v := c.FormValue("TailscaleHostname"); v != "" {
			request.TailscaleHostname = v
		}
		if v := c.FormValue("TailscalePort"); v != "" {
			if p, perr := strconv.Atoi(v); perr == nil {
				request.TailscalePort = p
			}
		}
		if v := c.FormValue("TailscaleFunnelMode"); v != "" {
			if b, perr := strconv.ParseBool(v); perr == nil {
				request.TailscaleFunnelMode = b
			}
		}
		// 对缺失字段使用现有配置作为默认值，避免零值覆盖
		if request.TailscaleHostname == "" {
			request.TailscaleHostname = config.GetTailscaleHostname()
		}
		if request.TailscalePort == 0 {
			request.TailscalePort = config.GetTailscalePort()
		}
	}

	// 验证输入
	if request.EnableTailscale {
		if request.TailscaleHostname == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "TailscaleHostname cannot be empty when Tailscale is enabled")
		}
		if request.TailscalePort < 0 || request.TailscalePort > 65535 {
			return echo.NewHTTPError(http.StatusBadRequest, "TailscalePort must be between 0 and 65535")
		}
		if request.TailscaleFunnelMode && (request.TailscalePort != 443 && request.TailscalePort != 8443 && request.TailscalePort != 10000) {
			return echo.NewHTTPError(http.StatusBadRequest, "Port must be 443, 8443, or 10000 when Funnel Mode is enabled")
		}
	}
	// 旧配置做个备份（后面需要对比）
	oldConfig := config.CopyCfg()
	// 更新Tailscale配置
	config.GetCfg().EnableTailscale = request.EnableTailscale
	config.GetCfg().TailscaleAuthKey = request.TailscaleAuthKey
	config.GetCfg().TailscaleHostname = request.TailscaleHostname
	config.GetCfg().TailscalePort = request.TailscalePort
	config.GetCfg().TailscaleFunnelMode = request.TailscaleFunnelMode
	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}

	// 根据配置的变化，做相应操作
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	// 返回成功响应
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Tailscale configuration updated successfully",
	})
}

func AddArrayConfigHandler(c echo.Context) error {
	// 解析htmx的form数据
	// 原始数据类似：configName=SupportMediaType&addValue=.test
	// 解析后的数据类似：ConfigName=SupportMediaType, AddValue=.test
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
	configName := c.FormValue("configName")
	addValue := c.FormValue("addValue")

	logger.Infof("AddArrayConfigHandler: %s = %s\n", configName, addValue)

	values, err := doAdd(configName, addValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "add error")
	}
	saveSuccessHint := false
	if configName == "StoreUrls" {
		saveSuccessHint = true
	}
	updatedHTML := StringArrayConfig(configName, values, configName+"_Description", saveSuccessHint)
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
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
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
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}

	configName := c.FormValue("configName")
	deleteValue := c.FormValue("deleteValue")

	logger.Infof("DeleteArrayConfigHandler: %s = %s\n", configName, deleteValue)

	values, err := doDelete(configName, deleteValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, configName+" delete Failed")
	}

	updatedHTML := StringArrayConfig(configName, values, configName+"_Description", false)
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
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
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
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
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
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
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
