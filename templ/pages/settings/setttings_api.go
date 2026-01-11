package settings

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// -------------------------
// 更新配置的公共逻辑
// -------------------------

// parseSingleHTMXFormPair 提取并返回表单中的"第一对" key/value
func parseSingleHTMXFormPair(c echo.Context) (string, string, error) {
	if !htmx.IsHTMX(c.Request()) {
		return "", "", errors.New(locale.GetString("err_non_htmx_request"))
	}
	formData, err := c.FormParams()
	if err != nil {
		return "", "", fmt.Errorf("parseForm error: %v", err)
	}
	if len(formData) == 0 {
		return "", "", errors.New(locale.GetString("err_no_form_data"))
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

	logger.Infof(locale.GetString("log_update_config"), name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", "", setErr
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return name, newValue, nil
}

// -------------------------
// 各类配置的更新 PageHandler
// -------------------------

// updateStringConfigFromJSON 从 JSON 请求中更新字符串配置的通用逻辑
func updateStringConfigFromJSON(c echo.Context) (string, string, error) {
	// 解析 JSON 请求体
	var request struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err := c.Bind(&request); err != nil {
		return "", "", fmt.Errorf("invalid JSON request: %v", err)
	}

	if request.Name == "" {
		return "", "", errors.New("name is required")
	}

	logger.Infof(locale.GetString("log_update_config"), request.Name, request.Value)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, request.Value); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", "", setErr
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return request.Name, request.Value, nil
}

// UpdateStringConfigHandler 处理 String 类型的 JSON API
func UpdateStringConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}
	// 调用通用逻辑更新配置
	name, newValue, err := updateStringConfigFromJSON(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
		saveSuccessHint = true
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"value":           newValue,
		"saveSuccessHint": saveSuccessHint,
	})
}

// updateBoolConfigFromJSON 从 JSON 请求中更新布尔配置的通用逻辑
func updateBoolConfigFromJSON(c echo.Context) (string, bool, error) {
	// 解析 JSON 请求体
	var request struct {
		Name  string `json:"name"`
		Value bool   `json:"value"`
	}
	if err := c.Bind(&request); err != nil {
		return "", false, fmt.Errorf("invalid JSON request: %v", err)
	}

	if request.Name == "" {
		return "", false, errors.New("name is required")
	}

	// 将 bool 转换为 string
	newValue := strconv.FormatBool(request.Value)

	logger.Infof(locale.GetString("log_update_config"), request.Name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", false, setErr
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return request.Name, request.Value, nil
}

// UpdateBoolConfigHandler 处理 Bool 类型的 JSON API
func UpdateBoolConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	name, boolVal, err := updateBoolConfigFromJSON(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" || name == "Debug" {
		saveSuccessHint = true
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"value":           boolVal,
		"saveSuccessHint": saveSuccessHint,
	})
}

// updateNumberConfigFromJSON 从 JSON 请求中更新数字配置的通用逻辑
func updateNumberConfigFromJSON(c echo.Context) (string, int, *config.Config, error) {
	// 解析 JSON 请求体
	var request struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	if err := c.Bind(&request); err != nil {
		return "", 0, nil, fmt.Errorf("invalid JSON request: %v", err)
	}

	if request.Name == "" {
		return "", 0, nil, errors.New("name is required")
	}

	// 将 int 转换为 string
	newValue := strconv.Itoa(request.Value)

	logger.Infof(locale.GetString("log_update_config"), request.Name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", 0, nil, setErr
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	return request.Name, request.Value, &oldConfig, nil
}

// UpdateNumberConfigHandler 处理 Number 类型的配置
func UpdateNumberConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}
	// 调用通用逻辑更新配置
	name, intVal, oldConfig, err := updateNumberConfigFromJSON(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	if name == "Port" {
		go func() {
			// 延迟1秒执行
			time.Sleep(1 * time.Second)
			beforeConfigUpdate(oldConfig, config.GetCfg())
		}()
	} else {
		beforeConfigUpdate(oldConfig, config.GetCfg())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" || name == "AutoRescanIntervalMinutes" {
		saveSuccessHint = true
	}

	// 构建响应
	response := map[string]interface{}{
		"value":           intVal,
		"saveSuccessHint": saveSuccessHint,
	}
	return c.JSON(http.StatusOK, response)
}

