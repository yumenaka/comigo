package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yumenaka/comi/common"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/routers"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var vip *viper.Viper

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     locale.GetString("comigo_use"),
	Short:   locale.GetString("short_description"),
	Example: locale.GetString("comigo_example"),
	Version: locale.GetString("comigo_version"),
	Long:    locale.GetString("long_description"),
	Run: func(cmd *cobra.Command, args []string) {
		routers.ParseCommands(args)
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

	//简单认证
	rootCmd.PersistentFlags().StringVarP(&common.Config.UserName, "username", "u", "admin", "用户名")
	rootCmd.PersistentFlags().StringVarP(&common.Config.Password, "password", "k", "", "密码")

	//简单认证
	rootCmd.PersistentFlags().StringVar(&common.Config.CertFile, "certfile", "", "tls certfile")
	rootCmd.PersistentFlags().StringVar(&common.Config.KeyFile, "keyfile", "", "tls keyfile")

	//指定配置文件
	rootCmd.PersistentFlags().StringVarP(&common.ConfigFile, "config", "c", "", locale.GetString("CONFIG"))
	//在当前目录生成示例配置文件
	rootCmd.PersistentFlags().BoolVar(&common.Config.NewConfig, "new-config", false, locale.GetString("NewConfig"))
	//服务端口
	rootCmd.PersistentFlags().IntVarP(&common.Config.Port, "port", "p", 1234, locale.GetString("PORT"))
	//本地Host名
	rootCmd.PersistentFlags().StringVar(&common.Config.Host, "host", "", locale.GetString("LOCAL_HOST"))

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
	rootCmd.PersistentFlags().IntVarP(&common.Config.MaxDepth, "max-depth", "m", 2, locale.GetString("MAX_DEPTH"))
	////服务器解析分辨率
	rootCmd.PersistentFlags().BoolVar(&common.Config.CheckImage, "check-image", false, locale.GetString("CHECK_IMAGE"))
	//打印所有可用网卡ip
	rootCmd.PersistentFlags().BoolVar(&common.Config.PrintAllIP, "print-all", false, locale.GetString("PRINT_ALL_IP"))
	//至少有几张图片，才认定为漫画压缩包
	rootCmd.PersistentFlags().IntVar(&common.Config.MinImageNum, "min-image-num", 1, locale.GetString("MIN_MEDIA_NUM"))

	////webp相关 打算拆分成子命令
	//启用webp传输
	//rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableWebpServer, "webp", "w", false, locale.GetString("ENABLE_WEBP"))
	//webp-server命令
	//rootCmd.PersistentFlags().StringVar(&common.Config.WebpConfig.WebpCommand, "webp-command", "webp-server", locale.GetString("WEBP_COMMAND"))
	//webp压缩质量
	//rootCmd.PersistentFlags().IntVarP(&common.Config.WebpConfig.QUALITY, "webp-quality", "q", 85, locale.GetString("WEBP_QUALITY"))

	////Frpc相关  打算拆分成子命令
	//frp反向代理
	rootCmd.PersistentFlags().BoolVarP(&common.Config.EnableFrpcServer, "frpc", "f", false, locale.GetString("ENABLE_FRPC"))
	//frps_addr
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.ServerAddr, "frps-addr", "frps.example.com", locale.GetString("FRP_SERVER_ADDR"))
	//frps server_port
	rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.ServerPort, "frps-port", 7000, locale.GetString("FRP_SERVER_PORT"))
	//frp token
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.Token, "token", "token_secretSAMPLE", locale.GetString("FRP_TOKEN"))
	//frpc命令,或frpc可执行文件路径
	rootCmd.PersistentFlags().StringVar(&common.Config.FrpConfig.FrpcCommand, "frpc-command", "frpc", locale.GetString("FRP_COMMAND"))
	//frpc random remote_port
	rootCmd.PersistentFlags().BoolVar(&common.Config.FrpConfig.RandomRemotePort, "frps-random-remote", true, locale.GetString("FRP_RANDOM_REMOTE_PORT"))
	//frpc remote_port
	rootCmd.PersistentFlags().IntVar(&common.Config.FrpConfig.RemotePort, "frps-remote-port", 50000, locale.GetString("FRP_REMOTE_PORT"))
	//输出log文件
	rootCmd.PersistentFlags().BoolVar(&common.Config.LogToFile, "log", false, locale.GetString("LOG_TO_FILE"))
	//默认web模板
	//rootCmd.PersistentFlags().StringVarP(&common.Config.Template, "template", "t", "scroll", locale.GetString("TEMPLATE"))
	//sketch模式的倒计时秒数
	rootCmd.PersistentFlags().IntVar(&common.Config.SketchCountSeconds, "sketch_count_seconds", 90, locale.GetString("SKETCH_COUNT_SECONDS"))
	//图片文件排序
	rootCmd.PersistentFlags().StringVarP(&common.Config.SortImage, "sort", "s", "none", locale.GetString("SORT"))
	//临时图片解压路径
	rootCmd.PersistentFlags().StringVar(&common.Config.TempPATH, "temp-path", "", locale.GetString("TEMP_PATH"))
	//退出时清除临时文件
	rootCmd.PersistentFlags().BoolVar(&common.Config.CleanAllTempFileOnExit, "clean", false, locale.GetString("CLEAN_ALL_TEMP_FILE"))
	//手动指定zip文件编码(gbk、shiftjis……etc）
	//rootCmd.PersistentFlags().StringVar(&common.Config.ZipFilenameEncoding, "zip-encode", "", locale.GetString("ZIP_ENCODE"))
	////访问密码，还没做
	//	rootCmd.PersistentFlags().StringVar(&common.Config.Auth, "auth", "user:comigo", locale.GetString("AUTH"))
	//尚未写完的功能
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFileName, "log_name", "comigo", "log文件名")
	//rootCmd.PersistentFlags().StringVar(&common.Config.LogFilePath, "log_path", "~", "log文件位置")
	//rootCmd.PersistentFlags().BoolVarP(&common.PrintVersion, "version", "vip", false, "输出版本号")

	//sample:https://qiita.com/nirasan/items/cc2ab5bc2889401fe596

	// rootCmd.Run() 运行前的初始化定义。
	// 运行前后顺序：rootCmd.Execute → 命令行参数的处理 → cobra.OnInitialize → rootCmd.Run、
	// 于是可以通过CMD读取配置文件、按照配置文件的设定值执行。不一致的时候，配置文件优先于CMD参数
	//cobra.OnInitialize(initConfig)
	cobra.OnInitialize(func() {
		vip = viper.New()
		//自动读取环境变量，改写对应值
		vip.AutomaticEnv()
		//设置环境变量的前缀，将 PORT变为 COMI_PORT
		vip.SetEnvPrefix("comi")
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		homeConfigPath := path.Join(home, ".config/comigo")
		vip.AddConfigPath(homeConfigPath)
		// 当前执行目录
		nowPath, _ := os.Getwd()
		vip.AddConfigPath(nowPath)
		vip.SetConfigType("toml")
		if common.ConfigFile == "" {
			vip.SetConfigName("config.toml")
		} else {
			vip.SetConfigName(common.ConfigFile)
		}
		// 設定ファイルを読み込む
		if err := vip.ReadInConfig(); err != nil {
			if common.ConfigFile == "" && common.Config.Debug {
				fmt.Println(err)
			}
		}
		// 把设定文件的内容，解析到构造体里面。
		if err := vip.Unmarshal(&common.Config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//var tomlExample = []byte(`DebugMode = true
		//OpenBrowser = true
		//Host = "localhost"
		//Port = "1234"
		//Debug = false
		//`)
		//vip.ReadConfig(bytes.NewBuffer(tomlExample))
		//保存配置並退出
		if common.Config.NewConfig {
			common.Config.NewConfig = false
			//if common.Config.NewConfig {
			//	vip.SafeWriteConfigAs("config.toml")
			//}
			bytes, err := toml.Marshal(common.Config)
			if err != nil {
				fmt.Println("toml.Marshal Error")
			}
			//在命令行打印
			fmt.Println(string(bytes))
			err = ioutil.WriteFile("config.toml", bytes, 0644)
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	})

}
