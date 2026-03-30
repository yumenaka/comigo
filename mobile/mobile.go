package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
)

type StartConfig struct {
	Port                      int      `json:"port"`
	DisableLAN                bool     `json:"disableLAN"`
	EnableUpload              bool     `json:"enableUpload"`
	ReadOnlyMode              bool     `json:"readOnlyMode"`
	EnableDatabase            bool     `json:"enableDatabase"`
	EnablePlugin              bool     `json:"enablePlugin"`
	Debug                     bool     `json:"debug"`
	Language                  string   `json:"language"`
	Username                  string   `json:"username"`
	Password                  string   `json:"password"`
	Host                      string   `json:"host"`
	CacheDir                  string   `json:"cacheDir"`
	ConfigDir                 string   `json:"configDir"`
	StoreUrls                 []string `json:"storeUrls"`
	AutoRescanIntervalMinutes int      `json:"autoRescanIntervalMinutes"`
	ScanOnStart               bool     `json:"scanOnStart"`
	WaitReadyTimeoutSeconds   int      `json:"waitReadyTimeoutSeconds"`
}

type StartResult struct {
	Success   bool     `json:"success"`
	Running   bool     `json:"running"`
	BaseURL   string   `json:"baseUrl,omitempty"`
	HealthURL string   `json:"healthUrl,omitempty"`
	Port      int      `json:"port,omitempty"`
	StoreUrls []string `json:"storeUrls,omitempty"`
	Error     string   `json:"error,omitempty"`
}

var (
	runtimeMu     sync.Mutex
	lastStartInfo StartResult
)

func defaultStartConfig() StartConfig {
	return StartConfig{
		Port:                    1234,
		DisableLAN:              false,
		EnableUpload:            false,
		ReadOnlyMode:            false,
		EnableDatabase:          false,
		EnablePlugin:            false,
		Debug:                   false,
		Language:                "auto",
		ScanOnStart:             true,
		WaitReadyTimeoutSeconds: 20,
	}
}

// Start 启动嵌入式 Comigo 服务，供 gomobile bind 生成的 Android 包直接调用。
func Start(configJSON string) string {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	if lastStartInfo.Running {
		return encodeResult(lastStartInfo)
	}

	cfg := defaultStartConfig()
	if strings.TrimSpace(configJSON) != "" {
		if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
			return encodeError(fmt.Errorf("failed to parse mobile config: %w", err))
		}
	}

	if err := configureRuntime(cfg); err != nil {
		return encodeError(err)
	}
	if err := startEmbeddedServer(cfg); err != nil {
		return encodeError(err)
	}

	lastStartInfo = currentResult(true, "")
	return encodeResult(lastStartInfo)
}

// Stop 停止嵌入式 Comigo 服务。
func Stop() string {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	if !lastStartInfo.Running {
		return encodeResult(currentResult(false, ""))
	}

	cmd.SaveMetadata()
	config.GetCfg().AutoRescanIntervalMinutes = 0
	config.StartOrStopAutoRescan()
	if err := routers.StopWebServer(); err != nil {
		return encodeError(err)
	}

	lastStartInfo = currentResult(false, "")
	return encodeResult(lastStartInfo)
}

// GetServerInfo 返回当前嵌入式服务的运行信息。
func GetServerInfo() string {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	return encodeResult(currentResult(lastStartInfo.Running, ""))
}

// GetBaseURL 返回 Flutter WebView 应使用的本地地址。
func GetBaseURL() string {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()
	return localBaseURL()
}

// GetPort 返回当前服务端口。
func GetPort() int {
	return config.GetCfg().Port
}

