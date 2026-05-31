package settings

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// decodeBase64URLStrict 将 base64url（RawURLEncoding，无 padding）解码为原始字符串。
// 解码失败应视为客户端请求参数不合法。
func decodeBase64URLStrict(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty base64url value")
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func reloadHintForStringConfig(name string) (saveSuccessHint bool, reason string) {
	switch name {
	case "Host", "BasePath":
		return true, sse_hub.UISuggestReasonServerConfig
	default:
		return false, ""
	}
}

func reloadHintForBoolConfig(name string) (saveSuccessHint bool, reason string) {
	switch name {
	case "DisableLAN":
		return true, sse_hub.UISuggestReasonServerConfig
	case "Debug":
		return true, sse_hub.UISuggestReasonDebugToggle
	case "EnablePlugin":
		return true, sse_hub.UISuggestReasonPluginsChanged
	default:
		return false, ""
	}
}

func reloadHintForNumberConfig(name string) (saveSuccessHint bool, reason string) {
	switch name {
	case "Port":
		// 端口变更由前端自行跳转，不再弹刷新确认。
		return true, ""
	case "Timeout":
		return true, sse_hub.UISuggestReasonServerConfig
	case "AutoRescanIntervalMinutes":
		// 自动重扫仅调整调度器，不需要整页刷新。
		return true, ""
	default:
		return false, ""
	}
}

func renderArrayConfigComponent(configName string, values []string) templ.Component {
	if configName == "StoreUrls" {
		return StoreConfig(configName, values, configName+"_Description", GetStoreBookCounts())
	}
	return StringArrayConfig(configName, values, configName+"_Description")
}

// ensureWritableConfig 统一拦截只读模式下的设置写操作。
func ensureWritableConfig() error {
	if config.GetCfg().ReadOnlyMode {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_config_locked"))
	}
	return nil
}

// jsonBadRequest 统一 JSON 解析失败响应，避免每个 handler 重复拼错误。
func jsonBadRequest(err error) error {
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, locale.GetString("err_invalid_json_request"))
	}
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s: %v", locale.GetString("err_invalid_json_request"), err))
}

// writeConfigAndApply 写入配置后应用变更副作用；写入失败沿用旧行为，只记录日志不中断响应。
func writeConfigAndApply(oldConfig config.Config) {
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	beforeConfigUpdate(oldConfig, config.GetCfg())
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

	logger.Infof(locale.GetString("log_update_config"), request.Name)

	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, request.Value); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", "", setErr
	}

	writeConfigAndApply(oldConfig)

	return request.Name, request.Value, nil
}

// UpdateStringConfigHandler 处理 String 类型的 JSON API
func UpdateStringConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}
	// 调用通用逻辑更新配置
	name, newValue, err := updateStringConfigFromJSON(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint, reloadReason := reloadHintForStringConfig(name)
	if reloadReason != "" {
		sse_hub.BroadcastUISuggestReload(reloadReason)
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

	logger.Infof(locale.GetString("log_update_config"), request.Name)

	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", false, setErr
	}

	writeConfigAndApply(oldConfig)

	return request.Name, request.Value, nil
}

// UpdateBoolConfigHandler 处理 Bool 类型的 JSON API
func UpdateBoolConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	name, boolVal, err := updateBoolConfigFromJSON(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint, reloadReason := reloadHintForBoolConfig(name)
	if reloadReason != "" {
		sse_hub.BroadcastUISuggestReload(reloadReason)
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

	logger.Infof(locale.GetString("log_update_config"), request.Name)

	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新配置
	if setErr := config.GetCfg().SetConfigValue(request.Name, newValue); setErr != nil {
		logger.Errorf(locale.GetString("err_failed_to_set_config_value"), setErr)
		return "", 0, nil, setErr
	}

	// 端口变更需要延迟应用副作用，因此这里只写文件，不走 writeConfigAndApply。
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}
	return request.Name, request.Value, &oldConfig, nil
}

