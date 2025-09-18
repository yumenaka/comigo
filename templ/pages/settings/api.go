package settings

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/templ/state"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/tailscale_plugin"
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

// GetTailscaleStatus 处理Tailscale网络配置
func GetTailscaleStatus(c echo.Context) error {
	tailscaleStatus, err := tailscale_plugin.GetTailscaleStatus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tailscaleStatus)
}

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
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
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

// UpdateLoginSettingsHandler 更新登录相关信息
func UpdateLoginSettingsHandler(c echo.Context) error {
	if !htmx.IsHTMX(c.Request()) {
		return echo.NewHTTPError(http.StatusBadRequest, "non-htmx request")
	}
	// 如果配置被锁定
	if config.GetCfg().GetConfigLocked() {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
	}
	username := c.FormValue("Username")
	oldPassword := c.FormValue("OldPassword")
	password := c.FormValue("Password") // 注意：这里获取的是用户在表单中输入的新密码或原始密码
	reEnterPassword := c.FormValue("ReEnterPassword")
	// 密码通常不明文记录到日志，除非是调试模式
	if state.ServerConfig.Debug {
		logger.Infof("Update user info: Username=%s", username)
		logger.Infof("Update user info: OldPassword=%s", oldPassword)
		logger.Infof("Update user info: Password=%s", password)
		logger.Infof("Update user info: ReEnterPassword=%s", reEnterPassword) //ReEnterPassword
	}

	// 两次输入的密码不一致
	if password != reEnterPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "Passwords do not match")
	}
	// 用户名或密码为空
	if username == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and Password cannot be empty")
	}
	// 密码过短
	if len(password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 6 characters long")
	}
	// 用户名过短
	if len(username) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username must be at least 3 characters long")
	}

	if state.ServerConfig.Password != "" {
		// 旧密码不正确
		if state.ServerConfig.Password != oldPassword {
			return echo.NewHTTPError(http.StatusBadRequest, "Old password is incorrect")
		}
	}

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if err := state.ServerConfig.SetConfigValue("Username", username); err != nil {
		logger.Errorf("Failed to set Username: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update username")
	}
	// 只有当用户在表单中确实输入了新密码时才更新密码
	// HTMX表单提交时，如果密码字段为空，Password会是空字符串。
	// 根据 UserInfoConfig 的逻辑，如果密码字段为空，可能表示不修改密码，或者用户清空了密码。
	// 这里的逻辑假设：如果Password字段有值（用户输入了），则更新它。
	// 如果允许清空密码，这里的逻辑是正确的。
	// 如果密码字段不允许为空，前端应该有校验。
	// 或者，可以比较提交的密码是否与旧密码（加密/哈希后）相同，如果相同则不更新。
	// 但此处简单处理：只要提交了Password，就更新。
	if err := state.ServerConfig.SetConfigValue("Password", password); err != nil {
		logger.Errorf("Failed to set Password: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
		// 根据业务需求，这里可能需要回滚配置更改或返回错误
		// return echo.NewHTTPError(http.StatusInternalServerError, "Failed to write config file")
	}

	// 根据配置的变化，做相应操作。
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	// 渲染 UserInfoConfig 模板并返回
	// 注意：UserInfoConfig 模板期望的是 initPassword，这里传递的是用户表单提交的 password
	// 如果密码在配置中是加密存储的，而 UserInfoConfig 期望的是明文（或特定格式），这里需要适配
	// 假设 UserInfoConfig 可以直接使用表单提交的 username 和 password 初始化
	updatedHTML := UserInfoConfig(username, password, true) // 如果UserInfoConfig期望的是加密后的密码，这里需要调整
	if renderErr := htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, updatedHTML); renderErr != nil {
		logger.Errorf("Failed to render UserInfoConfig template: %v", renderErr)
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
