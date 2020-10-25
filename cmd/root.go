package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"

	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/routers"

	"github.com/spf13/cobra"
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

comi -w -q 50 book.zip  

指定多个参数：
comi -w -q 50 --frpc  --token aX4457d3O -p 23455 --frps-addr sh.example.com test.zip
`,
	Version: "v0.2.3",
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""       //屏蔽鼠标提示，支持拖拽、双击运行
	cobra.MousetrapDisplayDuration = 5 //"这是命令行程序"的提醒表示时间
	//根据配置或系统变量，初始化各种参数
	cobra.OnInitialize(initConfig)
	viper.AutomaticEnv()
	// 局部标签(local flag)，只在直接调用它时运行
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// persistent，任何命令下均可使用，适合全局flag

	//服务端口
	if viper.GetInt("COMI_PORT")!=0{
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", viper.GetInt("COMI_PORT"), "服务端口")
	}else{
		rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, "服务端口")
	}

	//配置文件位置
	if viper.GetString("COMI_CONFIG")!=""{
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", viper.GetString(""), "配置文件")
	}else{
		rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", ".", "配置文件")
	}

	//文件搜索深度
	if viper.GetInt("COMI_MAX_DEPTH")!=0{
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", viper.GetInt("COMI_MAX_DEPTH"), "最大搜索深度")
	}else{
		rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", 1, "最大搜索深度")
	}

	//不对局域网开放
	if viper.GetBool("COMI_DISABLE_LAN"){
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", viper.GetBool("COMI_DISABLE_LAN"), "禁用LAN分享")
	}else{
		rootCmd.PersistentFlags().BoolVarP(&common.Config.DisableLAN, "disable-lan", "d", false, "禁用LAN分享")
	}
	//服务器解析分辨率
	if viper.GetBool("COMI_CHECK_IMAGE"){
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", viper.GetBool("COMI_CHECK_IMAGE"), "在服务器端分析图片分辨率")
	}else{
		rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", true, "在服务器端分析图片分辨率")
	}

	if viper.GetString("COMI_LOCAL_HOST")!=""{
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", viper.GetString("COMI_LOCAL_HOST"), "自定义域名")
	}else{
		rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "local_host", "", "自定义域名")
	}
	//打印所有可用网卡ip
	if viper.GetBool("COMI_PRINT_ALL_IP"){
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", viper.GetBool("COMI_PRINT_ALL_IP"), "打印所有可用网卡ip")
	}else{
		rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print_all_ip", false, "打印所有可用网卡ip")
	}
	//至少有几张图片，才认定为漫画压缩包
	if viper.GetInt("COMI_MIN_IMAGE_NUM")!=0{
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", viper.GetInt("COMI_MIN_IMAGE_NUM"), "至少有几张图片，才认定为漫画压缩包")
	}else{
		rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "min-image-num", "i", 3, "至少有几张图片，才认定为漫画压缩包")
	}

	////webp相关
	//启用webp传输
	if viper.GetBool("COMI_ENABLE_WEBP"){
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", viper.GetBool("COMI_ENABLE_WEBP"), "webp传输，需要webp-server")
	}else{
		rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", false, "启用webp传输，需要webp-server")
	}
	//webp-server命令
	if viper.GetString("COMI_WEBP_COMMAND")!=""{
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", viper.GetString("COMI_WEBP_COMMAND"), "webp-server命令,或webp-server可执行文件路径，默认为“webp-server")
	}else{
		rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", "webp-server", "webp-server命令,或webp-server可执行文件路径，默认为“webp-server")
	}
	//webp压缩质量
	if viper.GetInt("COMI_WEBP_QUALITY")!=0{
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", viper.GetInt("COMI_WEBP_QUALITY"), "webp压缩质量（默认60）")
	}else{
		rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", 60, "webp压缩质量（默认60）")
	}

	////Frpc相关
	//启用frp反向代理
	if viper.GetBool("COMI_ENABLE_FRPC"){
		rootCmd.PersistentFlags().BoolVarP(&common.Config.UseFrpc, "frpc", "f", viper.GetBool("COMI_ENABLE_FRPC"), "启用frp反向代理")
	}else{
		rootCmd.PersistentFlags().BoolVarP(&common.Config.UseFrpc, "frpc", "f", false, "启用frp反向代理")
	}
	//frps_addr
	if viper.GetString("COMI_FRP_SERVER_ADDR")!=""{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr",  viper.GetString("COMI_FRP_SERVER_ADDR"), "frps_addr, frpc必须")
	}else{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr",  "frps.example.com", "frps_addr, frpc必须")
	}
	//frps server_port
	if viper.GetInt("COMI_FRP_SERVER_PORT")!=0{
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port",  viper.GetInt("COMI_FRP_SERVER_PORT"), "frps server_port, frpc必须")
	}else{
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port",  7000, "frps server_port, frpc必须")
	}
	//frp token
	if viper.GetString("COMI_FRP_TOKEN")!=""{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token",  viper.GetString("COMI_FRP_TOKEN"), "token, frpc必须")
	}else{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token",  "", "token, frpc必须")
	}
	//frpc命令,或frpc可执行文件路径
	if viper.GetString("COMI_FRP_COMMAND")!=""{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", viper.GetString("COMI_FRP_COMMAND"), "frpc命令,或frpc可执行文件路径，默认为“frpc")
	}else{
		rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", "frpc", "frpc命令,或frpc可执行文件路径，默认为“frpc")
	}
	//frpc remote_port
	if viper.GetInt("COMI_FRP_REMOTE_PORT")!=0{
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port",  viper.GetInt("COMI_FRP_REMOTE_PORT"), "frpc remote_port，默认与本地相同")
	}else{
		rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "remote_port",  65536, "frpc remote_port，默认与本地相同")
	}

	//打开浏览器相关
	if viper.GetBool("COMI_OPEN_BROWSER"){
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", viper.GetBool("COMI_OPEN_BROWSER"), "同时打开浏览器，windows=true")
	}else{
		rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "browser", "b", false, "同时打开浏览器，windows=true")
	}
	if runtime.GOOS == "windows" {
		common.Config.OpenBrowser = true
	}

	//尚未启用的功能，暂时无意义的设置
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "logname", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "logpath", "~", "log文件位置")
	//rootCmd.PersistentFlags().StringVarP(&common.Config.ZipFilenameEncoding, "zip-encoding", "e", "", "Zip non-utf8 Encoding(gbk、shiftjis、gb18030）")
	//	rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "v", false, "输出版本号")
	//还没做配置文件，暂时屏蔽
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.comicview.yaml)")
	if viper.GetBool("COMI_LOG.TO.FILE"){
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", viper.GetBool("COMI_LOG.TO.FILE"), "记录log文件")
	}else{
		rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, "记录log文件")
	}
}

// initConfig reads in config file and ENV variables if set.
//参考：https://www.loginradius.com/engineering/blog/environment-variables-in-golang/
func initConfig() {

	// // Set the path to look for the configurations file
	// if common.Config.ConfigPath != "" {
	// 	viper.AddConfigPath(common.Config.ConfigPath)

	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		viper.AddConfigPath(".")
	// 		fmt.Println(err)
	// 	}else {
	// 		viper.AddConfigPath(home)
	// 	}
	// }

	// viper.SetConfigType("yaml")
	// // Set the file name of the configurations file
	// viper.SetConfigName(".config/comigo")

	// //viper.SetConfigFile(common.Config.ConfigPath+"\/config.yaml")
	// //如果不存在，就写入
	// err := viper.SafeWriteConfig()
	// if err != nil {
	// 	fmt.Println("保存配置:", common.Config.ConfigPath)
	// }
	// //读取符合的环境变量
	// viper.AutomaticEnv() // read in environment variables that match
	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// } else {
	// 	fmt.Println("No config file:", common.Config.ConfigPath)
	// }

	// //读取案例：
	// // Set undefined variables
	// viper.SetDefault("COMI.HOST", "0.0.0.0")

	// // getting env variables DB.PORT
	// // viper.Get() returns an empty interface{}
	// // so we have to do the type assertion, to get the value
	// DBPort, ok := viper.Get("COMI.PORT").(string)

	// // if type assert is not valid it will throw an error
	// if !ok {
	// 	//log.Fatalf("Invalid type assertion")
	// 	fmt.Println("Invalid type assertion")
	// 	//os.Exit(0)
	// }
	// fmt.Printf("viper : %s = %s \n", "Database Port", DBPort)
}
