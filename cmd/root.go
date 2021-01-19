package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers"
	"os"
	"runtime"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  locale.GetString("comigo_use") ,
	Short:locale.GetString("short_description")  ,
	Example: locale.GetString("comigo_example"),
	Version: locale.GetString("comigo_version"),
	Long: locale.GetString("long_description"),
	Run: func(cmd *cobra.Command, args []string) {
		routers.StartServer(args)
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""       //屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 //"这是命令行程序"的提醒表示时间
	//根据配置或系统变量，初始化各种参数
	cobra.OnInitialize(readConfigFile)
	// 局部标签(local flag)，只在直接调用它时运行
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// persistent，任何命令下均可使用，适合全局flag
	//服务端口
	if viper.GetInt("COMI_PORT") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", viper.GetInt("COMI_PORT"), locale.GetString("COMI_PORT"))
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, locale.GetString("COMI_PORT"))
	}
	//指定配置文件
	if viper.GetString("COMI_CONFIG") != "" {
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", viper.GetString("COMI_CONFIG"), locale.GetString("COMI_CONFIG"))
		viper.SetConfigFile(viper.GetString("COMI_CONFIG"))
	} else {
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", ".", locale.GetString("COMI_CONFIG"))
		viper.SetConfigFile(common.Config.ConfigPath)
	}
	//打开浏览器
	if viper.GetBool("COMI_OPEN_BROWSER") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", viper.GetBool("COMI_OPEN_BROWSER"), locale.GetString("COMI_OPEN_BROWSER"))
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", false, locale.GetString("COMI_OPEN_BROWSER"))
	}
	if runtime.GOOS == "windows" {
		common.Config.OpenBrowser = true
	}
	//不对局域网开放
	if viper.GetBool("COMI_DISABLE_LAN") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", viper.GetBool("COMI_DISABLE_LAN"), locale.GetString("COMI_DISABLE_LAN"))
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", false, locale.GetString("COMI_DISABLE_LAN"))
	}
	//文件搜索深度
	if viper.GetInt("COMI_MAX_DEPTH") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", viper.GetInt("COMI_MAX_DEPTH"), locale.GetString("COMI_MAX_DEPTH"))
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", 1, locale.GetString("COMI_MAX_DEPTH"))
	}
	//服务器解析分辨率
	if viper.GetBool("COMI_CHECK_IMAGE") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", viper.GetBool("COMI_CHECK_IMAGE"), locale.GetString("COMI_CHECK_IMAGE"))
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", true, locale.GetString("COMI_CHECK_IMAGE"))
	}
	//本地Host名
	if viper.GetString("COMI_LOCAL_HOST") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", viper.GetString("COMI_LOCAL_HOST"), locale.GetString("COMI_LOCAL_HOST"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", "", locale.GetString("COMI_LOCAL_HOST"))
	}
	//打印所有可用网卡ip
	if viper.GetBool("COMI_PRINT_ALL_IP") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", viper.GetBool("COMI_PRINT_ALL_IP"), locale.GetString("COMI_PRINT_ALL_IP"))
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", false, locale.GetString("COMI_PRINT_ALL_IP"))
	}
	//至少有几张图片，才认定为漫画压缩包
	if viper.GetInt("COMI_MIN_IMAGE_NUM") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", viper.GetInt("COMI_MIN_IMAGE_NUM"), locale.GetString("COMI_MIN_IMAGE_NUM"))
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", 3, locale.GetString("COMI_MIN_IMAGE_NUM"))
	}
	////webp相关
	//启用webp传输
	if viper.GetBool("COMI_ENABLE_WEBP") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", viper.GetBool("COMI_ENABLE_WEBP"), locale.GetString("COMI_ENABLE_WEBP"))
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", false, locale.GetString("COMI_ENABLE_WEBP"))
	}
	//webp-server命令
	if viper.GetString("COMI_WEBP_COMMAND") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", viper.GetString("COMI_WEBP_COMMAND"), locale.GetString("COMI_WEBP_COMMAND"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", "webp-server", locale.GetString("COMI_WEBP_COMMAND"))
	}
	//webp压缩质量
	if viper.GetInt("COMI_WEBP_QUALITY") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", viper.GetInt("COMI_WEBP_QUALITY"), locale.GetString("COMI_WEBP_QUALITY"))
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", 95, locale.GetString("COMI_WEBP_QUALITY"))
	}
	////Frpc相关
	//frp反向代理
	if viper.GetBool("COMI_ENABLE_FRPC") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", viper.GetBool("COMI_ENABLE_FRPC"), locale.GetString("COMI_ENABLE_FRPC"))
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", false, locale.GetString("COMI_ENABLE_FRPC"))
	}
	//frps_addr
	if viper.GetString("COMI_FRP_SERVER_ADDR") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", viper.GetString("COMI_FRP_SERVER_ADDR"), locale.GetString("COMI_FRP_SERVER_ADDR"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", "frps.example.com", locale.GetString("COMI_FRP_SERVER_ADDR"))
	}
	//frps server_port
	if viper.GetInt("COMI_FRP_SERVER_PORT") != 0 {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", viper.GetInt("COMI_FRP_SERVER_PORT"), locale.GetString("COMI_FRP_SERVER_PORT"))
	} else {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", 7000, locale.GetString("COMI_FRP_SERVER_PORT"))
	}
	//frp token
	if viper.GetString("COMI_FRP_TOKEN") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", viper.GetString("COMI_FRP_TOKEN"), locale.GetString("COMI_FRP_TOKEN"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", "", locale.GetString("COMI_FRP_TOKEN"))
	}
	//frpc命令,或frpc可执行文件路径
	if viper.GetString("COMI_FRP_COMMAND") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", viper.GetString("COMI_FRP_COMMAND"), locale.GetString("COMI_FRP_COMMAND"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", "frpc", locale.GetString("COMI_FRP_COMMAND"))
	}
	//frpc random remote_port
	if viper.GetBool("COMI_FRP_RANDOM_REMOTE_PORT") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "random_remote_port", viper.GetBool("COMI_FRP_RANDOM_REMOTE_PORT"), locale.GetString("COMI_FRP_RANDOM_REMOTE_PORT"))
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "random_remote_port", true, locale.GetString("COMI_FRP_RANDOM_REMOTE_PORT"))
	}
	//frpc remote_port
	if viper.GetInt("COMI_FRP_REMOTE_PORT") != 0 {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port", viper.GetInt("COMI_FRP_REMOTE_PORT"), locale.GetString("COMI_FRP_REMOTE_PORT"))
	} else {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port", -1, locale.GetString("COMI_FRP_REMOTE_PORT"))
	}
	//输出log文件
	if viper.GetBool("COMI_LOG_TO_FILE") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", viper.GetBool("COMI_LOG_TO_FILE"), locale.GetString("COMI_LOG_TO_FILE"))
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, locale.GetString("COMI_LOG_TO_FILE"))
	}

	//默认页面模式
	if viper.GetBool("COMI_TEMPLATE") {
		rootCmd.PersistentFlags().StringVar(&common.Config.Template, "template", viper.GetString("COMI_TEMPLATE"), locale.GetString("COMI_TEMPLATE"))
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.Template, "template", "auto", locale.GetString("COMI_TEMPLATE"))
	}

	//尚未启用的功能，暂时无意义的设置
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "logname", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "logpath", "~", "log文件位置")
	//rootCmd.PersistentFlags().StringVarP(&common.Config.ZipFilenameEncoding, "zip-encoding", "e", "", "Zip non-utf8 Encoding(gbk、shiftjis、gb18030）")
	//	rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "v", false, "输出版本号")
}

