package config

import (
	"os"
	"path"

	"github.com/yumenaka/comigo/config/stores"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

func GetCfg() *Config {
	return &cfg
}

func CopyCfg() Config {
	return cfg
}

func GetConfigPath() string {
	return cfg.ConfigPath
}

func SetConfigPath(path string) {
	// 检查路径是否存在
	if !util.PathExists(path) {
		logger.Info("Invalid config file path.")
		return
	}
	cfg.ConfigPath = path
}

func GetCachePath() string {
	return cfg.CachePath
}

func SetCachePath(path string) {
	if !util.PathExists(path) {
		logger.Info("Invalid cache path.")
		return
	}
	cfg.CachePath = path
}

func AutoSetCachePath() {
	// 手动设置的临时文件夹
	if cfg.CachePath != "" && util.IsExist(cfg.CachePath) && util.ChickIsDir(cfg.CachePath) {
		cfg.CachePath = path.Join(cfg.CachePath)
	} else {
		cfg.CachePath = path.Join(os.TempDir(), "comigo_cache") // 使用系统文件夹
	}
	err := os.MkdirAll(cfg.CachePath, os.ModePerm)
	if err != nil {
		logger.Infof("%s", locale.GetString("temp_folder_error"))
	} else {
		logger.Infof("%s", locale.GetString("temp_folder_path")+cfg.CachePath)
	}
}

func GetClearDatabaseWhenExit() bool {
	return cfg.ClearDatabaseWhenExit
}
func SetClearDatabaseWhenExit(clearDatabaseWhenExit bool) {
	cfg.ClearDatabaseWhenExit = clearDatabaseWhenExit
}

func GetDebug() bool {
	return cfg.Debug
}

func SetDebug(debug bool) {
	cfg.Debug = debug
}

func GetEnableUpload() bool {
	return cfg.EnableUpload
}

func GetEnableDatabase() bool {
	return cfg.EnableDatabase
}

func SetEnableDatabase(enableDatabase bool) {
	cfg.EnableDatabase = enableDatabase
}

func GetEnableTLS() bool {
	return cfg.EnableTLS
}

func GetUploadPath() string {
	return cfg.UploadPath
}

func SetUploadPath(path string) {
	if (!util.IsDir(path)) || (!util.PathExists(path)) {
		logger.Info("Invalid upload path.")
		return
	}
	cfg.UploadPath = path
}

func GetUseCache() bool {
	return cfg.UseCache
}

func SetUseCache(useCache bool) {
	cfg.UseCache = useCache
}

func GetCertFile() string {
	return cfg.CertFile
}

func GetClearCacheExit() bool {
	return cfg.ClearCacheExit
}

func GetLogToFile() bool {
	return cfg.LogToFile
}

func GetLogFilePath() string {
	return cfg.LogFilePath
}

func GetLogFileName() string {
	return cfg.LogFileName
}

func GetStores() []stores.Store {
	return cfg.Stores
}

func GetMaxScanDepth() int {
	return cfg.MaxScanDepth
}

func SetClearCacheExit(clearCacheExit bool) {
	cfg.ClearCacheExit = clearCacheExit
}

func GetDefaultMode() string {
	return cfg.DefaultMode
}

func GetDisableLAN() bool {
	return cfg.DisableLAN
}

func SetDisableLAN(disableLAN bool) {
	cfg.DisableLAN = disableLAN
}

func GetMinImageNum() int {
	return cfg.MinImageNum
}

func GetTimeoutLimitForScan() int {
	return cfg.TimeoutLimitForScan
}

func GetExcludePath() []string {
	return cfg.ExcludePath
}

func GetSupportMediaType() []string {
	return cfg.SupportMediaType
}

func GetSupportFileType() []string {
	return cfg.SupportFileType
}

func GetSupportTemplateFile() []string {
	return cfg.SupportTemplateFile
}

func GetZipFileTextEncoding() string {
	return cfg.ZipFileTextEncoding
}

func GetOpenBrowser() bool {
	return cfg.OpenBrowser
}

func SetOpenBrowser(openBrowser bool) {
	cfg.OpenBrowser = openBrowser
}

func GetPrintAllPossibleQRCode() bool {
	return cfg.PrintAllPossibleQRCode
}

func GetPort() int {
	return cfg.Port
}

func GetUsername() string {
	return cfg.Username
}

// GetJwtSigningKey JWT令牌签名key，目前是用户名+密码
func GetJwtSigningKey() string {
	return cfg.Username + cfg.Password
}

func GetPassword() string {
	return cfg.Password
}

func GetTimeout() int {
	return cfg.Timeout
}

func SetPort(port int) {
	if port < 0 || port > 65535 {
		port = 1234
		logger.Infof("Invalid port number. Using default port: %d", port)
	}
	cfg.Port = port
}

func GetHost() string {
	return cfg.Host
}

func SetHost(host string) {
	// 如果主机名为空，使用默认主机名
	if host == "" {
		host = ""
		logger.Infof("Invalid host name. Using default host: %s", host)
	}
	cfg.Host = host
}

func GetKeyFile() string {
	return cfg.KeyFile
}
