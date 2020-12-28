package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/routers"
)

//var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "comi",
	Short: "A simple comic book reader.",
	Example: `
comi book.zip

设定网页服务端口（默认为1234）：
comi -p 2345 book.zip 

不打开浏览器（windows）：
comi -b=false book.zip 

本机浏览，不对外开放：
comi -l book.zip  

webp传输，需要webp-server配合：
comi -w book.zip  

指定多个参数：
comi -w -q 70 --frpc  --token aX4457d3O -p 23455 --frps-addr sh.example.com test.zip
`,
	Version: "v0.2.4",
	Long: `comigo 一款简单的漫画阅读器
`,
	Run: func(cmd *cobra.Command, args []string) {
		routers.StartComicServer(args)
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
	//考虑试下termui
}

// 读取配置，参考下面三篇文章
//1、https://www.loginradius.com/engineering/blog/environment-variables-in-golang/
//2、https://www.liwenzhou.com/posts/Go/viper_tutorial/
//3、https://ovh.github.io/tat/sdk/golang-full-example/
func initConfig() {
	//Viper优先级顺序： 显式调用 Set 函数 > 命令行参数 > 环境变量 > 配置文件 > 远程 key/value 存储系统 > 默认值
	//读取环境变量
	viper.AutomaticEnv()
	//viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")      // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/etc/comigo") // 查找配置文件所在的路径，可以使用相对路径，也可以使用绝对路径
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(home+"/.config/comigo") // 多次调用以添加多个搜索路径
	viper.AddConfigPath("./.config/comigo")     // 多次调用以添加多个搜索路径
	viper.AddConfigPath("./.comigo")     // 多次调用以添加多个搜索路径
	viper.AddConfigPath("./.config")            // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")                    // 还可以在工作目录中查找配置
	//查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Debugf("配置文件未找到", viper.ConfigFileUsed())
		} else {
			logrus.Debugf("解析配置失败r", common.Config.ConfigPath)
		}
	}
	//应用配置文件
	if err := viper.Unmarshal(&common.Config); err != nil { // 读取配置文件转化成对应的结构体错误
		panic(fmt.Errorf("read config file to struct err: %s \n", err))
	}
	//// 设置默认值
	//viper.SetDefault("COMI_HOST", "0.0.0.0")
	//将当前的viper配置写入预定义的路径。如果没有预定义的路径，则报错。如果存在，将不会覆盖当前的配置文件。
	err = viper.SafeWriteConfig()
	if err != nil {
		fmt.Println("保存配置:", common.Config.ConfigPath)
	} else {
		fmt.Println("保存失败:", err.Error())
	}
	//监听配置变化，运行时动态加载配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("comigo配置发生变更：", e.Name)
	})
}

