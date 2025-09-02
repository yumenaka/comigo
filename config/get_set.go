package config

import (
	"os"
	"path"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
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
	if !tools.PathExists(path) {
		logger.Info("Invalid config file path.")
		return
	}
	cfg.ConfigPath = path
}

func GetCachePath() string {
	return cfg.CachePath
}

func SetCachePath(path string) {
	if !tools.PathExists(path) {
		logger.Info("Invalid cache path.")
		return
	}
	cfg.CachePath = path
}

func AutoSetCachePath() {
	// 手动设置的临时文件夹
	if cfg.CachePath != "" && tools.IsExist(cfg.CachePath) && tools.ChickIsDir(cfg.CachePath) {
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

func GetStoreUrls() []string {
	return cfg.StoreUrls
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
	if (!tools.IsDir(path)) || (!tools.PathExists(path)) {
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

func GetTailscaleEnable() bool {
	return cfg.EnableTailscale
}

func GetTailscaleHostname() string {
	return cfg.TailscaleHostname
}

func GetTailscalePort() int {
	return cfg.TailscalePort
}

// GetUsername 获取用户名
func GetUsername() string {
	return cfg.Username
}

// GetPassword 获取密码
func GetPassword() string {
	return cfg.Password
}

// GetJwtSigningKey JWT令牌签名key，目前是用户名+密码(如果两者都设置了的话)
func GetJwtSigningKey() string {
	if cfg.Username == "" || cfg.Password == "" {
		logger.Infof("Username or password is empty. Using default Jwt Signing key.")
		tempStr := cfg.Username + cfg.Password + GetVersion()
		for _, store := range cfg.StoreUrls {
			tempStr = tempStr + store
		}
		base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(tempStr))))
	}
	return cfg.Username + cfg.Password
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
