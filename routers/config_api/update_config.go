package config_api

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/service"
)

// UpdateConfig 修改服务器配置(post json)
func UpdateConfig(c echo.Context) error {
	// 读取请求体中的 JSON 数据
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return apiresp.BadRequest(c, "read_body_failed", "Failed to read request body", err.Error())
	}
	// 将 JSON 数据转换为字符串并打印
	jsonString := string(body)
	logger.Infof(locale.GetString("log_received_json_data"), jsonString)
	// 如果配置被锁定，返回错误
	if config.GetCfg().ReadOnlyMode {
		return apiresp.Forbidden(c, "config_locked", "Config is locked, cannot be modified", nil)
	}
	// 复制当前配置以便后续比较
	oldConfig := config.CopyCfg()
	// 解析 JSON 数据并更新服务器配置
	err = config.UpdateConfigByJson(jsonString)
	if err != nil {
		logger.Infof("%s", err.Error())
		return apiresp.BadRequest(c, "invalid_json", locale.GetString("log_failed_to_parse_json"), err.Error())
	}
	// 如果 Language 配置发生变化，重新初始化语言设置
	if oldConfig.Language != config.GetCfg().Language {
		locale.InitLanguageFromConfig(config.GetCfg().Language)
	}
	err = config.UpdateConfigFile()
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_update_local_config"), err)
	}
	// 根据配置变化执行副作用：打开浏览器、重扫、自动扫描调度等。
	service.ApplyConfigChange(&oldConfig, config.GetCfg(), nil)
	// 返回成功消息
	return apiresp.Success(c, "ok", "Server settings updated successfully", nil)
}