// 读取配置，参考下面三篇文章
//1、https://www.loginradius.com/engineering/blog/environment-variables-in-golang/
//2、https://www.liwenzhou.com/posts/Go/viper_tutorial/
//3、https://ovh.github.io/tat/sdk/golang-full-example/
func readConfigFile() {
	//Viper优先级顺序： 显式调用 Set 函数 > 命令行参数 > 环境变量 > 配置文件 > 远程 key/value 存储系统 > 默认值
	//读取环境变量
	viper.AutomaticEnv()
	//viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("comigo")      // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")           // 在当前目录中查找配置
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		viper.AddConfigPath(home) // 在Home目录中查找配置
	}
	viper.AddConfigPath(home + "/.config/comigo") // 多次调用以添加多个搜索路径,home
	viper.AddConfigPath("./.config")          // 多次调用以添加多个搜索路径
	//查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Debugf(locale.GetString("config_file_not_found") , viper.ConfigFileUsed())

		} else {
			logrus.Debugf(locale.GetString("config_file_not_resolve") , common.Config.ConfigPath)
		}
	}
	//应用配置文件
	if err := viper.Unmarshal(&common.Config); err != nil { // 读取配置文件转化成对应的结构体错误
		panic(fmt.Errorf(locale.GetString("config_file_not_found")+" %s \n", err))
	}
	//// 设置默认值
	//viper.SetDefault("COMI_HOST", "0.0.0.0")
	//监听配置变化，运行时动态加载配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(locale.GetString("config_change") , e.Name)
	})
	//保存配置並退出：
	////将当前的viper配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
	//err = viper.SafeWriteConfigAs("sample_config.yaml")
	//if err != nil {
	//	fmt.Println(locale.GetString("save_config_failed") ,err.Error())
	//} else {
	//	fmt.Println(locale.GetString("save_config_file"), common.Config.ConfigPath)
	//}
}
