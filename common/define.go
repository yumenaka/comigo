package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-homedir"
)

type ServerConfig struct {
	OpenBrowser         bool
	OnlyLocal           bool
	PrintAllIP          bool
	Port                int
	ConfigPath          string
	UseGO               bool
	LogToFile           bool
	UseWebpServer       bool
	LogFilePath         string
	LogFileName         string
	MaxDepth            int
	MinImageNum         int
	ZipFilenameEncoding string
	WebpCommand         string
	WebpConfig          WebPServerConfig
	ServerHost          string
}

var Config = ServerConfig{
	OpenBrowser:         true,
	OnlyLocal:           false,
	Port:                1234,
	UseGO:               true,
	LogToFile:           false,
	MaxDepth:            2,
	MinImageNum:         3,
	ZipFilenameEncoding: "",
	WebpCommand:         "webp-server",
	WebpConfig: WebPServerConfig{
		HOST:         "127.0.0.1",
		PORT:         "3333",
		ImgPath:      "",
		QUALITY:      "70",
		AllowedTypes: []string{"jpg", "png", "jpeg", "bmp"},
		ExhaustPath:  "",
	},
	ServerHost: "",
}

func init() {
	var JAVAHOME string
	JAVAHOME = os.Getenv("JAVA_HOME")
	fmt.Println(JAVAHOME)
}

type WebPServerConfig struct {
	HOST         string
	PORT         string
	ImgPath      string `json:"IMG_PATH"`
	QUALITY      string
	AllowedTypes []string `json:"ALLOWED_TYPES"`
	ExhaustPath  string   `json:"EXHAUST_PATH"`
}

func StartWebPServer(configPath string, imgPath string, exhaustPath string, port int) error {
	//Config.WebpCommand = wepBinaryPath
	Config.WebpConfig.ImgPath = imgPath
	Config.WebpConfig.ExhaustPath = exhaustPath
	Config.WebpConfig.PORT = strconv.Itoa(port)
	//Config.WebpConfig.QUALITY = quality
	if Config.WebpCommand == "" || Config.WebpConfig.ImgPath == "" || Config.WebpConfig.ExhaustPath == "" {
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
	err = webpCMD(configPath, Config.WebpCommand)
	return err
}

func webpCMD(configPath string, wepCommand string) (err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(wepCommand, "--config", configPath+"\\config.json")
		fmt.Println(cmd)
		if err = cmd.Start(); err != nil {
			return err
		}
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command(wepCommand, "--config", configPath+"/config.json")
		fmt.Println(cmd)
		if err = cmd.Start(); err != nil {
			return err
		}
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command(wepCommand, "--config", configPath+"/config.json")
		fmt.Println(cmd)
		if err = cmd.Start(); err != nil {
			return err
		}
	}
	return err
}

var ReadingBook Book
var BookList []Book
var (
	//ReadFileName           string
	TempDir         string
	PictureDir      string
	PrintVersion    bool
	Version         string = "v0.1.6"
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
	for i := 0; i < len(b.PageInfo); i++ { //此处不能用range，因为需要修改
		SetImageType(&b.PageInfo[i])
	}
}

//并发分析
func (b *Book) ScanAllImageGo() {
	var wg sync.WaitGroup
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
}

func SetImageType(p *ImageInfo) {
	err := p.GetImageSize()
	fmt.Println("分析图片分辨率中：", p.LocalPath)
	if err != nil {
		fmt.Println("读取分辨率出错：" + err.Error())
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
		fmt.Println("failed to open image: %v", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
	}
	return err
}
