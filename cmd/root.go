package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/routers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
comi -lw -webp-command=C:\Users\test\Desktop\webp-server-windows-amd64.exe -p 3344 -q 45  test.zip
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
	//还没做配置文件，暂时屏蔽
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.comicview.yaml)")
	// 局部标签(local flag)，只在直接调用它时运行
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// persistent，任何命令下均可使用，适合全局flag
	rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, "服务端口")
	rootCmd.PersistentFlags().StringVarP(&common.Config.ConfigPath, "config", "c", ".", "配置文件")
	rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", 2, "最大搜索深度")
	rootCmd.PersistentFlags().BoolVarP(&common.Config.OnlyLocal, "local-only", "l", false, "禁用LAN分享")
	rootCmd.PersistentFlags().BoolVarP(&common.Config.UseWebpServer, "webp", "w", false, "webp传输，需要webp-server")
	rootCmd.PersistentFlags().StringVar(&common.Config.WebpCommand, "webp-command", "webp-server", "webp-server命令,或webp-server可执行文件路径，默认为“webp-server")
	rootCmd.PersistentFlags().StringVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", "60", "webp压缩质量（默认60）")
	rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImageInServer, "checkimage", true, "在服务器端分析图片分辨率")
	rootCmd.PersistentFlags().BoolVarP(&common.Config.OpenBrowser, "broswer", "b", false, "同时打开浏览器，windows=true")
	rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "v", false, "输出版本号")
	rootCmd.PersistentFlags().StringVar(&common.Config.ServerHost, "host", "", "自定义域名")
	rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, "记录log文件")
	rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print-allip", false, "打印所有可用网卡ip")
	rootCmd.PersistentFlags().StringVarP(&common.Config.ZipFilenameEncoding, "zip-encoding", "e", "", "Zip non-utf8 Encoding(gbk、shiftjis、gb18030）")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "logname", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "logpath", "~", "log文件位置")
	rootCmd.PersistentFlags().IntVarP(&common.Config.MinImageNum, "imagenum", "i", 3, "至少有几张图片，才认定为漫画压缩包")
	if runtime.GOOS == "windows" {
		common.Config.OpenBrowser = true
	}
}

// initConfig reads in config file and ENV variables if set.
//参考：https://www.loginradius.com/engineering/blog/environment-variables-in-golang/
func initConfig() {

	// Set the path to look for the configurations file
	if common.Config.ConfigPath != "" {
		viper.AddConfigPath(common.Config.ConfigPath)

	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			viper.AddConfigPath(".")
			fmt.Println(err)
		}else {
			viper.AddConfigPath(home)
		}
	}

	viper.SetConfigType("yaml")
	// Set the file name of the configurations file
	viper.SetConfigName(".config/comigo")

	//viper.SetConfigFile(common.Config.ConfigPath+"\/config.yaml")
	//如果不存在，就写入
	err := viper.SafeWriteConfig()
	if err != nil {
		fmt.Println("保存配置:", common.Config.ConfigPath)
	}
	//读取符合的环境变量
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file:", common.Config.ConfigPath)
	}

	//读取案例：
	// Set undefined variables
	viper.SetDefault("COMI.HOST", "0.0.0.0")

	// getting env variables DB.PORT
	// viper.Get() returns an empty interface{}
	// so we have to do the type assertion, to get the value
	DBPort, ok := viper.Get("COMI.PORT").(string)

	// if type assert is not valid it will throw an error
	if !ok {
		//log.Fatalf("Invalid type assertion")
		fmt.Println("Invalid type assertion")
		//os.Exit(0)
	}
	fmt.Printf("viper : %s = %s \n", "Database Port", DBPort)
}
