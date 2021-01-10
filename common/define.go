package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

type ServerConfig struct {
	OpenBrowser         bool	`json:"-"` //不要解析这个字段
	DisableLAN          bool    `json:"-"` //不要解析这个字段
	DefaultPageMode     string  `json:"default_page_mode"`
	PrintAllIP          bool    `json:"-"` //不要解析这个字段
	Port                int
	ConfigPath          string  `json:"-"` //不要解析这个字段
	CheckImageInServer  bool
	LogToFile           bool    `json:"-"` //不要解析这个字段
	LogFilePath         string  `json:"-"` //不要解析这个字段
	LogFileName         string  `json:"-"` //不要解析这个字段
	MaxDepth            int     `json:"-"` //不要解析这个字段
	MinImageNum         int
	ServerHost          string
	EnableWebpServer    bool
	WebpConfig          WebPServerConfig    `json:"-"` //不要解析这个字段
	EnableFrpcServer    bool
	FrpConfig           FrpClientConfig     `json:"-"` //不要解析这个字段
	ZipFilenameEncoding string              `json:"-"` //不要解析这个字段
}

var Config = ServerConfig{
	OpenBrowser:         true,
	DisableLAN:          false,
	DefaultPageMode:     "random",//multi、single、random etc.
	Port:                1234,
	CheckImageInServer:  false,
	LogToFile:           false,
	MaxDepth:            2,
	MinImageNum:         3,
	ZipFilenameEncoding: "",
	WebpConfig: WebPServerConfig{
		WebpCommand:  "webp-server",
		HOST:         "127.0.0.1",
		PORT:         "3333",
		ImgPath:      "",
		QUALITY:      70,
		AllowedTypes: []string{"jpg", "png", "jpeg", "bmp"},
		ExhaustPath:  "",
	},
	EnableFrpcServer: false,
	FrpConfig: FrpClientConfig{
		FrpcCommand:      "frpc",
		ServerAddr:       "localhost", //server_addr
		ServerPort:       7000,        //server_port
		Token:            "&&%%!2356",
		FrpType:          "tcp",
		RemotePort:       -1, //remote_port
		RandomRemotePort: true,
		//AdminAddr:   "127.0.0.1",
		//AdminPort:   "12340",
		//AdminUser:   "",
		//AdminPwd :   "",
	},
	ServerHost: "",
}

type WebPServerConfig struct {
	WebpCommand  string
	HOST         string
	PORT         string
	ImgPath      string `json:"IMG_PATH"`
	QUALITY      int
	AllowedTypes []string `json:"ALLOWED_TYPES"`
	ExhaustPath  string   `json:"EXHAUST_PATH"`
}

type FrpClientConfig struct {
	//frp，服务器端
	FrpcCommand string
	ServerAddr  string
	ServerPort  int
	Token       string
	////本地管理界面，暂不开启
	//AdminAddr   string
	//AdminPort   string
	//AdminUser   string
	//AdminPwd    string
	//本地转发端口设置
	FrpType          string
	RemotePort       int
	RandomRemotePort bool
}

func StartFrpC(configPath string) error {
	//借助ini库，保存一个ini文件
	cfg := ini.Empty()
	//配置文件类似：
	//[common]
	//server_addr = frp.example.net
	//server_port = 7000
	//token = Nscffaass
	//[comi]
	//type = tcp
	//local_ip = 127.0.0.1
	//local_port = 1234
	//remote_port = 23456
	_, err := cfg.NewSection("common")
	_, err = cfg.Section("common").NewKey("server_addr", Config.FrpConfig.ServerAddr)
	_, err = cfg.Section("common").NewKey("server_port", strconv.Itoa(Config.FrpConfig.ServerPort))
	_, err = cfg.Section("common").NewKey("token", Config.FrpConfig.Token)
	FrpConfigName := ReadingBook.Name + "(" + "comi " + Version + " " + time.Now().Format("2006-01-02 15:04:05") + ")"
	_, err = cfg.NewSection(FrpConfigName)
	_, err = cfg.Section(FrpConfigName).NewKey("type", Config.FrpConfig.FrpType)
	_, err = cfg.Section(FrpConfigName).NewKey("local_ip", "127.0.0.1")
	_, err = cfg.Section(FrpConfigName).NewKey("local_port", strconv.Itoa(Config.Port))
	_, err = cfg.Section(FrpConfigName).NewKey("remote_port", strconv.Itoa(Config.FrpConfig.RemotePort))
	//保存文件
	err = cfg.SaveToIndent(configPath+"/frpc.ini", "\t")
	if err != nil {
		fmt.Println("frpc ini初始化错误")
		return err
	} else {
		fmt.Println("成功保存frpc设定.", configPath, cfg)
	}
	//实际执行
	var cmd *exec.Cmd
	cmd = exec.Command(Config.FrpConfig.FrpcCommand, "-c", configPath+"/frpc.ini")
	fmt.Println(cmd)
	if err = cmd.Start(); err != nil {
		return err
	}
	return err
}

