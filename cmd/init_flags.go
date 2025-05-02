package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/util/locale"
)

// cobra & viper sample:
// https://qiita.com/nirasan/items/cc2ab5bc2889401fe596
var runtimeViper *viper.Viper

func init() {
	runtimeViper = viper.New()
}

func InitFlags() {
	// 加载环境变量，改写对应值
	runtimeViper.AutomaticEnv()
	// 设置环境变量的前缀，将 PORT变为 COMI_PORT
	runtimeViper.SetEnvPrefix("COMI")
	cfg := config.GetCfg()
	cobra.MousetrapHelpText = ""       // 屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 // "这是命令行程序"的提醒表示时间
	RootCmd.PersistentFlags().BoolVar(&cfg.AutoRescan, "rescan", true, locale.GetString("rescan"))
	// 启用登陆保护，需要输入用户名、密码。
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableLogin, "login", false, locale.GetString("enable_login"))
	RootCmd.PersistentFlags().StringVarP(&cfg.Username, "username", "u", "admin", locale.GetString("username"))
	RootCmd.PersistentFlags().StringVarP(&cfg.Password, "password", "k", "", locale.GetString("password"))
	RootCmd.PersistentFlags().IntVarP(&cfg.Timeout, "timeout", "t", 60*24*30, locale.GetString("timeout"))
	// TLS设定
	RootCmd.PersistentFlags().BoolVar(&cfg.EnableTLS, "tls", false, locale.GetString("tls_enable"))
	RootCmd.PersistentFlags().StringVar(&cfg.CertFile, "tls-crt", "", locale.GetString("tls_crt"))
	RootCmd.PersistentFlags().StringVar(&cfg.KeyFile, "tls-key", "", locale.GetString("tls_key"))
	// 指定配置文件
	RootCmd.PersistentFlags().StringVarP(&cfg.ConfigPath, "config", "c", "", locale.GetString("config"))
	// 启用数据库，保存扫描数据
	RootCmd.PersistentFlags().BoolVarP(&cfg.EnableDatabase, "database", "e", false, locale.GetString("enable_database"))
	// 服务端口
	RootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 1234, locale.GetString("port"))
	// 本地Host
	RootCmd.PersistentFlags().StringVar(&cfg.Host, "host", "", locale.GetString("local_host"))
	// DEBUG
	RootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, locale.GetString("debug_mode"))
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
	if runtime.GOOS == "windows" {
		cfg.OpenBrowser = true
	}
	// 不对局域网开放
	RootCmd.PersistentFlags().BoolVarP(&cfg.DisableLAN, "disable-lan", "d", false, locale.GetString("disable_lan"))
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
	RootCmd.PersistentFlags().StringVar(&cfg.CachePath, "cache-path", "", locale.GetString("cache_file_path"))
	// 退出时清除缓存
	RootCmd.PersistentFlags().BoolVar(&cfg.ClearCacheExit, "cache-clean", true, locale.GetString("cache_file_clean"))
	// 手动指定zip文件编码 gbk、shiftjis……
	RootCmd.PersistentFlags().StringVar(&cfg.ZipFileTextEncoding, "zip-encode", "gbk", locale.GetString("zip_encode"))
}