// updateLoginSettingsFromJSON 从 JSON 请求中更新登录设置的通用逻辑
func updateLoginSettingsFromJSON(c echo.Context) error {
	// 解析 JSON 请求体
	var request struct {
		Username        string `json:"username"`
		CurrentPassword string `json:"currentPassword"`
		Password        string `json:"password"`
		ReEnterPassword string `json:"reEnterPassword"`
	}
	if err := c.Bind(&request); err != nil {
		return fmt.Errorf("invalid JSON request: %v", err)
	}

	// 除非是调试模式, 密码不明文记录到日志
	if config.GetCfg().Debug {
		logger.Infof(locale.GetString("log_update_user_info_username"), request.Username)
		logger.Infof(locale.GetString("log_update_user_info_current_password"), request.CurrentPassword)
		logger.Infof(locale.GetString("log_update_user_info_password"), request.Password)
		logger.Infof(locale.GetString("log_update_user_info_reenter_password"), request.ReEnterPassword)
	}

	// 两次输入的密码不一致
	if request.Password != request.ReEnterPassword {
		return errors.New("Password do not match")
	}
	// 用户名或密码为空
	if request.Username == "" || request.Password == "" {
		return errors.New("Username and Password cannot be empty")
	}

	// 当前密码不正确（如果已有密码）
	if config.GetCfg().Password != "" && config.GetCfg().Password != request.CurrentPassword {
		return errors.New("Current Password is incorrect")
	}

	// 旧配置做个备份（后面需要对比）
	oldConfig := config.CopyCfg()

	// 更新用户名
	if err := config.GetCfg().SetConfigValue("Username", request.Username); err != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_username"), err)
		return fmt.Errorf("failed to update username: %v", err)
	}
	// 更新密码
	if err := config.GetCfg().SetConfigValue("Password", request.Password); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	// 根据配置的变化，做相应操作
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return nil
}

// UpdateLoginSettingsHandler 处理登录设置的 JSON API
func UpdateLoginSettingsHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	if err := updateLoginSettingsFromJSON(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// UpdateTailscaleConfigHandler 处理Tailscale配置更新的JSON API
func UpdateTailscaleConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}
	// 解析请求体（JSON格式）
	var request struct {
		EnableTailscale   bool   `json:"EnableTailscale"`
		TailscaleAuthKey  string `json:"TailscaleAuthKey"`
		TailscaleHostname string `json:"TailscaleHostname"`
		TailscalePort     int    `json:"TailscalePort"`
		FunnelTunnel      bool   `json:"FunnelTunnel"`
		FunnelLoginCheck  bool   `json:"FunnelLoginCheck"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}
	fmt.Printf("Received Tailscale config update: %+v\n", request)

	// 验证输入
	if request.EnableTailscale {
		if request.TailscaleHostname == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "TailscaleHostname cannot be empty when Tailscale is enabled")
		}
		if request.TailscalePort < 0 || request.TailscalePort > 65535 {
			return echo.NewHTTPError(http.StatusBadRequest, "TailscalePort must be between 0 and 65535")
		}
		if request.FunnelTunnel && (request.TailscalePort != 443 && request.TailscalePort != 8443 && request.TailscalePort != 10000) {
			return echo.NewHTTPError(http.StatusBadRequest, "Port must be 443, 8443, or 10000 when Funnel Mode is enabled")
		}
		// 如果Funnel模式强制密码，但当前没有设置密码，则返回错误
		if request.FunnelTunnel && request.FunnelLoginCheck && !config.GetCfg().RequiresAuth() {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot Turn on FunnelMode when no password not set")
		}
	}
	// 旧配置做个备份（后面需要对比）
	oldConfig := config.CopyCfg()
	// 更新Tailscale配置
	config.GetCfg().EnableTailscale = request.EnableTailscale
	config.GetCfg().TailscaleAuthKey = request.TailscaleAuthKey
	config.GetCfg().TailscaleHostname = request.TailscaleHostname
	config.GetCfg().TailscalePort = request.TailscalePort
	config.GetCfg().FunnelTunnel = request.FunnelTunnel
	config.GetCfg().FunnelLoginCheck = request.FunnelLoginCheck
	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	// 根据配置的变化，做相应操作
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	// 返回成功响应
	return c.NoContent(http.StatusOK)
}

// AddArrayConfigHandler 处理添加数组元素的 JSON API
func AddArrayConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName string `json:"configName"`
		AddValue   string `json:"addValue"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.ConfigName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is required")
	}

	logger.Infof(locale.GetString("log_add_array_config_handler")+"\n", request.ConfigName, request.AddValue)

	values, err := doAdd(request.ConfigName, request.AddValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_add_config_failed"))
	}

	saveSuccessHint := false
	if request.ConfigName == "StoreUrls" {
		saveSuccessHint = true
	}

	// 渲染更新后的 HTML
	updatedHTML := StringArrayConfig(request.ConfigName, values, request.ConfigName+"_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"saveSuccessHint": saveSuccessHint,
		"html":            htmlString,
	})
}

func doAdd(configName, addValue string) ([]string, error) {
	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	values, err := config.GetCfg().AddStringArrayConfig(configName, addValue)
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_add_config_value"), err)
		return nil, err
	}
	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	return values, nil
}

// EnablePluginHandler 处理启用插件的 JSON API
func EnablePluginHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		PluginName string `json:"pluginName"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.PluginName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pluginName is required")
	}

	logger.Infof("启用插件: %s\n", request.PluginName)

	// 启用插件
	err := config.GetCfg().AddPlugin(request.PluginName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_update_config_failed"))
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "插件已启用",
	})
}

// DisablePluginHandler 处理禁用插件的 JSON API
func DisablePluginHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		PluginName string `json:"pluginName"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.PluginName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pluginName is required")
	}

	logger.Infof("禁用插件: %s\n", request.PluginName)

	// 禁用插件
	err := config.GetCfg().DisablePlugin(request.PluginName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_update_config_failed"))
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "插件已禁用",
	})
}