func StartWebPServer(configPath string, imgPath string, exhaustPath string, port int) error {
	//Config.WebpCommand = wepBinaryPath
	Config.WebpConfig.ImgPath = imgPath
	Config.WebpConfig.ExhaustPath = exhaustPath
	Config.WebpConfig.PORT = strconv.Itoa(port)
	//Config.WebpConfig.QUALITY = quality
	if Config.WebpConfig.WebpCommand == "" || Config.WebpConfig.ImgPath == "" || Config.WebpConfig.ExhaustPath == "" {
		return errors.New("webp设定错误")
	}
	jsonObject, err := os.OpenFile(configPath+"/config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer jsonObject.Close()
	content, err := json.Marshal(Config.WebpConfig)
	if err != nil {
		return err
	}
	if _, err := jsonObject.Write(content); err == nil {
		fmt.Println("成功保存webp设定.", configPath, content)
	}
	//err = webpCMD(configPath, Config.WebpCommand)
	var cmd *exec.Cmd
	cmd = exec.Command(Config.WebpConfig.WebpCommand, "--config", configPath+"/config.json")
	fmt.Println(cmd)
	if err = cmd.Start(); err != nil {
		return err
	}
	return err
}

var ReadingBook Book
var BookList []Book
var (
	//ReadFileName           string
	TempDir    string
	PictureDir string
	//PrintVersion    bool
	Version         string = "v0.2.4"
	SupportPicType         = [...]string{".png", ".jpg", ".jpeg", "bmp", ".gif", ".webp"}
	SupportFileType        = [...]string{
		".zip",
		".tar",
		".rar",
		".tar.gz",
		".tgz",
		".tar.bz2",
		".tbz2",
		".tar.xz",
		".txz",
		".tar.lz4",
		".tlz4",
		".tar.sz",
		".tsz",
		".bz2",
		".gz",
		".lz4",
		".sz",
		".xz"}
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	Config.LogFilePath = home
	Config.LogFileName = "comigo.ini"
}

type Book struct {
	Name            string      `json:"name"`
	Author          string      `json:"author"`
	Title           string      `json:"title"`
	FilePath        string      `json:"-"` //不要解析这个字段
	PageNum         int         `json:"page_num"`
	PageInfo        []ImageInfo `json:"pages"`
	FileType        string      `json:"file_type"`
	FileSize        int64       `json:"file_size"`
	Modified        time.Time   `json:"modified_time"`
	UUID            string      `json:"uuid"`
	IsFolder        bool        `json:"is_folder"`
	ExtractNum      int         `json:"extract_num"`
	ExtractComplete bool        `json:"extract_complete"`
	ReadPercent     float64     `json:"read_percent"`
}

type ImageInfo struct {
	Height        int    `json:"height"`
	Width         int    `json:"width"`
	UrlPath       string `json:"url"`
	LocalPath     string `json:"-"` //不要解析这个字段
	InArchiveName string `json:"-"` //不要解析这个字段
	ImgType       string `json:"class"`
}

//一些绑定到Book结构体的方法
func (b *Book) SetArchiveBookName(name string) {
	post := strings.LastIndex(name, "/") //Unix路径分隔符
	if post == -1 {
		post = strings.LastIndex(name, "\\") //windows分隔符
	}
	if post != -1 {
		//name = string([]rune(name)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来?
		name = name[post:]
		name = strings.ReplaceAll(name, "\\", "")
		name = strings.ReplaceAll(name, "/", "")
	}
	b.Name = name
}

func (b *Book) SetImageFolderBookName(name string) {
	b.Name = name
}

func (b *Book) SetPageNum() {
	//页数，目前只支持漫画
	b.PageNum = len(b.PageInfo)
}

func (b *Book) SetFilePath(path string) {
	b.FilePath = path
}

func (b *Book) GetName() string { //绑定到Book结构体的方法
	return b.Name
}

func (b *Book) GetPicNum() int {
	var PicNum = 0
	for _, p := range b.PageInfo {
		if checkPicExt(p.UrlPath) {
			PicNum++
		}
	}
	return PicNum
}

//服务器端分析单双页
func (b *Book) ScanAllImage() {
	log.Println("开始解析图片")
	for i := 0; i < len(b.PageInfo); i++ { //此处不能用range，因为需要修改
		SetImageType(&b.PageInfo[i])
	}
	log.Println("图片解析完成")
}

//并发分析
func (b *Book) ScanAllImageGo() {
	var wg sync.WaitGroup
	log.Println("开始分析图片分辨率")
	for i := 0; i < len(b.PageInfo); i++ { //此处不能用range，因为需要修改
		wg.Add(1)
		//并发处理，提升图片分析速度
		go func(i int) {
			defer wg.Done()
			SetImageType(&b.PageInfo[i])
		}(i)
		//if i < 10 {//为了优化打开速度，即便并发分析，前10张也要单线程做
		//	SetImageType(&b.PageInfo[i])
		//} else {
		//	wg.Add(1)
		//	//并发处理，提升图片分析速度
		//	go func(i int) {
		//		defer wg.Done()
		//		SetImageType(&b.PageInfo[i])
		//	}(i)
		//}
	}
	wg.Wait()
	log.Println("图片分辨率分析完成")
}

func SetImageType(p *ImageInfo) {
	err := p.GetImageSize()
	//log.Println("分析图片分辨率中：", p.LocalPath)
	if err != nil {
		log.Println("读取分辨率出错：" + err.Error())
	}
	if p.Width == 0 && p.Height == 0 {
		p.ImgType = "UnKnow"
		return
	}
	if p.Width > p.Height {
		p.ImgType = "Horizontal"
	} else {
		p.ImgType = "Vertical"
	}
}

//获取图片分辨率
func (i *ImageInfo) GetImageSize() (err error) {
	var img image.Image
	img, err = imaging.Open(i.LocalPath)
	if err != nil {
		log.Printf("failed to open image: %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
	}
	return err
}