func startEmbeddedServer(startCfg StartConfig) error {
	if err := routers.StartWebServer(); err != nil {
		return err
	}

	cmd.LoadUserPlugins()
	cmd.LoadMetadata()
	if startCfg.ScanOnStart {
		cmd.ScanStore()
		cmd.SaveMetadata()
	}
	config.StartOrStopAutoRescan()

	timeout := time.Duration(startCfg.WaitReadyTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	if err := waitForReady(timeout); err != nil {
		_ = routers.StopWebServer()
		return err
	}
	return nil
}

func configureRuntime(startCfg StartConfig) error {
	config.ResetConfigForRuntime()
	cfg := config.GetCfg()

	if err := prepareConfigDir(strings.TrimSpace(startCfg.ConfigDir)); err != nil {
		return err
	}
	if err := loadConfigFileIfExists(strings.TrimSpace(startCfg.ConfigDir)); err != nil {
		return err
	}

	// 移动端固定由宿主通过 WebView 访问本地 HTTP 服务，不需要桌面专用行为。
	cfg.OpenBrowser = false
	cfg.EnableTailscale = false
	cfg.EnableTLS = false
	cfg.AutoTLSCertificate = false
	cfg.ClearCacheExit = false

	cfg.Port = startCfg.Port
	cfg.DisableLAN = startCfg.DisableLAN
	cfg.EnableUpload = startCfg.EnableUpload
	cfg.ReadOnlyMode = startCfg.ReadOnlyMode
	cfg.EnableDatabase = startCfg.EnableDatabase
	cfg.EnablePlugin = startCfg.EnablePlugin
	cfg.Debug = startCfg.Debug
	cfg.Language = fallbackString(startCfg.Language, cfg.Language)
	cfg.Username = strings.TrimSpace(startCfg.Username)
	cfg.Password = strings.TrimSpace(startCfg.Password)
	cfg.Host = strings.TrimSpace(startCfg.Host)
	cfg.AutoRescanIntervalMinutes = startCfg.AutoRescanIntervalMinutes

	if trimmedCacheDir := strings.TrimSpace(startCfg.CacheDir); trimmedCacheDir != "" {
		cfg.CacheDir = trimmedCacheDir
	}
	config.AutoSetCacheDir()

	cfg.StoreUrls = []string{}
	for _, storeURL := range startCfg.StoreUrls {
		trimmedStoreURL := strings.TrimSpace(storeURL)
		if trimmedStoreURL == "" {
			continue
		}
		if err := cfg.AddStoreUrl(trimmedStoreURL); err != nil {
			return err
		}
	}

	locale.InitLanguageFromConfig(cfg.Language)
	return nil
}

func prepareConfigDir(configDir string) error {
	if configDir == "" {
		os.Unsetenv("COMIGO_CONFIG_DIR")
		return nil
	}
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return err
	}
	return os.Setenv("COMIGO_CONFIG_DIR", configDir)
}

func loadConfigFileIfExists(configDir string) error {
	if configDir == "" {
		return nil
	}
	configFilePath := filepath.Join(configDir, "config.toml")
	config.GetCfg().ConfigFile = configFilePath
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return nil
	}
	return toml.Unmarshal(content, config.GetCfg())
}

func waitForReady(timeout time.Duration) error {
	healthURL := localBaseURL() + "/healthz"
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		response, err := client.Get(healthURL)
		if err == nil {
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("unexpected health status: %d", response.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(200 * time.Millisecond)
	}

	if lastErr == nil {
		lastErr = errors.New("health check timed out")
	}
	return fmt.Errorf("comigo embedded server not ready: %w", lastErr)
}

func currentResult(running bool, errMessage string) StartResult {
	return StartResult{
		Success:   errMessage == "",
		Running:   running,
		BaseURL:   localBaseURL(),
		HealthURL: localBaseURL() + "/healthz",
		Port:      config.GetCfg().Port,
		StoreUrls: append([]string(nil), config.GetCfg().StoreUrls...),
		Error:     errMessage,
	}
}

func localBaseURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d", config.GetCfg().Port)
}

func encodeResult(result StartResult) string {
	bytes, err := json.Marshal(result)
	if err != nil {
		return `{"success":false,"running":false,"error":"failed to marshal mobile result"}`
	}
	return string(bytes)
}

func encodeError(err error) string {
	lastStartInfo = currentResult(false, err.Error())
	return encodeResult(lastStartInfo)
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
