package cmd

import (
	"os"
	"path"
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
	cobra.MousetrapHelpText = ""       // 屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 //"这是命令行程序"的提醒表示时间
	//rootCmd.PersistentFlags().BoolVar(&config.Config.AutoRescan, "rescan", true, locale.GetString("rescan"))
	// 启用登陆保护，需要输入用户名、密码。
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableLogin, "login", false, locale.GetString("enable_login"))
	rootCmd.PersistentFlags().StringVarP(&config.Config.Username, "username", "u", "admin", locale.GetString("username"))
	rootCmd.PersistentFlags().StringVarP(&config.Config.Password, "password", "k", "", locale.GetString("password"))
	rootCmd.PersistentFlags().IntVarP(&config.Config.Timeout, "timeout", "t", 65535, locale.GetString("timeout"))
	// TLS设定
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableTLS, "tls", false, locale.GetString("tls_enable"))
	rootCmd.PersistentFlags().StringVar(&config.Config.CertFile, "tls-crt", "", locale.GetString("tls_crt"))
	rootCmd.PersistentFlags().StringVar(&config.Config.KeyFile, "tls-key", "", locale.GetString("tls_key"))
	// 指定配置文件
	rootCmd.PersistentFlags().StringVarP(&config.Config.ConfigPath, "config", "c", "", locale.GetString("config"))
	// 启用数据库，保存扫描数据
	rootCmd.PersistentFlags().BoolVarP(&config.Config.EnableDatabase, "database", "e", false, locale.GetString("enable_database"))
	// 服务端口
	rootCmd.PersistentFlags().IntVarP(&config.Config.Port, "port", "p", 1234, locale.GetString("port"))
	// 本地Host
	rootCmd.PersistentFlags().StringVar(&config.Config.Host, "host", "DefaultHost", locale.GetString("local_host"))
	// DEBUG
	rootCmd.PersistentFlags().BoolVar(&config.Config.Debug, "debug", false, locale.GetString("debug_mode"))
	// 启用文件上传功能
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableUpload, "enable-upload", true, locale.GetString("enable_file_upload"))
	// 上传文件的保存路径
	rootCmd.PersistentFlags().StringVar(&config.Config.UploadPath, "upload-path", "", locale.GetString("upload_path"))
	if config.Config.EnableUpload && config.Config.UploadPath == "" {
		// 获取当前目录
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			config.Config.UploadPath = path.Join(dir, "upload")
		}
	}
	// 打开浏览器
	rootCmd.PersistentFlags().BoolVarP(&config.Config.OpenBrowser, "open-browser", "o", false, locale.GetString("open_browser"))
	if runtime.GOOS == "windows" {
		config.Config.OpenBrowser = true
	}
	// 不对局域网开放
	rootCmd.PersistentFlags().BoolVarP(&config.Config.DisableLAN, "disable-lan", "d", false, locale.GetString("disable_lan"))
	// 文件搜索深度
	rootCmd.PersistentFlags().IntVarP(&config.Config.MaxScanDepth, "max-depth", "m", 5, locale.GetString("max_depth"))
	////服务器解析书籍元数据，如果生成blurhash，需要消耗大量资源
	rootCmd.PersistentFlags().BoolVar(&config.Config.GenerateMetaData, "generate-metadata", false, locale.GetString("generate_metadata"))
	// 打印所有可用网卡ip
	rootCmd.PersistentFlags().BoolVar(&config.Config.PrintAllPossibleQRCode, "print-all", false, locale.GetString("print_all_ip"))
	// 至少有几张图片，才认定为漫画压缩包
	rootCmd.PersistentFlags().IntVar(&config.Config.MinImageNum, "min-image", 1, locale.GetString("min_media_num"))
	// 输出log文件
	rootCmd.PersistentFlags().BoolVar(&config.Config.LogToFile, "log", false, locale.GetString("log_to_file"))
	// web图片缓存
	rootCmd.PersistentFlags().BoolVar(&config.Config.UseCache, "use-cache", false, locale.GetString("cache_file_enable"))
	// 图片缓存路径
	rootCmd.PersistentFlags().StringVar(&config.Config.CachePath, "cache-path", "", locale.GetString("cache_file_path"))
	// 退出时清除缓存
	rootCmd.PersistentFlags().BoolVar(&config.Config.ClearCacheExit, "cache-clean", true, locale.GetString("cache_file_clean"))
	// 手动指定zip文件编码 gbk、shiftjis……
	rootCmd.PersistentFlags().StringVar(&config.Config.ZipFileTextEncoding, "zip-encode", "gbk", locale.GetString("zip_encode"))
}