// DeleteArrayConfigHandler 处理删除数组元素的 JSON API
func DeleteArrayConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName  string `json:"configName"`
		DeleteValue string `json:"deleteValue"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.ConfigName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is required")
	}

	if request.DeleteValue == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "deleteValue is required")
	}

	logger.Infof(locale.GetString("log_delete_array_config_handler")+"\n", request.ConfigName, request.DeleteValue)

	values, err := doDelete(request.ConfigName, request.DeleteValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_config_failed"))
	}

	// 渲染更新后的 HTML
	updatedHTML := StringArrayConfig(request.ConfigName, values, request.ConfigName+"_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"html": htmlString,
	})
}

func doDelete(configName string, deleteValue string) ([]string, error) {
	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	values, err := config.GetCfg().DeleteStringArrayConfig(configName, deleteValue)
	if err != nil {
		return nil, err
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	return values, nil
}

// renderTemplToString 将 templ 组件渲染为 HTML 字符串
func renderTemplToString(ctx context.Context, component templ.Component) (string, error) {
	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// HandleConfigSave 处理 /api/config-save 的 JSON API
func HandleConfigSave(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.SelectedDir == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No directory selected")
	}

	if request.SelectedDir != config.WorkingDirectory && request.SelectedDir != config.HomeDirectory && request.SelectedDir != config.ProgramDirectory {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid directory selected")
	}

	// 保存配置
	if err := config.SaveConfig(request.SelectedDir); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_save_config_failed"))
	}

	// 渲染更新后的 HTML
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"html": htmlString,
	})
}

// HandleConfigDelete 处理 /api/config-delete 的 JSON API
func HandleConfigDelete(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.SelectedDir == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No directory selected")
	}

	if request.SelectedDir != config.WorkingDirectory && request.SelectedDir != config.HomeDirectory && request.SelectedDir != config.ProgramDirectory {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid directory selected")
	}

	// 删除配置
	if err := config.DeleteConfigIn(request.SelectedDir); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_config_failed"))
	}

	// 渲染更新后的 HTML
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"html": htmlString,
	})
}

// RescanStoreHandler 处理重新扫描书库的 JSON API
func RescanStoreHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		StoreUrl string `json:"storeUrl"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.StoreUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is required")
	}

	logger.Infof("重新扫描书库: %s\n", request.StoreUrl)

	// 记录扫描前的书籍数量
	beforeCount := model.GetAllBooksNumber()

	// 调用扫描功能
	err := scan.InitStore(request.StoreUrl, config.GetCfg())
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_scan_store_path"), err)
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_rescan_store_failed"))
	}

	// 如果启用数据库，保存扫描结果
	if config.GetCfg().EnableDatabase {
		if err := scan.SaveBooksToDatabase(config.GetCfg()); err != nil {
			logger.Infof(locale.GetString("log_failed_to_save_results_to_database"), err)
			return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_rescan_store_failed"))
		}
	}

	// 计算新增的书籍数量
	afterCount := model.GetAllBooksNumber()
	newBooksCount := afterCount - beforeCount
	if newBooksCount < 0 {
		newBooksCount = 0
	}

	logger.Infof("书库扫描完成，新增 %d 本书\n", newBooksCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"newBooksCount": newBooksCount,
		"message":       locale.GetString("rescan_store_success"),
	})
}

// DeleteStoreHandler 处理删除书库的 JSON API
func DeleteStoreHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}

	// 解析 JSON 请求体
	var request struct {
		StoreUrl string `json:"storeUrl"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	if request.StoreUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is required")
	}

	logger.Infof("删除书库: %s\n", request.StoreUrl)

	// 先删除该书库的所有书籍数据
	targetStoreAbs, err := filepath.Abs(request.StoreUrl)
	if err != nil {
		logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
		targetStoreAbs = request.StoreUrl
	}

	// 遍历所有书籍，删除属于该书库的书籍
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_store_failed"))
	}

	deletedCount := 0
	for _, book := range allBooks {
		bookStoreAbs, err := filepath.Abs(book.StoreUrl)
		if err != nil {
			logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
			bookStoreAbs = book.StoreUrl
		}

		if bookStoreAbs == targetStoreAbs {
			err := model.IStore.DeleteBook(book.BookID)
			if err != nil {
				logger.Infof(locale.GetString("log_error_deleting_book"), book.BookID, err)
			} else {
				deletedCount++
			}
		}
	}

	logger.Infof("删除了 %d 本书籍\n", deletedCount)

	// 从配置中移除该书库 URL
	values, err := doDelete("StoreUrls", request.StoreUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_store_failed"))
	}

	// 重新生成书组
	if err := model.IStore.GenerateBookGroup(); err != nil {
		logger.Infof(locale.GetString("log_error_initializing_main_folder"), err)
	}

	// 渲染更新后的 HTML
	updatedHTML := StoreConfig("StoreUrls", values, "StoreUrls_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"html":    htmlString,
		"message": locale.GetString("delete_store_success"),
	})
}