func init() {
	cobra.MousetrapHelpText = ""       //屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 //"这是命令行程序"的提醒表示时间
	//根据配置或系统变量，初始化各种参数
	cobra.OnInitialize(initConfig)
	// 局部标签(local flag)，只在直接调用它时运行
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// persistent，任何命令下均可使用，适合全局flag
	//服务端口
	if viper.GetInt("COMI_PORT") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", viper.GetInt("COMI_PORT"), "服务端口")
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, "服务端口")
	}

	//指定配置文件
	if viper.GetString("COMI_CONFIG") != "" {
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", viper.GetString("COMI_CONFIG"), "指定配置文件")
		viper.SetConfigFile(viper.GetString("COMI_CONFIG"))
	} else {
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", ".", "指定配置文件")
		viper.SetConfigFile(common.Config.ConfigPath)
	}
	//打开浏览器
	if viper.GetBool("COMI_OPEN_BROWSER") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", viper.GetBool("COMI_OPEN_BROWSER"), "同时打开浏览器，windows=true")
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", false, "同时打开浏览器，windows=true")
	}
	if runtime.GOOS == "windows" {
		common.Config.OpenBrowser = true
	}
	//不对局域网开放
	if viper.GetBool("COMI_DISABLE_LAN") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", viper.GetBool("COMI_DISABLE_LAN"), "禁用LAN分享")
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", false, "禁用LAN分享")
	}

	//文件搜索深度
	if viper.GetInt("COMI_MAX_DEPTH") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", viper.GetInt("COMI_MAX_DEPTH"), "最大搜索深度")
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", 1, "最大搜索深度")
	}
	//服务器解析分辨率
	if viper.GetBool("COMI_CHECK_IMAGE") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", viper.GetBool("COMI_CHECK_IMAGE"), "在服务器端分析图片分辨率")
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", true, "在服务器端分析图片分辨率")
	}
	//本地Host名
	if viper.GetString("COMI_LOCAL_HOST") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", viper.GetString("COMI_LOCAL_HOST"), "自定义域名")
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", "", "自定义域名")
	}
	//打印所有可用网卡ip
	if viper.GetBool("COMI_PRINT_ALL_IP") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", viper.GetBool("COMI_PRINT_ALL_IP"), "打印所有可用网卡ip")
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", false, "打印所有可用网卡ip")
	}
	//至少有几张图片，才认定为漫画压缩包
	if viper.GetInt("COMI_MIN_IMAGE_NUM") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", viper.GetInt("COMI_MIN_IMAGE_NUM"), "至少有几张图片，才认定为漫画压缩包")
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", 3, "至少有几张图片，才认定为漫画压缩包")
	}
	////webp相关
	//启用webp传输
	if viper.GetBool("COMI_ENABLE_WEBP") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", viper.GetBool("COMI_ENABLE_WEBP"), "webp传输，需要webp-server")
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", false, "启用webp传输，需要webp-server")
	}
	//webp-server命令
	if viper.GetString("COMI_WEBP_COMMAND") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", viper.GetString("COMI_WEBP_COMMAND"), "webp-server命令,或webp-server可执行文件路径，默认为“webp-server")
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", "webp-server", "webp-server命令,或webp-server可执行文件路径，默认为“webp-server")
	}
	//webp压缩质量
	if viper.GetInt("COMI_WEBP_QUALITY") != 0 {
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", viper.GetInt("COMI_WEBP_QUALITY"), "webp压缩质量（默认60）")
	} else {
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", 60, "webp压缩质量（默认60）")
	}
	////Frpc相关
	//启用frp反向代理
	if viper.GetBool("COMI_ENABLE_FRPC") {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", viper.GetBool("COMI_ENABLE_FRPC"), "启用frp反向代理")
	} else {
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", false, "启用frp反向代理")
	}
	//frps_addr
	if viper.GetString("COMI_FRP_SERVER_ADDR") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", viper.GetString("COMI_FRP_SERVER_ADDR"), "frps_addr, frpc必须")
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", "frps.example.com", "frps_addr, frpc必须")
	}
	//frps server_port
	if viper.GetInt("COMI_FRP_SERVER_PORT") != 0 {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", viper.GetInt("COMI_FRP_SERVER_PORT"), "frps server_port, frpc必须")
	} else {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", 7000, "frps server_port, frpc必须")
	}
	//frp token
	if viper.GetString("COMI_FRP_TOKEN") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", viper.GetString("COMI_FRP_TOKEN"), "token, frpc必须")
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", "", "token, frpc必须")
	}
	//frpc命令,或frpc可执行文件路径
	if viper.GetString("COMI_FRP_COMMAND") != "" {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", viper.GetString("COMI_FRP_COMMAND"), "frpc命令,或frpc可执行文件路径，默认为“frpc")
	} else {
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", "frpc", "frpc命令,或frpc可执行文件路径，默认为“frpc")
	}
	//frpc random remote_port
	if viper.GetBool("COMI_FRP_RANDOM_REMOTE_PORT") {
		rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "random_remote_port", viper.GetBool("COMI_FRP_RANDOM_REMOTE_PORT"), "frpc 远程随机端口，默认启用（40000~50000）")
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "random_remote_port", true, "frpc 远程随机端口，默认启用（40000~50000）")
	}
	//frpc remote_port
	if viper.GetInt("COMI_FRP_REMOTE_PORT") != 0 {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port", viper.GetInt("COMI_FRP_REMOTE_PORT"), "frpc remote_port，默认与本地相同")
	} else {
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port", -1, "frpc remote_port，只关闭随机功能，则frp远程端口与本地相同")
	}
	//尚未启用的功能，暂时无意义的设置
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "logname", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "logpath", "~", "log文件位置")
	//rootCmd.PersistentFlags().StringVarP(&common.Config.ZipFilenameEncoding, "zip-encoding", "e", "", "Zip non-utf8 Encoding(gbk、shiftjis、gb18030）")
	//	rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "v", false, "输出版本号")

	if viper.GetBool("COMI_LOG.TO.FILE"){
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", viper.GetBool("COMI_LOG.TO.FILE"), "记录log文件")
	} else {
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, "记录log文件")
	}
}
