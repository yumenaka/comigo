package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// cobra & viper sample:
// https://qiita.com/nirasan/items/cc2ab5bc2889401fe596
var runtimeViper *viper.Viper

func init() {
	runtimeViper = viper.New()
	// 为了避免循环引用，把 model.IStore 指向 store.MainStoreGroup
	model.IStore = store.RamStore
}

func InitFlags() {
	// 加载环境变量，改写对应值
	runtimeViper.AutomaticEnv()
	// 设置环境变量的前缀，将 PORT变为 COMIGO_PORT
	runtimeViper.SetEnvPrefix("COMIGO")
	cfg := config.GetCfg()
	cobra.MousetrapHelpText = ""       // 屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 // "这是命令行程序"的提醒表示时间
	// 指定配置文件
	RootCmd.PersistentFlags().StringVarP(&cfg.ConfigFile, "config", "c", "", locale.GetString("config"))
	// 启用登陆保护，需要设定用户名
	RootCmd.PersistentFlags().StringVar(&cfg.Username, "username", "", locale.GetString("username"))
	RootCmd.PersistentFlags().StringVar(&cfg.Password, "password", "", locale.GetString("password"))
	RootCmd.PersistentFlags().IntVar(&cfg.Timeout, "timeout", 60*24*30, locale.GetString("timeout"))
	// 启用自动扫描间隔，单位分钟，0为禁用自动扫描
	RootCmd.PersistentFlags().IntVar(&cfg.AutoRescanIntervalMinutes, "auto-rescan-min", 0, locale.GetString("auto_rescan_interval_minutes"))
	// 启用数据库，保存扫描数据
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableDatabase, "database", false, locale.GetString("enable_database"))
	// 服务端口
	RootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 1234, locale.GetString("port"))
	// 本地Host
	RootCmd.PersistentFlags().StringVar(&cfg.Host, "host", "", locale.GetString("local_host"))
	// TLS设定
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableTLS, "tls", false, locale.GetString("tls_enable"))
	RootCmd.PersistentFlags().BoolVar(&cfg.AutoTLSCertificate, "auto-tls", false, locale.GetString("auto_https_cert"))
	RootCmd.PersistentFlags().StringVar(&cfg.CertFile, "tls-crt", "", locale.GetString("tls_crt"))
	RootCmd.PersistentFlags().StringVar(&cfg.KeyFile, "tls-key", "", locale.GetString("tls_key"))
	// 启用文件上传功能
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableUpload, "enable-upload", true, locale.GetString("enable_file_upload"))
	// 上传文件的保存路径
	RootCmd.PersistentFlags().StringVar(&cfg.UploadPath, "upload-path", "", locale.GetString("upload_path"))
	if cfg.EnableUpload && cfg.UploadPath == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			config.SetUploadPath(dir)
		}
	}
	// 打开浏览器
	RootCmd.PersistentFlags().BoolVarP(&cfg.OpenBrowser, "open-browser", "o", false, locale.GetString("open_browser"))
	// 不对局域网开放
	RootCmd.PersistentFlags().BoolVar(&cfg.DisableLAN, "local", false, locale.GetString("disable_lan"))
	// 文件搜索深度
	RootCmd.PersistentFlags().IntVarP(&cfg.MaxScanDepth, "max-depth", "m", 5, locale.GetString("max_depth"))
	// //服务器解析书籍元数据，如果生成blurhash，需要消耗大量资源
	RootCmd.PersistentFlags().BoolVar(&cfg.GenerateMetaData, "generate-metadata", false, locale.GetString("generate_metadata"))
	// 打印所有可用网卡ip
	RootCmd.PersistentFlags().BoolVar(&cfg.PrintAllPossibleQRCode, "print-all", false, locale.GetString("print_all_ip"))
	// 至少有几张图片，才认定为漫画压缩包
	RootCmd.PersistentFlags().IntVar(&cfg.MinImageNum, "min-image", 1, locale.GetString("min_media_num"))
	// 输出log文件
	RootCmd.PersistentFlags().BoolVar(&cfg.LogToFile, "log-file", false, locale.GetString("log_to_file"))
	// web图片缓存
	RootCmd.PersistentFlags().BoolVar(&cfg.UseCache, "use-cache", false, locale.GetString("cache_file_enable"))
	// 图片缓存路径
	RootCmd.PersistentFlags().StringVar(&cfg.CacheDir, "cache-dir", "", locale.GetString("cache_file_dir"))
	// 退出时清除缓存
	RootCmd.PersistentFlags().BoolVar(&cfg.ClearCacheExit, "cache-clean", false, locale.GetString("cache_file_clean"))
	// 手动指定zip文件编码 gbk、shiftjis……
	RootCmd.PersistentFlags().StringVar(&cfg.ZipFileTextEncoding, "zip-encode", "gbk", locale.GetString("zip_encode"))
	// 启用Tailscale服务
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableTailscale, "tailscale", false, locale.GetString("EnableTailscale"))
	// Tailscale服务 启用Funnel模式
	RootCmd.PersistentFlags().BoolVar(&cfg.FunnelTunnel, "tailscale-funnel", false, locale.GetString("FunnelTunnel"))
	// FunnelLoginCheck Funnel密码保护检查
	RootCmd.PersistentFlags().BoolVar(&cfg.FunnelLoginCheck, "funnel-password-check", true, locale.GetString("FunnelLoginCheck"))
	// Tailscale服务主机名,用于 Tailscale 网络中的标识节点
	RootCmd.PersistentFlags().StringVar(&cfg.TailscaleHostname, "tailscale-hostname", "comigo", locale.GetString("TailscaleHostname"))
	// Tailscale服务端口号
	RootCmd.PersistentFlags().IntVar(&cfg.TailscalePort, "tailscale-port", 443, locale.GetString("TailscalePort"))
	// Tailscale AuthKey
	RootCmd.PersistentFlags().StringVar(&cfg.TailscaleAuthKey, "tailscale-authKey", "", locale.GetString("TailscaleAuthKey"))
	// ReadOnlyMode 只读模式，禁止网页端修改配置或上传文件
	RootCmd.PersistentFlags().BoolVar(&cfg.ReadOnlyMode, "read-only", false, locale.GetString("read_only_mode"))
	// EnableSingleInstance 启用单实例模式
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableSingleInstance, "single-instance", false, locale.GetString("enable_single_instance"))
	// Language 语言设置
	RootCmd.PersistentFlags().StringVar(&cfg.Language, "lang", "auto", locale.GetString("lang"))
	// Windows 右键菜单注册/卸载，仅在 Windows 下生效
	if runtime.GOOS == "windows" {
		RootCmd.PersistentFlags().BoolVar(&cfg.RegisterContextMenu, "register-context-menu", false, locale.GetString("register_context_menu"))
		RootCmd.PersistentFlags().BoolVar(&cfg.UnregisterContextMenu, "unregister-context-menu", false, locale.GetString("unregister_context_menu"))
	}
	// EnablePlugin
	RootCmd.PersistentFlags().BoolVar(&cfg.EnablePlugin, "plugin", true, locale.GetString("plugin_enable"))
	// DEBUG
	RootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, locale.GetString("debug_mode"))
}

