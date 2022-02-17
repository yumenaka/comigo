package common

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-homedir"
	"github.com/xxjwxc/gowp/workpool"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
	"image"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

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
	////本地管理界面，现在用不着
	//AdminAddr   string
	//AdminPort   string
	//AdminUser   string
	//AdminPwd    string
	//本地转发端口设置
	FrpType          string
	RemotePort       int
	RandomRemotePort bool
}

var ConfigFile string //服务器配置文件路径

type ServerConfig struct {
	UserName    string `json:"-"` //不要解析这个字段
	Password    string `json:"-"` //不要解析这个字段
	CertFile    string `json:"-"` //不要解析这个字段
	KeyFile     string `json:"-"` //不要解析这个字段
	OpenBrowser bool   `json:"-"` //不要解析这个字段
	DisableLAN  bool   `json:"-"` //不要解析这个字段
	Template    string `json:"template"`
	//Auth                   string `json:"-"` //不要解析这个字段 访问密码，还没做
	PrintAllIP             bool `json:"-"` //不要解析这个字段
	Port                   int
	CheckImage             bool
	Debug                  bool   `json:"-"` //不要解析这个字段
	LogToFile              bool   `json:"-"` //不要解析这个字段
	LogFilePath            string `json:"-"` //不要解析这个字段
	LogFileName            string `json:"-"` //不要解析这个字段
	MaxDepth               int    `json:"-"` //不要解析这个字段
	MinImageNum            int
	Host                   string
	EnableWebpServer       bool
	WebpConfig             WebPServerConfig `json:"-"` //不要解析这个字段
	EnableFrpcServer       bool
	FrpConfig              FrpClientConfig `json:"-"` //不要解析这个字段
	ZipFilenameEncoding    string          `json:"-"` //不要解析这个字段
	SketchCountSeconds     int             `json:"sketch_count_seconds"`
	SortImage              string
	TempPATH               string `json:"-"` //不要解析这个字段
	CleanAllTempFileOnExit bool   `json:"-"` //不要解析这个字段
	CleanAllTempFile       bool   `json:"-"` //不要解析这个字段
	NewConfig              bool   `json:"-"` //不要解析这个字段
}

