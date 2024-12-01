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
	rootCmd.PersistentFlags().BoolVar(&config.Config.AutoRescan, "rescan", true, locale.GetString("RESCAN"))
	// 启用登陆保护，需要输入用户名、密码。
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableLogin, "login", false, locale.GetString("ENABLE_LOGIN"))
	rootCmd.PersistentFlags().StringVarP(&config.Config.Username, "username", "u", "admin", locale.GetString("USERNAME"))
	rootCmd.PersistentFlags().StringVarP(&config.Config.Password, "password", "k", "", locale.GetString("PASSWORD"))
	rootCmd.PersistentFlags().IntVarP(&config.Config.Timeout, "timeout", "t", 65535, locale.GetString("TIMEOUT"))
	// TLS设定
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableTLS, "tls", false, locale.GetString("TLS_ENABLE"))
	rootCmd.PersistentFlags().StringVar(&config.Config.CertFile, "tls-crt", "", locale.GetString("TLS_CRT"))
	rootCmd.PersistentFlags().StringVar(&config.Config.KeyFile, "tls-key", "", locale.GetString("TLS_KEY"))
	// 指定配置文件
	rootCmd.PersistentFlags().StringVarP(&config.Config.ConfigPath, "config", "c", "", locale.GetString("CONFIG"))
	// 启用数据库，保存扫描数据
	rootCmd.PersistentFlags().BoolVarP(&config.Config.EnableDatabase, "database", "e", false, locale.GetString("EnableDatabase"))
	// 服务端口
	rootCmd.PersistentFlags().IntVarP(&config.Config.Port, "port", "p", 1234, locale.GetString("PORT"))
	// 本地Host
	rootCmd.PersistentFlags().StringVar(&config.Config.Host, "host", "DefaultHost", locale.GetString("LOCAL_HOST"))
	// DEBUG
	rootCmd.PersistentFlags().BoolVar(&config.Config.Debug, "debug", false, locale.GetString("DEBUG_MODE"))
	// 启用文件上传功能
	rootCmd.PersistentFlags().BoolVar(&config.Config.EnableUpload, "enable-upload", true, locale.GetString("ENABLE_FILE_UPLOAD"))
	// 上传文件的保存路径
	rootCmd.PersistentFlags().StringVar(&config.Config.UploadPath, "upload-path", "", locale.GetString("UPLOAD_PATH"))
	if config.Config.EnableUpload && config.Config.UploadPath == "" {
		// 获取当前目录
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			config.Config.UploadPath = path.Join(dir, "upload")
		}
	}
	// 打开浏览器
	rootCmd.PersistentFlags().BoolVarP(&config.Config.OpenBrowser, "open-browser", "o", false, locale.GetString("OPEN_BROWSER"))
	if runtime.GOOS == "windows" {
		config.Config.OpenBrowser = true
	}
	// 不对局域网开放
	rootCmd.PersistentFlags().BoolVarP(&config.Config.DisableLAN, "disable-lan", "d", false, locale.GetString("DISABLE_LAN"))
	// 文件搜索深度
	rootCmd.PersistentFlags().IntVarP(&config.Config.MaxScanDepth, "max-depth", "m", 5, locale.GetString("MAX_DEPTH"))
	////服务器解析书籍元数据，如果生成blurhash，需要消耗大量资源
	rootCmd.PersistentFlags().BoolVar(&config.Config.GenerateMetaData, "generate-metadata", false, locale.GetString("GENERATE_METADATA"))
	// 打印所有可用网卡ip
	rootCmd.PersistentFlags().BoolVar(&config.Config.PrintAllPossibleQRCode, "print-all", false, locale.GetString("PRINT_ALL_IP"))
	// 至少有几张图片，才认定为漫画压缩包
	rootCmd.PersistentFlags().IntVar(&config.Config.MinImageNum, "min-image", 1, locale.GetString("MIN_MEDIA_NUM"))
	// 输出log文件
	rootCmd.PersistentFlags().BoolVar(&config.Config.LogToFile, "log", false, locale.GetString("LOG_TO_FILE"))
	// web图片缓存
	rootCmd.PersistentFlags().BoolVar(&config.Config.UseCache, "use-cache", false, locale.GetString("CACHE_FILE_ENABLE"))
	// 图片缓存路径
	rootCmd.PersistentFlags().StringVar(&config.Config.CachePath, "cache-path", "", locale.GetString("CACHE_FILE_PATH"))
	// 退出时清除缓存
	rootCmd.PersistentFlags().BoolVar(&config.Config.ClearCacheExit, "cache-clean", true, locale.GetString("CACHE_FILE_CLEAN"))
	// 手动指定zip文件编码 gbk、shiftjis……
	rootCmd.PersistentFlags().StringVar(&config.Config.ZipFileTextEncoding, "zip-encode", "gbk", locale.GetString("ZIP_ENCODE"))
}
