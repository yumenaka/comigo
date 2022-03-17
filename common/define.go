package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bbrks/go-blurhash"
	"github.com/cheggaaa/pb/v3"
	"github.com/disintegration/imaging"
	"github.com/mitchellh/go-homedir"
	"github.com/xxjwxc/gowp/workpool"
	"github.com/yumenaka/comi/arch"
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

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	Config.LogFilePath = home
	Config.LogFileName = "comigo.log"
	//退出时清理
	SetupCloseHander()
}

var (
	ReadingBook          Book
	BookList             []Book
	CacheFilePath        = ""
	ConfigFile           = ""
	Version              = "v0.5.2"
	ExcludeFileOrFolders = []string{".comigo", "node_modules", "flutter_ui", "$RECYCLE.BIN", "Config.Msi"}
	SupportMediaType     = []string{".jpg", ".jpeg", ".JPEG", ".jpe", ".jpf", ".jfif", ".jfi", ".png", ".bmp", ".webp", ".ico", ".heic", ".pdf", ".mp4", ".webm"}
	SupportFileType      = [...]string{
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
	Config = ServerConfig{
		OpenBrowser:         true,
		DisableLAN:          false,
		Template:            "scroll", //multi、single、random etc.
		Port:                1234,
		GenerateMetaData:    false,
		LogToFile:           false,
		MaxDepth:            3,
		MinImageNum:         3,
		ZipFileTextEncoding: "",
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

type ServerConfig struct {
	MinImageNum            int
	Host                   string
	EnableWebpServer       bool
	EnableFrpcServer       bool
	Port                   int
	GenerateMetaData       bool
	SortImage              string
	SketchCountSeconds     int              `json:"sketch_count_seconds"`
	Template               string           `json:"template"`
	UserName               string           `json:"-"` //不要解析这个字段
	Password               string           `json:"-"` //不要解析这个字段
	CertFile               string           `json:"-"` //不要解析这个字段
	KeyFile                string           `json:"-"` //不要解析这个字段
	OpenBrowser            bool             `json:"-"` //不要解析这个字段
	DisableLAN             bool             `json:"-"` //不要解析这个字段
	PrintAllIP             bool             `json:"-"` //不要解析这个字段
	Debug                  bool             `json:"-"` //不要解析这个字段
	LogToFile              bool             `json:"-"` //不要解析这个字段
	LogFilePath            string           `json:"-"` //不要解析这个字段
	LogFileName            string           `json:"-"` //不要解析这个字段
	MaxDepth               int              `json:"-"` //不要解析这个字段
	ZipFileTextEncoding    string           `json:"-"` //不要解析这个字段
	TempPATH               string           `json:"-"` //不要解析这个字段
	CleanAllTempFileOnExit bool             `json:"-"` //不要解析这个字段
	CleanAllTempFile       bool             `json:"-"` //不要解析这个字段
	GenerateConfig         bool             `json:"-"` //不要解析这个字段
	WebpConfig             WebPServerConfig `json:"-"` //不要解析这个字段
	FrpConfig              FrpClientConfig  `json:"-"` //不要解析这个字段
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

type Book struct {
	Name            string    `json:"name"`
	Author          string    `json:"author"`
	Title           string    `json:"title"`
	filePath        string    `json:"-"` //不要解析这个字段
	ExtractPath     string    `json:"-"` //不要解析这个字段
	AllPageNum      int       `json:"all_page_num"`
	FileType        string    `json:"file_type"`
	FileSize        int64     `json:"file_size"`
	Modified        time.Time `json:"modified_time"`
	BookID          string    `json:"uuid"` //根据FilePath计算
	IsDir           bool      `json:"is_folder"`
	ExtractNum      int       `json:"extract_num"`
	ExtractComplete bool      `json:"extract_complete"`
	ReadPercent     float64   `json:"read_percent"`
	//NonUTF8Zip 表示 Name 和 Comment 未以 UTF-8 编码。根据规范，唯一允许的其他编码应该是 CP-437，但从历史上看，许多 ZIP 阅读器将 Name 和 Comment 解释为系统的本地字符编码。仅当用户打算为特定本地化区域编码不可移植的 ZIP 文件时，才应设置此标志。否则，Writer 会自动为有效的 UTF-8 字符串设置 ZIP 格式的 UTF-8 标志。
	NonUTF8Zip      bool           `json:"non_utf8_zip"`
	ZipTextEncoding string         `json:"zip_text_encoding"`
	Cover           SinglePageInfo `json:"cover"`
	Pages           AllPageInfo    `json:"pages"`
}

type SinglePageInfo struct {
	NameInArchive string    `json:"filename"` //用于解压的压缩文件内文件路径，或图片名，为了适应特殊字符，经过一次转义
	Url           string    `json:"url"`      //远程用户读取图片的URL，为了适应特殊字符，经过一次转义
	Blurhash      string    `json:"blurhash"` //blurhash占位符。需要扫描图片生成（tools.GetImageDataBlurHash）
	Height        int       `json:"height"`   //blurhash用，图片的高
	Width         int       `json:"width"`    //blurhash用，图片的宽
	ModeTime      time.Time `json:"-"`        //不要解析这个字段
	FileSize      int64     `json:"-"`        //不要解析这个字段

	RealImageFilePATH string `json:"-"` //不要解析这个字段  书籍为文件夹的时候，实际图片的路径
	ImgType           string `json:"-"` //不要解析这个字段
}

// BookInfo 与Book唯一的区别是没有AllPageInfo,而是封面图URL
type BookInfo struct {
	Name            string         `json:"name"`
	Author          string         `json:"author"`
	Title           string         `json:"title"`
	FilePath        string         `json:"-"` //不要解析这个字段
	ExtractPath     string         `json:"-"` //不要解析这个字段
	AllPageNum      int            `json:"all_page_num"`
	FileType        string         `json:"file_type"`
	FileSize        int64          `json:"file_size"`
	Modified        time.Time      `json:"modified_time"`
	BookID          string         `json:"uuid"` //根据FilePath计算
	IsDir           bool           `json:"is_folder"`
	ExtractNum      int            `json:"extract_num"`
	ExtractComplete bool           `json:"extract_complete"`
	ReadPercent     float64        `json:"read_percent"`
	NonUTF8Zip      bool           `json:"non_utf_8_zip"`
	ZipTextEncoding string         `json:"zip_text_encoding"`
	Cover           SinglePageInfo `json:"cover"`
	//Pages         AllPageInfo `json:"pages"`
}

// NewBookInfo BookInfo的模拟构造函数
func NewBookInfo(b Book) *BookInfo {
	return &BookInfo{
		Name:            b.Name,
		Author:          b.Author,
		Title:           b.Title,
		FilePath:        b.GetFilePath(),
		ExtractPath:     b.ExtractPath,
		AllPageNum:      b.GetAllPageNum(),
		FileType:        b.FileType,
		FileSize:        b.FileSize,
		Modified:        b.Modified,
		BookID:          b.BookID,
		IsDir:           b.IsDir,
		ExtractNum:      b.ExtractNum,
		ExtractComplete: b.ExtractComplete,
		ReadPercent:     b.ReadPercent,
		NonUTF8Zip:      b.NonUTF8Zip,
		Cover:           b.Cover,
	}
}

func GetBookShelf() (*[]BookInfo, error) {
	var bookShelf []BookInfo
	for _, b := range BookList {
		info := NewBookInfo(b)
		bookShelf = append(bookShelf, *info)
	}
	if len(bookShelf) > 0 {
		return &bookShelf, nil
	}
	return nil, errors.New("can not found bookshelf")
}

// InitBook 初始化一本书，设置文件路径、书名、BookID等等
func InitBook(allPageNum int, filePath string, modified time.Time, isDir bool, fileSize int64, extractComplete bool) *Book {
	var b = Book{
		AllPageNum:      allPageNum,
		Modified:        modified,
		IsDir:           isDir,
		FileSize:        fileSize,
		ExtractComplete: extractComplete,
	}
	//书名直接用路径
	b.Name = filePath
	b.SetFilePath(filePath)
	//压缩文件的话，去除路径，取文件名
	if !b.IsDir {
		post := strings.LastIndex(filePath, "/") //Unix路径分隔符
		if post == -1 {
			post = strings.LastIndex(filePath, "\\") //windows分隔符
		}
		if post != -1 {
			//filePath = string([]rune(filePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
			filePath = filePath[post:]
			filePath = strings.ReplaceAll(filePath, "\\", "")
			filePath = strings.ReplaceAll(filePath, "/", "")
		}
		b.Name = filePath
		b.FileType = path.Ext(filePath) //获取文件后缀
	}
	b.setBookID()
	return &b
}

// GetBookByUUID 获取特定书籍，复制一份数据
// TODO: 只获取、不改变原始数据。
func GetBookByUUID(uuid string, sort bool) (Book, error) {
	for _, b := range BookList {
		if b.BookID == uuid {
			if sort {
				b.SortPages()
			}
			return b, nil
		}
	}
	//为了调试方便，支持模糊查找，可以使用UUID的开头来查找书籍，当然这样有可能出错
	for _, b := range BookList {
		if strings.HasPrefix(b.BookID, uuid) {
			if sort {
				b.SortPages()
			}
			return b, nil
		}
	}
	return Book{}, errors.New("can not found book,uuid=" + uuid)
}

// GetBookByAuthor 获取同一作者的书籍。
// TODO: 只获取、不改变原始数据。
func GetBookByAuthor(author string, sort bool) ([]Book, error) {
	var bookList []Book
	for _, b := range BookList {
		if b.Author == author {
			if sort {
				b.SortPages()
			}
			bookList = append(bookList, b)
		}
	}
	if len(bookList) > 0 {
		return bookList, nil
	}
	return nil, errors.New("can not found book,author=" + author)
}

// AllPageInfo Slice
type AllPageInfo []SinglePageInfo

func (s AllPageInfo) Len() int {
	return len(s)
}

// Less 按时间或URL，将图片排序
func (s AllPageInfo) Less(i, j int) (less bool) {
	//如何定义 s[i] < s[j]  根据文件名
	//numI, err1 := getNumberFromString(s[i].NameInArchive)
	//if err1 != nil {
	//	less = strings.Compare(s[i].NameInArchive, s[j].NameInArchive) > 0
	//	return less
	//}
	//numJ, err2 := getNumberFromString(s[j].NameInArchive)
	//if err2 != nil {
	//	less = strings.Compare(s[i].NameInArchive, s[j].NameInArchive) > 0
	//	return less
	//}
	////fmt.Println("numI:",numI)
	////fmt.Println("numJ:",numJ)
	//less = numI < numJ //如果有的话，比较文件名里的数字

	//如何定义 s[i] < s[j]  根据文件名(第三方库、自然语言字符串)
	less = tools.Compare(s[i].NameInArchive, s[j].NameInArchive)

	////如何定义 s[i] < s[j]  根据修改时间
	//if Config.SortImage == "time" {
	//	less = s[i].ModeTime.After(s[j].ModeTime) // s[i] 的年龄（修改时间），是否比 s[j] 小？
	//}
	return less
}

func (s AllPageInfo) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// SortPages 上面三个函数定义好了，终于可以使用sort包排序了
func (b *Book) SortPages() {
	sort.Sort(b.Pages)
}

func md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// setBookID  根据路径的MD5，生成书籍ID。初始化时调用。
func (b *Book) setBookID() {
	//fmt.Println("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	fileAbaPath, err := filepath.Abs(b.filePath)
	if err != nil {
		b.BookID = md5string(b.filePath)
		fmt.Println(err, fileAbaPath)
	} else {
		b.BookID = md5string(b.GetFilePath())
	}
}

// GetBookID  根据路径的MD5，生成书籍ID
func (b *Book) GetBookID() string {
	//防止未初始化，最好不要用到
	if b.BookID == "" {
		fmt.Println("BookID未初始化，一定是哪里写错了")
		b.setBookID()
	}
	return b.BookID
}

//设置页数
func (b *Book) setPageNum() {
	b.AllPageNum = len(b.Pages)
}

func (b *Book) GetAllPageNum() int {
	//设置页数
	b.setPageNum()
	return b.AllPageNum
}

func (b *Book) SetFilePath(path string) {
	fileAbaPath, err := filepath.Abs(path)
	if err != nil {
		//因为权限问题，无法取得绝对路径的情况下，用相对路径
		fmt.Println(err, fileAbaPath)
		b.filePath = path
	} else {
		b.filePath = fileAbaPath
	}
}

func (b *Book) GetFilePath() string {
	return b.filePath
}

func (b *Book) GetName() string { //绑定到Book结构体的方法
	return b.Name
}

func (b *Book) GetPicNum() int {
	var PicNum = 0
	for _, p := range b.Pages {
		if isSupportMedia(p.Url) {
			PicNum++
		}
	}
	return PicNum
}

// ScanAllImage 服务器端分析分辨率、漫画单双页，只适合已解压文件
func (b *Book) ScanAllImage() {
	log.Println(locale.GetString("check_image_start"))
	// Console progress bar
	bar := pb.StartNew(b.GetAllPageNum())
	tmpl := `{{ red "With funcs:" }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}} {{string . "my_green_string" | green}} {{string . "my_blue_string" | blue}}`
	bar.SetTemplateString(tmpl)
	for i := 0; i < len(b.Pages); i++ { //此处不能用range，因为需要修改
		analyzePageImages(&b.Pages[i], b.filePath)
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
	bar := pb.StartNew(b.GetAllPageNum())
	for i := 0; i < len(b.Pages); i++ { //此处不能用range，因为需要修改
		//wg.Add(1)
		count++
		ii := i
		//并发处理，提升图片分析速度
		wp.Do(func() error {
			//defer wg.Done()
			analyzePageImages(&b.Pages[ii], b.filePath)
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

//analyzePageImages 解析漫画的分辨率与blurhash
func analyzePageImages(p *SinglePageInfo, bookPath string) {
	err := p.analyzeImage(bookPath)
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

// analyzeImage 获取某页漫画的分辨率与blurhash
func (i *SinglePageInfo) analyzeImage(bookPath string) (err error) {
	var img image.Image
	//img, err = imaging.Open(i.RealImageFilePATH)

	imgData, err := arch.GetSingleFile(bookPath, i.NameInArchive, "gbk")
	if err != nil {
		fmt.Println(err)
	}
	buf := bytes.NewBuffer(imgData)
	img, err = imaging.Decode(buf)
	if err != nil {
		log.Printf(locale.GetString("check_image_error")+" %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
		//很耗费服务器资源，以后再研究。
		str, err := blurhash.Encode(1, 1, img)
		if err != nil {
			// Handle errors
			log.Printf(locale.GetString("check_image_error")+" %v\n", err)
		}
		i.Blurhash = str
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