var Config = ServerConfig{
	OpenBrowser:         true,
	DisableLAN:          false,
	Template:            "scrool", //multi、single、random etc.
	Port:                1234,
	CheckImage:          false,
	LogToFile:           false,
	MaxDepth:            3,
	MinImageNum:         3,
	ZipFilenameEncoding: "",
	WebpConfig: WebPServerConfig{
		WebpCommand:  "webp-server",
		HOST:         "127.0.0.1",
		PORT:         "3333",
		ImgPath:      "",
		QUALITY:      70,
		AllowedTypes: []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp"},
		ExhaustPath:  "",
	},
	EnableFrpcServer: false,
	FrpConfig: FrpClientConfig{
		FrpcCommand:      "frpc",
		ServerAddr:       "localhost", //server_addr
		ServerPort:       7000,        //server_port
		Token:            "&&%%!2356",
		FrpType:          "tcp",
		RemotePort:       50000, //remote_port
		RandomRemotePort: true,
		//AdminAddr:   "127.0.0.1",
		//AdminPort:   "12340",
		//AdminUser:   "",
		//AdminPwd :   "",
	},
	Host:                   "",
	SketchCountSeconds:     90,
	SortImage:              "",
	TempPATH:               "",
	CleanAllTempFileOnExit: true,
	CleanAllTempFile:       true,
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (config *ServerConfig) SetByExecutableFilename() {
	// 当前执行目录
	//targetPath, _ := os.Getwd()
	//fmt.Println(locale.GetString("target_path"), targetPath)
	// 带后缀的执行文件名 comi.exe  sketch.exe
	filenameWithSuffix := path.Base(os.Args[0])
	// 执行文件名后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	// 去掉后缀后的执行文件名
	filenameWithOutSuffix := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//fmt.Println("filenameWithOutSuffix =", filenameWithOutSuffix)
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	extPath := filepath.Dir(ex)
	//fmt.Println("extPath =",extPath)
	ExtFileName := strings.TrimPrefix(filenameWithOutSuffix, extPath)
	//fmt.Println("ExtFileName =", ExtFileName)
	//如果执行文件名包含 scroll 等关键字，选择卷轴模板
	if haveKeyWord(ExtFileName, []string{"scroll", "スクロール", "默认", "下拉", "卷轴"}) {
		config.Template = "scroll"
	}
	//如果执行文件名包含 sketch 等关键字，选择速写模板
	if haveKeyWord(ExtFileName, []string{"sketch", "croquis", "クロッキー", "素描", "速写"}) {
		config.Template = "sketch"
	}
	//根据文件名设定倒计时秒数,不管默认是不是sketch模式
	Seconds, err := getNumberFromString(ExtFileName)
	if err != nil {
		if config.Template == "sketch" {
			//fmt.Println(Seconds)
		}
	} else {
		config.SketchCountSeconds = Seconds
	}
	//如果执行文件名包含 single 等关键字，选择 flip分页漫画模板
	if haveKeyWord(ExtFileName, []string{"flip", "翻页", "めく"}) {
		config.Template = "flip"
	}

	//选择模式以后，打印提示
	switch config.Template {
	case "scroll":
		fmt.Println(locale.GetString("scroll_template"))
	case "flip":
		fmt.Println(locale.GetString("single_page_template"))
	case "sketch":
		fmt.Println(locale.GetString("sketch_template"))
		//速写倒计时秒数
		fmt.Println(locale.GetString("SKETCH_COUNT_SECONDS"), config.SketchCountSeconds)
	default:
	}
}

//从字符串中提取数字,如果有几个数字，就简单地加起来
func getNumberFromString(s string) (int, error) {
	var err error
	num := 0
	//同时设定倒计时秒数
	valid := regexp.MustCompile("[0-9]+")
	numbers := valid.FindAllStringSubmatch(s, -1)
	if len(numbers) > 0 {
		//循环取出多维数组
		for _, value := range numbers {
			for _, v := range value {
				temp, errTemp := strconv.Atoi(v)
				if errTemp != nil {
					fmt.Println("error num value:" + v)
				} else {
					num = num + temp
				}
			}
		}
		//fmt.Println("get Number:",num," form string:",s,"numbers[]=",numbers)
	} else {
		err = errors.New("number not found")
		return 0, err
	}
	return num, err
}

//检测字符串中是否有关键字
func haveKeyWord(checkString string, list []string) bool {
	//转换为小写，使Sketch、DOUBLE也生效
	checkString = strings.ToLower(checkString)
	for _, key := range list {
		if strings.Contains(checkString, key) {
			return true
		}
	}
	return false
}

var ReadingBook Book
var BookList []Book
var (
	CacheFilePath    string
	Version          = "v0.5.1"
	ExcludeFolders   = []string{".comigo", "node_modules", "flutter_ui", "$RECYCLE.BIN", "Config.Msi"}
	SupportMediaType = []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp", ".webp", ".ico", ".heic", ".pdf", ".mp4", ".webm"}
	SupportFileType  = [...]string{
		".zip",
		".tar",
		".rar",
		".cbr",
		".cbz",
		".epub",
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
	Config.LogFileName = "comigo.log"
}

type Book struct {
	Name            string      `json:"name"`
	Author          string      `json:"author"`
	Title           string      `json:"title"`
	FilePath        string      `json:"-"` //不要解析这个字段
	ExtractPath     string      `json:"-"` //不要解析这个字段
	AllPageNum      int         `json:"all_page_num"`
	FileType        string      `json:"file_type"`
	FileSize        int64       `json:"file_size"`
	Modified        time.Time   `json:"modified_time"`
	BookID          string      `json:"uuid"` //根据FilePath计算
	IsDir           bool        `json:"is_folder"`
	ExtractNum      int         `json:"extract_num"`
	ExtractComplete bool        `json:"extract_complete"`
	ReadPercent     float64     `json:"read_percent"`
	PageInfo        AllPageInfo `json:"pages"`
}

type SinglePageInfo struct {
	ModeTime          time.Time `json:"-"` //不要解析这个字段
	FileSize          int64     `json:"-"` //不要解析这个字段
	Height            int       `json:"height"`
	Width             int       `json:"width"`
	Url               string    `json:"url"`      //远程用户读取图片的实际URL，为了适应特殊字符，经过一次转义
	InArchiveName     string    `json:"filename"` //书籍为压缩文件的时候，用于解压的压缩文件内文件路径
	RealImageFilePATH string    `json:"-"`        //不要解析这个字段  书籍为文件夹的时候，实际图片的路径
	ImgType           string    `json:"image_type"`
	Blurhash          string    `json:"blurhash"`
}

// AllPageInfo Slice
type AllPageInfo []SinglePageInfo

func (s AllPageInfo) Len() int {
	return len(s)
}

// Less 按时间或URL，将图片排序
func (s AllPageInfo) Less(i, j int) (less bool) {
	//如何定义 s[i] < s[j]  根据文件名
	numI, err1 := getNumberFromString(s[i].InArchiveName)
	if err1 != nil {
		less = strings.Compare(s[i].InArchiveName, s[j].InArchiveName) > 0
		return less
	}
	numJ, err2 := getNumberFromString(s[j].InArchiveName)
	if err2 != nil {
		less = strings.Compare(s[i].InArchiveName, s[j].InArchiveName) > 0
		return less
	}
	//fmt.Println("numI:",numI)
	//fmt.Println("numJ:",numJ)
	less = numI < numJ //如果有的话，比较文件名里的数字

	//如何定义 s[i] < s[j]  根据修改时间
	if Config.SortImage == "time" {
		less = s[i].ModeTime.After(s[j].ModeTime) // s[i] 的年龄（修改时间），是否比 s[j] 小？
	}
	return less
}

func (s AllPageInfo) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// SortPages 上面三个函数定义好了，终于可以使用sort包排序了
func (b *Book) SortPages() {
	sort.Sort(b.PageInfo)
}

// SetBookID  根据路径的MD5，生成书籍ID
func (b *Book) SetBookID() {
	//fmt.Println("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	fileAbaPath, err := filepath.Abs(b.FilePath)
	if err != nil {
		fmt.Println(err, fileAbaPath)
	}
	//fmt.Println(md5s(fileAbaPath))
	b.BookID = md5string(fileAbaPath)
}

// GetBookID  根据路径的MD5，生成书籍ID
func (b *Book) GetBookID() string {
	if b.BookID == "" {
		b.SetBookID()
	}
	return b.BookID
}

// ScanArchive  设置书名与Book ID等
func (b *Book) InitBook(name string) {
	//文件夹直接用路径
	if b.IsDir {
		b.Name = name
	} else {
		//压缩文件的话，去除路径，取文件名
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
		b.FileType = path.Ext(name) //获取文件后缀
	}
	b.SetBookID()
}

func (b *Book) SetPageNum() {
	//页数，目前只支持漫画
	b.AllPageNum = len(b.PageInfo)
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
		if isSupportMedia(p.Url) {
			PicNum++
		}
	}
	return PicNum
}

// ScanAllImage 服务器端分析单双页
func (b *Book) ScanAllImage() {
	log.Println(locale.GetString("check_image_start"))
	// Console progress bar
	bar := pb.StartNew(b.AllPageNum)
	tmpl := `{{ red "With funcs:" }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}} {{string . "my_green_string" | green}} {{string . "my_blue_string" | blue}}`
	bar.SetTemplateString(tmpl)
	for i := 0; i < len(b.PageInfo); i++ { //此处不能用range，因为需要修改
		SetImageType(&b.PageInfo[i])
		//进度条计数
		bar.Increment()
	}
	// 进度条跑完
	bar.Finish()
	log.Println(locale.GetString("check_image_completed"))
}

// ScanAllImageGo 并发分析
func (b *Book) ScanAllImageGo() {
	//var wg sync.WaitGroup
	log.Println(locale.GetString("check_image_start"))
	wp := workpool.New(10) //设置最大线程数
	//res := make(chan string)
	count := 0
	// Console progress bar
	bar := pb.StartNew(b.AllPageNum)
	for i := 0; i < len(b.PageInfo); i++ { //此处不能用range，因为需要修改
		//wg.Add(1)
		count++
		ii := i
		//并发处理，提升图片分析速度
		wp.Do(func() error {
			//defer wg.Done()
			SetImageType(&b.PageInfo[ii])
			bar.Increment()
			//res <- fmt.Sprintf("Finished %d", i)
			return nil
		})
	}
	//wg.Wait()
	_ = wp.Wait()
	// finish bar
	bar.Finish()
	log.Println(locale.GetString("check_image_completed"))
}

func SetImageType(p *SinglePageInfo) {
	err := p.GetImageSize()
	//log.Println(locale.GetString("check_image_ing"), p.RealImageFilePATH)
	if err != nil {
		log.Println(locale.GetString("check_image_error") + err.Error())
	}
	if p.Width == 0 && p.Height == 0 {
		p.ImgType = "UnKnow"
		return
	}
	if p.Width > p.Height {
		p.ImgType = "DoublePage"
	} else {
		p.ImgType = "SinglePage"
	}
}

// GetImageSize 获取图片分辨率
func (i *SinglePageInfo) GetImageSize() (err error) {
	var img image.Image
	img, err = imaging.Open(i.RealImageFilePATH)
	if err != nil {
		log.Printf(locale.GetString("check_image_error")+" %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
		//简单占位符，很耗费服务器资源，以后再研究。
		////"github.com/buckket/go-blurhash"
		//str, err := blurhash.Encode(4, 3, img)
		//if err != nil {
		//	// Handle errors
		//	log.Printf(locale.GetString("check_image_error")+" %v\n", err)
		//}
		//i.Blurhash = str
	}
	return err
}

// SetupCloseHander 中断处理：程序被中断的时候，清理临时文件
func SetupCloseHander() {
	c := make(chan os.Signal, 2)
	//SIGHUP（挂起）, SIGINT（中断）或 SIGTERM（终止）默认会使得程序退出。
	//1、SIGHUP 信号在用户终端连接(正常或非正常)结束时发出。
	//2、syscall.SIGINT 和 os.Interrupt 是同义词,按下 CTRL+C 时发出。
	//3、SIGTERM（终止）:kill终止进程,允许程序处理问题后退出。
	//4.syscall.SIGHUP,终端控制进程结束(终端连接断开)
	//5、syscall.SIGQUIT，CTRL+\ 退出
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-c
		if Config.CleanAllTempFileOnExit {
			fmt.Println("\r" + locale.GetString("start_clear_file"))
			clearTempFilesALL()
		} else {
			clearTempFilesOne(&ReadingBook)
		}
		os.Exit(0)
	}()
}

//func InitReadingBook() (err error) {
//	//准备解压，设置图片文件夹
//	if ReadingBook.IsDir {
//		ReadingBook.ExtractPath = ReadingBook.FilePath
//		ReadingBook.ExtractComplete = true
//		ReadingBook.ExtractNum = ReadingBook.AllPageNum
//	} else {
//		//setTempDir()
//		ReadingBook.ExtractPath = path.Join(CacheFilePath, ReadingBook.GetBookID()) //extraFolder
//		//err = LsArchive(&ReadingBook)
//		//if err != nil {
//		//	fmt.Println(locale.GetString("scan_archive_error"))
//		//	return err
//		//}
//		//err = UnArchive(&ReadingBook)
//		//if err != nil {
//		//	fmt.Println(locale.GetString("un_archive_error"))
//		//	return err
//		//}
//		ReadingBook.InitBook(ReadingBook.FilePath) //设置书名
//	}
//	//服务器分析图片，新版默认不做
//	if Config.CheckImage {
//		ReadingBook.ScanAllImageGo() //扫描所有图片，取得分辨率信息，使用了协程
//	}
//	//服务器排序图片
//	if Config.SortImage != "" {
//		if Config.SortImage == "name" {
//			ReadingBook.SortPages()
//			fmt.Println(locale.GetString("SORT_BY_NAME"))
//		}
//		if Config.SortImage == "time" {
//			ReadingBook.SortPages()
//			fmt.Println(locale.GetString("SORT_BY_TIME"))
//		}
//		if Config.Debug {
//			//判断是否已经排好顺序，将会打印true
//			fmt.Println("IS Sorted?\t", sort.IsSorted(ReadingBook.PageInfo))
//			//打印排序后的数据
//			//litter.Dump(ReadingBook.PageInfo)
//		}
//	}
//	return err
//}

// setTempDir 设置临时文件夹，退出时会被清理
func setTempDir() {
	//手动设置的临时文件夹
	if Config.TempPATH != "" && tools.ChickExists(Config.TempPATH) && tools.ChickIsDir(Config.TempPATH) {
		CacheFilePath = path.Join(Config.TempPATH)
	} else {
		CacheFilePath = path.Join(os.TempDir(), "comigo_temp_files") //直接使用系统文件夹
	}
	err := os.MkdirAll(CacheFilePath, os.ModePerm)
	if err != nil {
		println(locale.GetString("temp_folder_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + CacheFilePath)
	}
}

// 清空解压缓存
func clearTempFilesALL() {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	for _, tempBook := range BookList {
		clearTempFilesOne(&tempBook)
	}
}

// 清空某一本压缩漫画的解压缓存
func clearTempFilesOne(book *Book) {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	haveThisBook := false
	for _, tempBook := range BookList {
		if tempBook.GetBookID() == book.GetBookID() {
			haveThisBook = true
		}
	}
	if haveThisBook {
		extractPath := path.Join(CacheFilePath, book.GetBookID())
		//避免删错文件,解压路径包含UUID，len不可能小于32
		PathLen := len(extractPath)
		if PathLen < 32 {
			return
		}
		err := os.RemoveAll(extractPath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + extractPath)
		} else {
			fmt.Println(locale.GetString("clear_temp_file_completed") + extractPath)
		}
	}
}
