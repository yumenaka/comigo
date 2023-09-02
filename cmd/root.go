package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/routers/handler"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers"
)

var runtimeViper *viper.Viper

func init() {
	cobra.MousetrapHelpText = ""       //屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 //"这是命令行程序"的提醒表示时间
	//jwt认证，tls证书
	rootCmd.PersistentFlags().StringVarP(&common.Config.Username, "username", "u", "comigo", locale.GetString("USERNAME"))
	rootCmd.PersistentFlags().StringVarP(&common.Config.Password, "password", "k", "", locale.GetString("PASSWORD"))
	rootCmd.PersistentFlags().IntVarP(&common.Config.Timeout, "timeout", "t", 65535, locale.GetString("TIMEOUT"))
	//TLS设定
	rootCmd.PersistentFlags().StringVar(&common.Config.CertFile, "tls-cert-file", "", locale.GetString("TLS_CERT_FILE"))              //--tls-cert-file
	rootCmd.PersistentFlags().StringVar(&common.Config.KeyFile, "tls-private-key-file", "", locale.GetString("TLS_PRIVATE_KEY_FILE")) //--tls-private-key-file
	//指定配置文件
	rootCmd.PersistentFlags().StringVarP(&common.ConfigFilePath, "config", "c", "", locale.GetString("CONFIG"))
	//在当前目录生成示例配置文件
	rootCmd.PersistentFlags().BoolVar(&common.GenerateConfig, "generate-config", false, locale.GetString("GenerateConfig"))
	//启用数据库，保存扫描数据
	rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableDatabase, "enable-database", "e", false, locale.GetString("EnableDatabase"))
	//服务端口
	rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, locale.GetString("PORT"))
	//本地Host
	rootCmd.PersistentFlags().StringVar(&common.Config.Host, "host", "OutboundIP", locale.GetString("LOCAL_HOST"))
	//DEBUG
	rootCmd.PersistentFlags().BoolVar(&common.Config.Debug, "debug", false, locale.GetString("DEBUG_MODE"))
	//打开浏览器
	rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "open-browser", "o", false, locale.GetString("OPEN_BROWSER"))
	if runtime.GOOS == "windows" {
		common.Config.OpenBrowser = true
	}
	//不对局域网开放
	rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", false, locale.GetString("DISABLE_LAN"))
	//文件搜索深度
	rootCmd.PersistentFlags().IntVarP(&common.Config.MaxScanDepth, "max-depth", "m", 3, locale.GetString("MAX_DEPTH"))
	////服务器解析书籍元数据，如果生成blurhash，需要消耗大量资源
	rootCmd.PersistentFlags().BoolVar(&common.Config.GenerateMetaData, "generate-metadata", false, locale.GetString("GENERATE_METADATA"))
	//打印所有可用网卡ip
	rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllPossibleQRCode, "print-all", false, locale.GetString("PRINT_ALL_IP"))
	//至少有几张图片，才认定为漫画压缩包
	rootCmd.PersistentFlags().IntVar(&common.Config.MinImageNum, "min-image-num", 1, locale.GetString("MIN_MEDIA_NUM"))
	////webp相关 拆分成子命令？
	//启用webp传输
	//rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", false, locale.GetString("ENABLE_WEBP"))
	//webp-server命令
	//rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", "webp-server", locale.GetString("WEBP_COMMAND"))
	//webp压缩质量
	//rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", 85, locale.GetString("WEBP_QUALITY"))
	////Frpc相关  拆分成子命令？
	//frp反向代理
	rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", false, locale.GetString("ENABLE_FRPC"))
	//frps_addr
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", "frps.example.com", locale.GetString("FRP_SERVER_ADDR"))
	//frps server_port
	rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", 7000, locale.GetString("FRP_SERVER_PORT"))
	//frp token
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", "token_secret_sample", locale.GetString("FRP_TOKEN"))
	//frpc命令,或frpc可执行文件路径
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", "frpc", locale.GetString("FRP_COMMAND"))
	//frpc random remote_port
	rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "frps-random-remote", true, locale.GetString("FRP_RANDOM_REMOTE_PORT"))
	//frpc remote_port
	rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "frps-remote-port", 50000, locale.GetString("FRP_REMOTE_PORT"))
	//输出log文件
	rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, locale.GetString("LOG_TO_FILE"))
	//sketch模式的倒计时秒数
	//rootCmd.PersistentFlags().IntVar(&common.Config.SketchCountSeconds, "sketch_count_seconds", 90, locale.GetString("SKETCH_COUNT_SECONDS"))

	rootCmd.PersistentFlags().BoolVar(&common.Config.EnableLocalCache, "enable-cache", true, locale.GetString("CACHE_FILE_ENABLE"))
	//web图片缓存路径
	rootCmd.PersistentFlags().StringVar(&common.Config.CachePath, "cache-path", "", locale.GetString("CACHE_FILE_PATH"))
	//退出时清除临时文件
	rootCmd.PersistentFlags().BoolVar(&common.Config.ClearCacheExit, "cache-clean", true, locale.GetString("CACHE_FILE_CLEAN"))

	//启用文件上传功能
	rootCmd.PersistentFlags().BoolVar(&common.Config.EnableUpload, "enable-upload", true, locale.GetString("ENABLE_FILE_UPLOAD"))
	//上传文件的保存路径
	rootCmd.PersistentFlags().StringVar(&common.Config.UploadPath, "upload-path", "", locale.GetString("UPLOAD_PATH"))
	if common.Config.EnableUpload && common.Config.UploadPath == "" {
		//获取当前目录
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			common.Config.UploadPath = path.Join(dir, "upload")
		}
	}
	handler.EnableUpload = &common.Config.EnableUpload
	handler.UploadPath = &common.Config.UploadPath

	//手动指定zip文件编码 gbk、shiftjis……
	rootCmd.PersistentFlags().StringVar(&common.Config.ZipFileTextEncoding, "zip-encode", "gbk", locale.GetString("ZIP_ENCODE"))
	//尚未写完的功能
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "log_name", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "log_path", "~", "log文件位置")
	//rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "runtimeViper", false, "输出版本号")
	//cobra & viper sample:https://qiita.com/nirasan/items/cc2ab5bc2889401fe596
	// rootCmd.Run() 运行前的初始化定义。
	// 运行前后顺序：rootCmd.Execute → 命令行参数的处理 → cobra.OnInitialize → rootCmd.Run、
	// 于是可以通过CMD读取配置文件、按照配置文件的设定值执行。不一致的时候，配置文件优先于CMD参数
	//cobra.OnInitialize(initConfig)
	cobra.OnInitialize(func() {
		runtimeViper = viper.New()
		//自动读取环境变量，改写对应值
		runtimeViper.AutomaticEnv()
		//设置环境变量的前缀，将 PORT变为 COMI_PORT
		runtimeViper.SetEnvPrefix("COMI")
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		}
		//需要在home目录下面搜索配置文件
		homeConfigPath := path.Join(home, ".config/comigo")
		runtimeViper.AddConfigPath(homeConfigPath)

		// 获取可执行程序自身的文件路径
		executablePath, err := os.Executable()
		if err != nil {
			fmt.Println("无法获取程序路径:", err)
			return
		}

		// 将可执行程序自身的文件路径转换为绝对路径
		absPath, err := filepath.Abs(executablePath)
		if err != nil {
			fmt.Println("Failed to get absolute path:", err)
			return
		}
		fmt.Println("Executable path:", absPath)
		runtimeViper.AddConfigPath(absPath)

		// 当前执行目录
		nowPath, err := os.Getwd()
		if err != nil {
			fmt.Println("无法获取程序执行目录:", err)
		}
		runtimeViper.AddConfigPath(nowPath)
		runtimeViper.SetConfigType("toml")
		runtimeViper.SetConfigName("config.toml")
		// 读取设定文件
		if err := runtimeViper.ReadInConfig(); err != nil {
			if common.ConfigFilePath == "" && common.Config.Debug {
				fmt.Println(err)
			}
		} else {
			//获取当前使用的配置文件路径
			//https://github.com/spf13/viper/issues/89
			common.ConfigFilePath = runtimeViper.ConfigFileUsed()
			fmt.Println(locale.GetString("FoundConfigFile") + common.ConfigFilePath)
		}

		// 把设定文件的内容，解析到构造体里面。
		if err := runtimeViper.Unmarshal(&common.Config); err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			os.Exit(1)
		}
		//保存配置示例並退出
		if common.GenerateConfig {
			common.Config.LogFilePath = ""
			common.Config.EnableDatabase = true
			common.Config.OpenBrowser = false
			common.Config.StoresPath = []string{"C:\\test\\Comic", "D:\\some_path\\book", "/home/username/Download"}
			common.Config.CachePath = ".comigo"
			bytes, err := toml.Marshal(common.Config)
			if err != nil {
				fmt.Println("toml.Marshal Error")
			}
			//在命令行打印
			fmt.Println(string(bytes))
			err = os.WriteFile("config.toml", bytes, 0644)
			if err != nil {
				panic(err)
			}
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}
		//监听文件修改
		runtimeViper.WatchConfig()
		//文件修改时，执行重载设置、服务重启的函数
		runtimeViper.OnConfigChange(configReloadHandler)
	})
}

// rootCmd 没有任何子命令的情况下调用时的基本命令
var rootCmd = &cobra.Command{
	Use:     locale.GetString("comigo_use"),
	Short:   locale.GetString("short_description"),
	Example: locale.GetString("comigo_example"),
	Version: common.Version,
	Long:    locale.GetString("long_description"),
	Run: func(cmd *cobra.Command, args []string) {
		//解析命令，扫描文件
		initBookStores(args)
		//设置临时文件夹
		common.SetTempDir()
		//CheckWebPort
		routers.CheckWebPort()
		//设置书籍API
		routers.StartWebServer()
		//退出时清理临时文件
		SetShutdownHandler()
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}
