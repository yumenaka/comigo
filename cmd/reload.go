package cmd

import (
	"errors"
	"fmt"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools/service"
)

// reloadDefaults 保存读取 TOML 之前的默认值；删除配置项时恢复默认值，而非保留上一次的旧值。
// runtimeViper 仍绑定启动 flag 和环境变量，所以重载遵守与启动相同的优先级。
var reloadDefaults config.Config

// readReloadConfig 先完整解析并校验候选配置，不改动正在使用的全局配置。
func readReloadConfig(file string) (config.Config, error) {
	candidate := reloadDefaults
	runtimeViper.SetConfigFile(file)
	if err := runtimeViper.ReadInConfig(); err != nil {
		return candidate, err
	}
	if err := runtimeViper.Unmarshal(&candidate); err != nil {
		return candidate, err
	}
	candidate.ConfigFile = file
	if candidate.Port < 1 || candidate.Port > 65535 {
		return candidate, fmt.Errorf(locale.GetString("cli_invalid_port"), candidate.Port)
	}
	if err := candidate.ValidateTLS(); err != nil {
		return candidate, err
	}
	// 重载切换缓存目录时同样先创建并校验，失败不修改运行配置。
	if candidate.CacheDir != "" && candidate.CacheDir != config.GetCfg().CacheDir {
		if err := candidate.InitCacheDir(); err != nil {
			return candidate, err
		}
	}

	// 启动时的位置参数仍然有效，不能因重载配置丢失命令行指定的书库。
	for _, path := range Args {
		// 网页保存的 TOML 可能已包含启动路径；与启动时一样跳过已覆盖的书库。
		if overlapping, existing, _ := candidate.IsPathOverlapping(path); overlapping && existing != "" {
			continue
		}
		if err := candidate.AddStoreUrl(path); err != nil {
			return candidate, err
		}
	}
	return candidate, nil
}

// reloadConfig 在原进程应用 TOML；监听重启失败时恢复旧配置和服务，并向 CLI 返回错误。
// 不写回 TOML，避免覆盖用户刚编辑的文件或将命令行 flag 固化到文件。
func reloadConfig(file string) error {
	old := config.CopyCfg()
	if file == "" {
		file = old.ConfigFile
	}
	if file == "" {
		return fmt.Errorf("%s", locale.GetString("cli_reload_no_config"))
	}
	candidate, err := readReloadConfig(file)
	if err != nil {
		return err
	}
	// 存储后端只能在启动时选择，避免配置已切到 SQLite 而请求仍写 JSON。
	if candidate.EnableDatabase != old.EnableDatabase || candidate.DBType != old.DBType || candidate.DBDSN != old.DBDSN {
		return fmt.Errorf("database config requires restart")
	}
	// 这些字段由启动流程或插件扫描计算，不属于配置文件的持久化内容。
	candidate.TemporaryReaderMode, candidate.NoTUI = old.TemporaryReaderMode, old.NoTUI
	candidate.CustomPlugins, candidate.UserPluginList = old.CustomPlugins, old.UserPluginList
	candidate.PluginDirectory = old.PluginDirectory
	if candidate.CacheDir == "" {
		candidate.CacheDir = old.CacheDir
	}
	if candidate.EnablePlugin {
		_ = candidate.AddPlugin("auto_flip")
		_ = candidate.AddPlugin("auto_scroll")
	}
	action := service.BuildConfigChangeAction(old, &candidate)
	config.Mutex.Lock()
	*config.GetCfg() = candidate
	config.Mutex.Unlock()
	if action.ReStartWebServer {
		if err := routers.RestartWebServer(); err != nil {
			config.Mutex.Lock()
			*config.GetCfg() = old
			config.Mutex.Unlock()
			restoreErr := routers.RestartWebServer()
			routers.StartTailscale()
			return errors.Join(err, restoreErr)
		}
	}
	if action.ReStartWebServer || action.StartTailscale || action.ReStartTailscale {
		routers.StartTailscale()
	} else if action.StopTailscale {
		routers.StopTailscale()
	}
	locale.InitLanguageFromConfig(candidate.Language)
	LoadUserPlugins()
	service.ApplyConfigChange(old, config.GetCfg(), nil)
	return nil
}
