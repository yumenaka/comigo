package common

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir" //不使用 cgo 获取用户主目录的第三方库，支持交叉编译

	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/settings"
	"github.com/yumenaka/comi/tools"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	Config.LogFilePath = home
	Config.LogFileName = "comigo.log"
}

var (
	ConfigFile  = ""
	Version     = "v0.7.3"
	ReadingBook *book.Book
	Srv         *http.Server

	Config = settings.ServerSettings{
		Port:                 1234,
		Host:                 "",
		StoresPath:           []string{},
		CacheFileEnable:      true,
		CacheFilePath:        "",
		CacheFileClean:       true,
		DatabaseFilePath:     "",
		OpenBrowser:          true,
		DisableLAN:           false,
		GenerateMetaData:     false,
		LogToFile:            false,
		MaxDepth:             3,
		MinImageNum:          3,
		ZipFileTextEncoding:  "",
		EnableFrpcServer:     false,
		SupportFileType:      []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".pdf", ".mp4", ".webm"},
		SupportMediaType:     []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
		ExcludeFileOrFolders: []string{".comigo", ".idea", ".vscode", ".git", "node_modules", "flutter_ui", "$RECYCLE.BIN", "System Volume Information", ".cache"},
		FrpConfig: settings.FrpClientConfig{
			FrpcCommand:      "frpc",
			ServerAddr:       "localhost", //server_addr
			ServerPort:       7000,        //server_port
			Token:            "&&%%!2356",
			FrpType:          "tcp",
			RemotePort:       50000, //remote_port
			RandomRemotePort: true,
		},
		//WebpConfig: WebPServerConfig{
		//	WebpCommand:  "webp-server",
		//	HOST:         "127.0.0.1",
		//	PORT:         "3333",
		//	ImgPath:      "",
		//	QUALITY:      70,
		//	AllowedTypes: []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp"},
		//	ExhaustPath:  "",
		//},
		//Template:             "scroll", //multi、single、random etc.
	}
)

// CheckPathSkip 检查路径是否应该跳过（排除文件，文件夹列表）。
func CheckPathSkip(path string) bool {
	for _, substr := range Config.ExcludeFileOrFolders {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}

//// SetupCloseHander 中断处理：程序被中断的时候，清理临时文件
//// 在一个新的 goroutine 上创建一个监听器。 如果接收到了一个 interrupt 信号，就会立即通知程序，做一些清理工作并退出
//func SetupCloseHander() {
//	//中断处理：程序被中断的时候，清理临时文件
//	//容量2(capacity)代表Channel容纳的最多的元素的数量，代表Channel的缓存的大小。如果设置了缓存，就有可能不发生阻塞， 只有buffer满了后 send才会阻塞， 而只有缓存空了后receive才会阻塞。
//	c := make(chan os.Signal, 2)
//	//SIGHUP（挂起）, SIGINT（中断）或 SIGTERM（终止）默认会使得程序退出。
//	//1、SIGHUP 信号在用户终端连接(正常或非正常)结束时发出。
//	//2、syscall.SIGINT 和 os.Interrupt 是同义词,按下 CTRL+C 时发出。
//	//3、SIGTERM（终止）:kill终止进程,允许程序处理问题后退出。
//	//4.syscall.SIGHUP,终端控制进程结束(终端连接断开)
//	//5、syscall.SIGQUIT，CTRL+\ 退出
//
//	// kill (no param) default send syscall.SIGTERM
//	// kill -2 is syscall.SIGINT
//	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
//	//signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
//	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
//	go func() {
//		var temp = <-c
//		fmt.Println(temp)
//		if Config.CacheFileClean {
//			fmt.Println("\r" + locale.GetString("start_clear_file"))
//			ClearTempFilesALL()
//		}
//		os.Exit(0)
//	}()
//}

// SetTempDir 设置临时文件夹，退出时会被清理
func SetTempDir() {
	//手动设置的临时文件夹
	if Config.CacheFilePath != "" && tools.ChickExists(Config.CacheFilePath) && tools.ChickIsDir(Config.CacheFilePath) {
		Config.CacheFilePath = path.Join(Config.CacheFilePath)
	} else {
		Config.CacheFilePath = path.Join(os.TempDir(), "comigo_cache") //直接使用系统文件夹
	}
	err := os.MkdirAll(Config.CacheFilePath, os.ModePerm)
	if err != nil {
		println(locale.GetString("temp_folder_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + Config.CacheFilePath)
	}
}
