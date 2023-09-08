package common

import (
	"github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/settings"
	"net/http"
)

var (
	ConfigFilePath = "" //yaml设置文件路径，数据库文件(comigo.db)在同一个文件夹。
	Version        = "v0.9.5"
	GenerateConfig = false
	ReadingBook    *book.Book
	Srv            *http.Server
	Config         = settings.ServerConfig{
		Port:                  1234,
		Host:                  "DefaultHost",
		StoresPath:            []string{},
		SupportFileType:       []string{".zip", ".tar", ".rar", ".cbr", ".cbz", ".epub", ".tar.gz", ".tgz", ".tar.bz2", ".tbz2", ".tar.xz", ".txz", ".tar.lz4", ".tlz4", ".tar.sz", ".tsz", ".bz2", ".gz", ".lz4", ".sz", ".xz", ".mp4", ".webm", ".pdf", ".m4v", ".flv", ".avi", ".mp3", ".wav", ".wma", ".ogg"},
		SupportMediaType:      []string{".jpg", ".jpeg", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".gif", ".apng", ".bmp", ".webp", ".ico", ".heic", ".heif", ".avif"},
		ExcludePath:           []string{".comigo", ".idea", ".vscode", ".git", "node_modules", "flutter_ui", "$RECYCLE.BIN", "System Volume Information", ".cache"},
		MaxScanDepth:          4,
		MinImageNum:           3,
		ZipFileTextEncoding:   "",
		OpenBrowser:           true,
		UseCache:              true,
		CachePath:             "",
		ClearCacheExit:        true,
		UploadPath:            "",
		EnableUpload:          true,
		EnableDatabase:        false,
		ClearDatabaseWhenExit: true,
		EnableTLS:             false,
		Username:              "comigo",
		Password:              "",
		DisableLAN:            false,
		DefaultMode:           "scroll",
		LogToFile:             false,
		//EnableFrpcServer:      false,
		//FrpConfig: settings.FrpClientConfig{
		//	FrpcCommand:      "frpc",
		//	ServerAddr:       "localhost", //server_addr
		//	ServerPort:       7000,        //server_port
		//	Token:            "&&%%!2356",
		//	FrpType:          "tcp",
		//	RemotePort:       50000, //remote_port
		//	RandomRemotePort: true,
		//},
		//WebpConfig: WebPServerConfig{
		//	WebpCommand:  "webp-server",
		//	HOST:         "127.0.0.1",
		//	PORT:         "3333",
		//	ImgPath:      "",
		//	QUALITY:      70,
		//	AllowedTypes: []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp"},
		//	ExhaustPath:  "",
		//},
		//DatabaseFilePath:     "",
	}
)

var DatabaseBookList []*book.Book
