package config

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

var (
	Server         *http.Server
	Mutex          sync.Mutex
	ConfigFileLock sync.Mutex
)

const (
	cliConfigFilename     = "config.toml"
	desktopConfigFilename = "desktop.toml"
	trayConfigFilename    = "tray.toml"
)

type configLocation struct {
	name string
	dir  string
}

// PlatformConfigFilename 返回当前启动壳默认使用的配置文件名。
func PlatformConfigFilename() string {
	switch configProfile() {
	case "desktop":
		return desktopConfigFilename
	case "tray":
		return trayConfigFilename
	default:
		return cliConfigFilename
	}
}

func configDirForLocation(location string) (string, error) {
	switch location {
	case HomeDirectory:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "comigo"), nil
	case WorkingDirectory:
		return ".", nil
	case ProgramDirectory:
		executable, err := os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Dir(executable), nil
	default:
		return "", errors.New("unknown config location")
	}
}

func configFilePathForLocation(location string) (string, error) {
	configDir, err := configDirForLocation(location)
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, PlatformConfigFilename()), nil
}

func configSearchLocations() []configLocation {
	locations := make([]configLocation, 0, 3)
	if home, err := os.UserHomeDir(); err == nil {
		locations = append(locations, configLocation{
			name: HomeDirectory,
			dir:  filepath.Join(home, ".config", "comigo"),
		})
	} else {
		logger.Infof(locale.GetString("log_warning_failed_to_get_homedir"), err)
	}
	if executable, err := os.Executable(); err == nil {
		locations = append(locations, configLocation{
			name: ProgramDirectory,
			dir:  filepath.Dir(executable),
		})
	} else {
		logger.Infof(locale.GetString("log_warning_failed_to_get_executable_path"), err)
	}
	if workingDir, err := os.Getwd(); err == nil {
		locations = append(locations, configLocation{
			name: WorkingDirectory,
			dir:  workingDir,
		})
	} else {
		logger.Infof(locale.GetString("log_failed_to_get_working_directory"), err)
	}
	return locations
}

func effectiveConfigPathInDir(dir string) string {
	platformPath := filepath.Join(dir, PlatformConfigFilename())
	if fileExists(platformPath) {
		return platformPath
	}
	return ""
}

// FindConfigFile 按目录顺序寻找当前启动壳真正会读取的配置文件。
func FindConfigFile() (location string, filePath string) {
	for _, loc := range configSearchLocations() {
		if configPath := effectiveConfigPathInDir(loc.dir); configPath != "" {
			return loc.name, configPath
		}
	}
	return "", ""
}

// canSaveConfigTo 限制配置管理器只覆盖当前生效目录，避免同时维护多份生效配置。
func canSaveConfigTo(location string) bool {
	activeLocation, _ := FindConfigFile()
	return activeLocation == "" || activeLocation == location
}

func writeConfigBytes(filePath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	// 配置可能包含登录密码、Tailscale key 和远程书库凭据，只允许当前用户读取。
	return writeFileAtomically(filePath, data, 0o600)
}

// UpdateConfigFile 更新当前正在使用的配置文件。
// 如果尚未存在配置文件，则会自动在默认位置创建一份。
func UpdateConfigFile() error {
	if runtime.GOOS == "js" || cfg.TemporaryReaderMode {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	bytes, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	targetPath := cfg.ConfigFile
	if targetPath == "" {
		targetPath, err = configFilePathForLocation(DefaultConfigLocation())
		if err != nil {
			return err
		}
	}

	if err := writeConfigBytes(targetPath, bytes); err != nil {
		return err
	}
	cfg.ConfigFile = targetPath
	return nil
}

const (
	HomeDirectory    = "HomeDirectory"
	WorkingDirectory = "WorkingDirectory"
	ProgramDirectory = "ProgramDirectory"
)

// DefaultConfigLocation 判断当前配置文件应该保存到哪里；没有配置时默认保存到 HomeDirectory。
func DefaultConfigLocation() string {
	// 只有在非js环境下才需要检查配置文件位置
	if runtime.GOOS == "js" {
		return ""
	}
	if location, _ := FindConfigFile(); location != "" {
		return location
	}
	return HomeDirectory
}

// fileExists 封装一个检查文件是否存在的辅助函数
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		// 其它错误比如权限问题，也直接返回false
		return false
	}
	// 确实存在且不是目录
	return !info.IsDir()
}

