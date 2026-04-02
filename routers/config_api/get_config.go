package config_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/tools/logger"
)

// GetConfig 获取json格式的当前配置，不做修改
func GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, config.GetCfg())
}

// GetConfigToml 下载服务器配置(toml)，修改关键值后上传
func GetConfigToml(c echo.Context) error {
	// golang结构体默认深拷贝（但是基本类型浅拷贝）
	tempConfig := config.GetCfg()
	tempConfig.LogFilePath = ""
	tempConfig.OpenBrowser = false
	tempConfig.EnableDatabase = true
	tempConfig.StoreUrls = []string{"C:\\test\\Comic", "D:\\some_path\\book", "/home/user/download"}
	tempConfig.LoginProtection = false
	tempConfig.Username = "You_can_change_this_username"
	tempConfig.Password = "Some_Secret-.Password密码"
	tempConfig.EnableOAuthLogin = false
	tempConfig.OAuthProviderName = "GitHub"
	tempConfig.OAuthClientID = "your-client-id"
	tempConfig.OAuthClientSecret = "your-client-secret"
	tempConfig.OAuthAuthURL = "https://example.com/oauth/authorize"
	tempConfig.OAuthTokenURL = "https://example.com/oauth/token"
	tempConfig.OAuthUserInfoURL = "https://example.com/oauth/userinfo"
	tempConfig.OAuthRedirectURL = "https://your-host.example.com/api/oauth/callback"
	tempConfig.OAuthScopes = []string{"openid", "profile", "email"}

	bytes, err := toml.Marshal(tempConfig)
	if err != nil {
		logger.Infof(locale.GetString("log_toml_marshal_error"))
		return err
	}

	// 在命令行打印
	logger.Infof("%s", string(bytes))

	// 设置响应头，指定文件下载名称和类型
	return c.Blob(
		http.StatusOK,
		"application/octet-stream",
		bytes,
	)
}