// UpdateNumberConfigHandler 处理 Number 类型的配置
func UpdateNumberConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
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
			beforeConfigUpdate(*oldConfig, config.GetCfg())
		}()
	} else {
		beforeConfigUpdate(*oldConfig, config.GetCfg())
	}

	// 判断是否需要显示保存成功提示并刷新页面
	saveSuccessHint, reloadReason := reloadHintForNumberConfig(name)
	if reloadReason != "" {
		sse_hub.BroadcastUISuggestReload(reloadReason)
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

	request.Username = strings.TrimSpace(request.Username)
	request.CurrentPassword = strings.TrimSpace(request.CurrentPassword)

	cfg := config.GetCfg()
	existingPasswordLogin := cfg.HasPasswordLoginConfigured()
	passwordLoginChanged := request.Username != cfg.Username ||
		request.Password != "" ||
		request.ReEnterPassword != ""

	// 仅在修改账号密码登录设置时校验当前密码。
	if existingPasswordLogin && passwordLoginChanged && cfg.Password != request.CurrentPassword {
		return errors.New(locale.GetString("err_current_password_incorrect"))
	}

	if request.Password != request.ReEnterPassword {
		return errors.New(locale.GetString("err_password_mismatch"))
	}

	var effectiveUsername string
	var effectivePassword string
	if request.Username == "" {
		if request.Password != "" || request.ReEnterPassword != "" {
			return errors.New(locale.GetString("prompt_set_username"))
		}
		effectiveUsername = ""
		effectivePassword = ""
	} else {
		effectiveUsername = request.Username
		switch {
		case request.Password != "":
			effectivePassword = request.Password
		case existingPasswordLogin:
			effectivePassword = cfg.Password
		default:
			return errors.New(locale.GetString("prompt_set_password"))
		}
	}

	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	cfg.Username = effectiveUsername
	cfg.Password = effectivePassword

	writeConfigAndApply(oldConfig)

	return nil
}

// UpdateLoginSettingsHandler 处理登录设置的 JSON API
func UpdateLoginSettingsHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	if err := updateLoginSettingsFromJSON(c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sse_hub.BroadcastUISuggestReload(sse_hub.UISuggestReasonLoginSettings)
	return c.NoContent(http.StatusOK)
}

// UpdateTailscaleConfigHandler 处理Tailscale配置更新的JSON API
func UpdateTailscaleConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
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
		return jsonBadRequest(err)
	}
	//fmt.Printf("Received Tailscale config update: %+v\n", request)

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
		// 如果 Funnel 模式要求登录保护，但当前没有可用登录方式，则返回错误
		if request.FunnelTunnel && request.FunnelLoginCheck && !config.GetCfg().RequiresAuth() {
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot turn on FunnelMode when login protection is unavailable")
		}
	}
	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新Tailscale配置
	config.GetCfg().EnableTailscale = request.EnableTailscale
	config.GetCfg().TailscaleAuthKey = request.TailscaleAuthKey
	config.GetCfg().TailscaleHostname = request.TailscaleHostname
	config.GetCfg().TailscalePort = request.TailscalePort
	config.GetCfg().FunnelTunnel = request.FunnelTunnel
	config.GetCfg().FunnelLoginCheck = request.FunnelLoginCheck
	writeConfigAndApply(oldConfig)

	// 返回成功响应
	return c.NoContent(http.StatusOK)
}

// AddArrayConfigHandler 处理添加数组元素的 JSON API
func AddArrayConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName string `json:"configName"`
		AddValue   string `json:"addValue"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.ConfigName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is required")
	}

	decodedConfigName, err := decodeBase64URLStrict(request.ConfigName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is not valid base64url")
	}
	decodedAddValue, err := decodeBase64URLStrict(request.AddValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "addValue is not valid base64url")
	}
	logger.Infof(locale.GetString("log_add_array_config_handler")+"\n", decodedConfigName)

	values, err := doAdd(decodedConfigName, decodedAddValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_add_config_failed"))
	}

	saveSuccessHint := false
	if decodedConfigName == "StoreUrls" {
		saveSuccessHint = true
	}

	// 渲染更新后的 HTML。StoreUrls 需要回传专用组件，保证空书架/设置页都能即时更新列表。
	updatedHTML := renderArrayConfigComponent(decodedConfigName, values)
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
	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新配置
	values, err := config.GetCfg().AddStringArrayConfig(configName, addValue)
	if err != nil {
		logger.Errorf(locale.GetString("err_failed_to_add_config_value"), err)
		return nil, err
	}
	writeConfigAndApply(oldConfig)
	return values, nil
}

// EnablePluginHandler 处理启用插件的 JSON API
func EnablePluginHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		PluginName string `json:"pluginName"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.PluginName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pluginName is required")
	}

	logger.Infof(locale.GetString("log_plugin_enabled")+"\n", request.PluginName)

	// 互斥逻辑：auto_flip 和 sketch_practice 不能同时启用
	if request.PluginName == "sketch_practice" && config.GetCfg().IsPluginEnabled("auto_flip") {
		// 启用 sketch_practice 时，禁用 auto_flip
		logger.Infof(locale.GetString("log_disable_mutex_plugin_auto_flip") + "\n")
		_ = config.GetCfg().DisablePlugin("auto_flip")
	} else if request.PluginName == "auto_flip" && config.GetCfg().IsPluginEnabled("sketch_practice") {
		// 启用 auto_flip 时，禁用 sketch_practice
		logger.Infof(locale.GetString("log_disable_mutex_plugin_sketch_practice") + "\n")
		_ = config.GetCfg().DisablePlugin("sketch_practice")
	}

	// 启用插件
	err := config.GetCfg().AddPlugin(request.PluginName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_update_config_failed"))
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	sse_hub.BroadcastUISuggestReload(sse_hub.UISuggestReasonPluginsChanged)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":         true,
		"message":         "插件已启用",
		"saveSuccessHint": true,
	})
}