// SetByExecutableFilename 根据可执行文件名或软链接名自动设置部分配置项的默认值。
// 该函数会检查可执行文件的名称（在类Unix系统中，如果通过软链接调用，会获取软链接名），
// 并根据文件名中包含的关键词自动设置相应的配置项。
//
// 支持的配置项及对应的文件名关键词：
//
//   - 阅读模式：
//
//   - "flip" → 设置为翻页模式
//
//   - "scroll" → 设置为滚动模式
//
//   - 布尔值配置项（文件名包含关键词时启用）：
//
//   - "rescan" → 启用自动扫描 (AutoRescan)
//
//   - "database" → 启用数据库 (EnableDatabase)
//
//   - "enable-upload" 或 "upload" → 启用文件上传 (EnableUpload)
//
//   - "local" → 仅本地访问，不对局域网开放 (DisableLAN)
//
//   - "generate-metadata" 或 "metadata" → 生成元数据 (GenerateMetaData)
//
//   - "log-file" 或 "logfile" → 输出日志文件 (LogToFile)
//
//   - "use-cache" 或 "cache" → 启用缓存 (UseCache)
//
//   - "cache-clean" 或 "cleancache" → 退出时清除缓存 (ClearCacheExit)
//
//   - "read-only" 或 "readonly" → 启用只读模式 (ReadOnlyMode)
//
//   - "single-instance" 或 "singleinstance" → 启用单实例模式 (EnableSingleInstance)
//
//   - "open-browser" → 打开浏览器 (OpenBrowser)
//
//   - "debug" → 启用调试模式 (Debug)
//
//   - 语言设置：
//
//   - "zh" 或 "chinese" → 设置为中文 (Language = "zh")
//
//   - "en" 或 "english" → 设置为英文 (Language = "en")
//
//   - "ja" 或 "japanese" → 设置为日文 (Language = "ja")
//
// 特殊行为：
//   - 在 Windows 系统上，默认会启用 OpenBrowser（打开浏览器）
//   - 在 Windows 系统上，文件名会移除 .exe 或 .EXE 后缀（保持原始大小写）
//   - 关键词匹配是大小写不敏感的
//
// 使用示例：
//   - 可执行文件名为 "comigo-rescan" → 启用自动扫描
//   - 可执行文件名为 "comigo-open-browser-zh" → 打开浏览器并设置语言为中文
//   - 可执行文件名为 "comigo-readonly" → 启用只读模式
func SetByExecutableFilename() {
	// 获取可执行文件的名称,如果在类Unix系统里通过软链接调用,会拿到"软链接名"（也就是别名）
	filename := filepath.Base(os.Args[0])
	// 在 Windows 系统上，移除 .exe 或 .EXE 后缀（保持原始大小写）
	if runtime.GOOS == "windows" {
		filename = strings.TrimSuffix(filename, ".exe")
		filename = strings.TrimSuffix(filename, ".EXE")
	}
	// 转换为小写用于大小写不敏感的匹配
	filenameLower := strings.ToLower(filename)
	cfg := config.GetCfg()
	// Windows 默认打开浏览器
	if runtime.GOOS == "windows" {
		cfg.OpenBrowser = true
	}
	// 打开浏览器
	if strings.Contains(filenameLower, "open-browser") {
		cfg.OpenBrowser = true
	}
	// Debug 模式
	if strings.Contains(filenameLower, "debug") {
		cfg.Debug = true
	}
	// 启用 EnablePlugin
	if strings.Contains(filenameLower, "plugin") {
		cfg.EnablePlugin = true
	}
	// 自动扫描间隔
	if strings.Contains(filenameLower, "autorescan") {
		if cfg.AutoRescanIntervalMinutes == 0 {
			cfg.AutoRescanIntervalMinutes = 60 // 默认60分钟
		}
	}
	// 启用数据库
	if strings.Contains(filenameLower, "database") {
		cfg.EnableDatabase = true
	}
	// 启用文件上传
	if strings.Contains(filenameLower, "enable-upload") || strings.Contains(filenameLower, "upload") {
		cfg.EnableUpload = true
	}
	// 仅本地访问
	if strings.Contains(filenameLower, "local") {
		cfg.DisableLAN = true
	}
	// 生成元数据
	if strings.Contains(filenameLower, "generate-metadata") || strings.Contains(filenameLower, "metadata") {
		cfg.GenerateMetaData = true
	}
	// 输出日志文件
	if strings.Contains(filenameLower, "log-file") || strings.Contains(filenameLower, "logfile") {
		cfg.LogToFile = true
	}
	// 启用缓存
	if strings.Contains(filenameLower, "use-cache") || strings.Contains(filenameLower, "cache") {
		cfg.UseCache = true
	}
	// 退出时清除缓存
	if strings.Contains(filenameLower, "cache-clean") || strings.Contains(filenameLower, "cleancache") {
		cfg.ClearCacheExit = true
	}
	// 只读模式
	if strings.Contains(filenameLower, "read-only") || strings.Contains(filenameLower, "readonly") {
		cfg.ReadOnlyMode = true
	}
	// 单实例模式
	if strings.Contains(filenameLower, "comigo") || strings.Contains(filenameLower, "single-instance") || strings.Contains(filenameLower, "singleinstance") {
		cfg.EnableSingleInstance = true
	}
	// 语言设置
	if strings.Contains(filenameLower, "zh") || strings.Contains(filenameLower, "chinese") {
		cfg.Language = "zh"
	} else if strings.Contains(filenameLower, "en") || strings.Contains(filenameLower, "english") {
		cfg.Language = "en"
	} else if strings.Contains(filenameLower, "ja") || strings.Contains(filenameLower, "japanese") {
		cfg.Language = "ja"
	}
	if cfg.Debug {
		logger.Infof(locale.GetString("log_executable_name"), filename)
		cfg.EnabledPluginList = []string{"clock", "auto_flip", "auto_scroll", "comigo_xyz", "sample"}
	}
}
