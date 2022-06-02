package common

import (
	"fmt"
	"net/http"
	"os"
	"path"

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
	ConfigFile  = "" //yaml设置文件路径，数据库文件(comigo.db)在同一个文件夹。
	Version     = "v0.7.3"
	ReadingBook *book.Book
	Srv         *http.Server
	Config      = settings.ServerSettings{
		Port:                 1234,
		Host:                 "",
		StoresPath:           []string{},
		CacheFileEnable:      true,
		CacheFilePath:        "",
		CacheFileClean:       true,
		EnableDatabase:       false,
		DatabaseFilePath:     "",
		OpenBrowser:          true,
		DisableLAN:           false,
		DefaultMode:          "scroll",
		GenerateMetaData:     false,
		LogToFile:            false,
		MaxDepth:             3,
		MinImageNum:          3,
		ZipFileTextEncoding:  "",
		EnableFrpcServer:     false,
		SupportFileType:      []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".mp4", ".webm"},
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
	}
)

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