// DisablePluginHandler 处理禁用插件的 JSON API
func DisablePluginHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		PluginName string `json:"pluginName"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.PluginName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "pluginName is required")
	}

	logger.Infof(locale.GetString("log_plugin_disabled")+"\n", request.PluginName)

	// 禁用插件
	err := config.GetCfg().DisablePlugin(request.PluginName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_update_config_failed"))
	}

	// 写入配置文件
	if writeErr := config.UpdateConfigFile(); writeErr != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), writeErr)
	}

	sse_hub.BroadcastUISuggestReload(sse_hub.UISuggestReasonPluginsChanged)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "插件已禁用",
	})
}

// DeleteArrayConfigHandler 处理删除数组元素的 JSON API
func DeleteArrayConfigHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		ConfigName  string `json:"configName"`
		DeleteValue string `json:"deleteValue"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.ConfigName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is required")
	}

	if request.DeleteValue == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "deleteValue is required")
	}

	decodedConfigName, err := decodeBase64URLStrict(request.ConfigName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "configName is not valid base64url")
	}
	decodedDeleteValue, err := decodeBase64URLStrict(request.DeleteValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "deleteValue is not valid base64url")
	}
	logger.Infof(locale.GetString("log_delete_array_config_handler")+"\n", decodedConfigName)

	values, err := doDelete(decodedConfigName, decodedDeleteValue)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_config_failed"))
	}

	// 渲染更新后的 HTML
	updatedHTML := StringArrayConfig(decodedConfigName, values, decodedConfigName+"_Description")
	htmlString, renderErr := renderTemplToString(c.Request().Context(), updatedHTML)
	if renderErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render template")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"html": htmlString,
	})
}

func doDelete(configName string, deleteValue string) ([]string, error) {
	// 更新前先保存旧配置，后续用于比较并触发副作用。
	oldConfig := config.CopyCfg()
	// 更新配置
	values, err := config.GetCfg().DeleteStringArrayConfig(configName, deleteValue)
	if err != nil {
		return nil, err
	}

	writeConfigAndApply(oldConfig)
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
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
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
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		SelectedDir string `json:"selectedDir"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
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
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		StoreUrl string `json:"storeUrl"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.StoreUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is required")
	}

	storeUrl, err := decodeBase64URLStrict(request.StoreUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is not valid base64url")
	}
	logger.Infof(locale.GetString("log_rescan_store")+"\n", storeUrl)

	// 记录扫描前的书籍数量
	beforeCount := model.GetAllBooksNumber()

	// 调用扫描功能
	err = scan.InitStore(storeUrl, config.GetCfg())
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

	logger.Infof(locale.GetString("log_rescan_store_completed_new_books")+"\n", newBooksCount)

	sse_hub.BroadcastUISuggestReload(sse_hub.UISuggestReasonSingleStoreRescan)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"newBooksCount": newBooksCount,
		"message":       locale.GetString("rescan_store_success"),
	})
}

// DeleteStoreHandler 处理删除书库的 JSON API
func DeleteStoreHandler(c echo.Context) error {
	if err := ensureWritableConfig(); err != nil {
		return err
	}

	// 解析 JSON 请求体
	var request struct {
		StoreUrl string `json:"storeUrl"`
	}
	if err := c.Bind(&request); err != nil {
		return jsonBadRequest(err)
	}

	if request.StoreUrl == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is required")
	}

	storeUrl, err := decodeBase64URLStrict(request.StoreUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "storeUrl is not valid base64url")
	}
	logger.Infof(locale.GetString("log_delete_store")+"\n", storeUrl)

	// 先删除该书库的所有书籍数据
	targetStoreAbs, err := filepath.Abs(storeUrl)
	if err != nil {
		logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
		targetStoreAbs = storeUrl
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

	logger.Infof(locale.GetString("log_deleted_books_count")+"\n", deletedCount)

	// 从配置中移除该书库 URL
	values, err := doDelete("StoreUrls", storeUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, locale.GetString("err_delete_store_failed"))
	}

	// 重新生成书组
	if err := model.IStore.GenerateBookGroup(); err != nil {
		logger.Infof(locale.GetString("log_error_initializing_main_folder"), err)
	}

	// 渲染更新后的 HTML
	updatedHTML := StoreConfig("StoreUrls", values, "StoreUrls_Description", GetStoreBookCounts())
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