func SaveConfig(to string) error {
	// 在js环境下
	if runtime.GOOS == "js" || cfg.TemporaryReaderMode {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	if !canSaveConfigTo(to) {
		return errors.New(locale.GetString("please_delete_other_config_first"))
	}

	// 保存配置
	bytes, errMarshal := toml.Marshal(cfg)
	if errMarshal != nil {
		return errMarshal
	}
	logger.Infof(locale.GetString("log_cfg_save_to"), to)
	configPath, err := configFilePathForLocation(to)
	if err != nil {
		return err
	}
	if err := writeConfigBytes(configPath, bytes); err != nil {
		return err
	}
	cfg.ConfigFile = configPath
	return nil
}

func GetWorkingDirectoryConfig() string {
	return getConfigInLocation(WorkingDirectory)
}

func GetHomeDirectoryConfig() string {
	// 在js环境下
	if runtime.GOOS == "js" {
		return ""
	}
	return getConfigInLocation(HomeDirectory)
}

func GetProgramDirectoryConfig() string {
	return getConfigInLocation(ProgramDirectory)
}

func getConfigInLocation(location string) string {
	configDir, err := configDirForLocation(location)
	if err != nil {
		return ""
	}
	configPath := effectiveConfigPathInDir(configDir)
	if configPath == "" {
		return ""
	}
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return configPath
	}
	return absPath
}

func DeleteConfigIn(in string) error {
	// 在非js环境下
	if runtime.GOOS == "js" || cfg.TemporaryReaderMode {
		return nil
	}
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()

	logger.Infof(locale.GetString("log_try_delete_cfg_in"), in)
	configDir, err := configDirForLocation(in)
	if err != nil {
		return err
	}
	configFile := effectiveConfigPathInDir(configDir)
	if configFile == "" {
		configFile = filepath.Join(configDir, PlatformConfigFilename())
	}
	return tools.DeleteFileIfExist(configFile)
}

func writeFileAtomically(filePath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filePath)
	tempFile, err := os.CreateTemp(dir, ".config-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tempFile.Name()
	cleanup := func() {
		_ = os.Remove(tmpName)
	}
	if _, err = tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		cleanup()
		return err
	}
	if err = tempFile.Chmod(perm); err != nil {
		_ = tempFile.Close()
		cleanup()
		return err
	}
	if err = tempFile.Close(); err != nil {
		cleanup()
		return err
	}
	if err = os.Rename(tmpName, filePath); err != nil {
		cleanup()
		return err
	}
	return nil
}

func GetQrcodeURL() string {
	enableTLS := cfg.CertFile != "" && cfg.KeyFile != ""
	protocol := "http://"
	if enableTLS {
		protocol = "https://"
	}
	if host := strings.TrimSpace(cfg.Host); host != "" {
		if tools.IsLoopbackHost(host) {
			host = tools.GetOutboundIP().String()
		}
		return protocol + host + ":" + strconv.Itoa(int(cfg.Port)) + PrefixPath("/")
	}
	// 取得本机的首选出站IP
	OutIP := tools.GetOutboundIP().String()
	return protocol + OutIP + ":" + strconv.Itoa(int(cfg.Port)) + PrefixPath("/")
}

// ToQrcodePublicURL 将页面 URL 中的 localhost/127.0.0.1 替换为二维码公开地址。
// 非 URL 文本或非本机地址保持原样，避免二维码接口误改普通文本。
func ToQrcodePublicURL(rawURL string) string {
	pageURL, err := url.Parse(rawURL)
	if err != nil || !pageURL.IsAbs() || pageURL.Hostname() == "" || !tools.IsLoopbackHost(pageURL.Hostname()) {
		return rawURL
	}
	publicURL, err := url.Parse(GetQrcodeURL())
	if err != nil || publicURL.Host == "" {
		return rawURL
	}
	pageURL.Scheme = publicURL.Scheme
	pageURL.Host = publicURL.Host
	return pageURL.String()
}

// GetLocalBrowserURL 返回启动时自动打开浏览器使用的本机入口，避免本机启动流程受 Host 或 0.0.0.0 影响。
func GetLocalBrowserURL() string {
	protocol := "http://"
	if cfg.EnableTLS {
		protocol = "https://"
	}
	return protocol + "127.0.0.1:" + strconv.Itoa(cfg.Port) + PrefixPath("/")
}

func OpenBrowserIfNeeded() {
	if cfg.OpenBrowser == true {
		go tools.OpenBrowserByURL(GetLocalBrowserURL())
	}
}
