package settings

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// -------------------------
// 抽取更新配置的公共逻辑
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}
	// 调用通用逻辑更新配置
	name, newValue, err := updateStringConfigFromJSON(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
		saveSuccessHint = true
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"message":         "Configuration updated successfully",
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	name, boolVal, err := updateBoolConfigFromJSON(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" || name == "Debug" {
		saveSuccessHint = true
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"message":         "Configuration updated successfully",
		"value":           boolVal,
		"saveSuccessHint": saveSuccessHint,
	})
}

// updateNumberConfigFromJSON 从 JSON 请求中更新数字配置的通用逻辑
func updateNumberConfigFromJSON(c echo.Context) (string, int, error) {
	// 解析 JSON 请求体
	var request struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	if err := c.Bind(&request); err != nil {
		return "", 0, fmt.Errorf("invalid JSON request: %v", err)
	}

	if request.Name == "" {
		return "", 0, errors.New("name is required")
	}

	// 将 int 转换为 string
	newValue := strconv.Itoa(request.Value)

	logger.Infof(locale.GetString("log_update_config"), request.Name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", 0, setErr
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())
	return request.Name, request.Value, nil
}

// UpdateNumberConfigHandler 处理 Number 类型的配置
func UpdateNumberConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}
	// 调用通用逻辑更新配置
	name, intVal, err := updateNumberConfigFromJSON(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint := false
	if name == "Username" || name == "Password" || name == "Port" || name == "Host" || name == "DisableLAN" || name == "Timeout" {
		saveSuccessHint = true
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"message":         "Configuration updated successfully",
		"value":           intVal,
		"saveSuccessHint": saveSuccessHint,
	})
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	if err := updateLoginSettingsFromJSON(c); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login settings updated successfully",
	})
}

// UpdateTailscaleConfigHandler 处理Tailscale配置更新的JSON API
func UpdateTailscaleConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, "Config is locked, cannot be modified")
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
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Tailscale configuration updated successfully",
	})
}

// AddArrayConfigHandler 处理添加数组元素的 JSON API
func AddArrayConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName string `json:"configName"`
		AddValue   string `json:"addValue"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON request",
		})
	}

	if request.ConfigName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "configName is required",
		})
	}

	logger.Infof(locale.GetString("log_add_array_config_handler")+"\n", request.ConfigName, request.AddValue)

	values, err := doAdd(request.ConfigName, request.AddValue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to add config value",
		})
	}

	saveSuccessHint := false
	if request.ConfigName == "StoreUrls" {
		saveSuccessHint = true
	}

	// 渲染更新后的 HTML
	updatedHTML := StringArrayConfig(request.ConfigName, values, request.ConfigName+"_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to render template",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"message":         "Configuration value added successfully",
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

// DeleteArrayConfigHandler 处理删除数组元素的 JSON API
func DeleteArrayConfigHandler(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName  string `json:"configName"`
		DeleteValue string `json:"deleteValue"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON request",
		})
	}

	if request.ConfigName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "configName is required",
		})
	}

	if request.DeleteValue == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "deleteValue is required",
		})
	}

	logger.Infof(locale.GetString("log_delete_array_config_handler")+"\n", request.ConfigName, request.DeleteValue)

	values, err := doDelete(request.ConfigName, request.DeleteValue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to delete config value",
		})
	}

	// 渲染更新后的 HTML
	updatedHTML := StringArrayConfig(request.ConfigName, values, request.ConfigName+"_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to render template",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Configuration value deleted successfully",
		"html":    htmlString,
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON request",
		})
	}

	if request.SelectedDir == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No directory selected",
		})
	}

	if request.SelectedDir != config.WorkingDirectory && request.SelectedDir != config.HomeDirectory && request.SelectedDir != config.ProgramDirectory {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid directory selected",
		})
	}

	// 保存配置
	if err := config.SaveConfig(request.SelectedDir); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to save config",
		})
	}

	// 渲染更新后的 HTML
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to render template",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Configuration saved successfully",
		"html":    htmlString,
	})
}

// HandleConfigDelete 处理 /api/config-delete 的 JSON API
func HandleConfigDelete(c echo.Context) error {
	// 如果配置被锁定
	if config.GetCfg().ReadOnlyMode {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Config is locked, cannot be modified",
		})
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON request",
		})
	}

	if request.SelectedDir == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No directory selected",
		})
	}

	if request.SelectedDir != config.WorkingDirectory && request.SelectedDir != config.HomeDirectory && request.SelectedDir != config.ProgramDirectory {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid directory selected",
		})
	}

	// 删除配置
	if err := config.DeleteConfigIn(request.SelectedDir); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to delete config",
		})
	}

	// 渲染更新后的 HTML
	updatedHTML := ConfigManager(config.DefaultConfigLocation(), config.GetWorkingDirectoryConfig(), config.GetHomeDirectoryConfig(), config.GetProgramDirectoryConfig())
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to render template",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Configuration deleted successfully",
		"html":    htmlString,
	})
}
