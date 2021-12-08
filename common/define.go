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

type ServerConfig struct {
	OpenBrowser         bool   `json:"-"` //不要解析这个字段
	DisableLAN          bool   `json:"-"` //不要解析这个字段
	Template            string `json:"template"`
	Auth                string `json:"-"` //不要解析这个字段 访问密码，还没做
	PrintAllIP          bool   `json:"-"` //不要解析这个字段
	Port                int
	ConfigPath          string `json:"-"` //不要解析这个字段
	CheckImageInServer  bool
	DebugMode           bool   `json:"-"` //不要解析这个字段
	LogToFile           bool   `json:"-"` //不要解析这个字段
	LogFilePath         string `json:"-"` //不要解析这个字段
	LogFileName         string `json:"-"` //不要解析这个字段
	MaxDepth            int    `json:"-"` //不要解析这个字段
	MinImageNum         int
	ServerHost          string
	EnableWebpServer    bool
	WebpConfig          WebPServerConfig `json:"-"` //不要解析这个字段
	EnableFrpcServer    bool
	FrpConfig           FrpClientConfig `json:"-"` //不要解析这个字段
	ZipFilenameEncoding string          `json:"-"` //不要解析这个字段
	SketchCountSeconds  int             `json:"sketch_count_seconds"`
	SortImage           string
	TempFolderSetting   string
	CleanOnExit         bool
	CleanNotAll         bool
	//SortByModTime       bool            //SortByModificationTime
	//SortByFileName      bool
}

var Config = ServerConfig{
	OpenBrowser:         true,
	DisableLAN:          false,
	Template:            "multi", //multi、single、random etc.
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
	ServerHost:         "",
	SketchCountSeconds: 90,
	SortImage:          "",
	TempFolderSetting:  "",
	CleanOnExit:        true,
	CleanNotAll:        true,
}

// SetByExecutableFilename 通过执行文件名设置默认网页模板参数
func (config *ServerConfig) SetByExecutableFilename() {
	// 当前执行目录
	targetPath, _ := os.Getwd()
	fmt.Println(locale.GetString("target_path"), targetPath)
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
	//如果执行文件名包含 single 等关键字，选择 single 分页漫画模板
	if haveKeyWord(ExtFileName, []string{"single", "单页", "シングル"}) {
		config.Template = "single"
	}
	//如果执行文件名包含 double 等关键字，选择 double 分页漫画模板
	if haveKeyWord(ExtFileName, []string{"double", "双页", "ダブルページ"}) {
		config.Template = "double"
	}
	//选择模式以后，打印提示
	switch config.Template {
	case "scroll":
		fmt.Println(locale.GetString("scroll_template"))
	case "sketch":
		fmt.Println(locale.GetString("sketch_template"))
		//速写倒计时秒数
		fmt.Println(locale.GetString("COMI_SKETCH_COUNT_SECONDS"), config.SketchCountSeconds)
	case "single":
		fmt.Println(locale.GetString("single_page_template"))
	case "double":
		fmt.Println(locale.GetString("double_page_template"))
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
	RealExtractPath  string
	WebImagePath     string
	Version          = "v0.2.4"
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
	AllPageNum      int         `json:"all_page_num"`
	PageInfo        AllPageInfo `json:"pages"`
	FileType        string      `json:"file_type"`
	FileSize        int64       `json:"file_size"`
	Modified        time.Time   `json:"modified_time"`
	FileID          string      `json:"uuid"`
	IsFolder        bool        `json:"is_folder"`
	ExtractNum      int         `json:"extract_num"`
	ExtractComplete bool        `json:"extract_complete"`
	ReadPercent     float64     `json:"read_percent"`
}

type SinglePageInfo struct {
	ModeTime  time.Time `json:"-"` //不要解析这个字段
	FileSize  int64     `json:"-"` //不要解析这个字段
	Height    int       `json:"height"`
	Width     int       `json:"width"`
	Url       string    `json:"url"`
	LocalPath string    `json:"-"` //不要解析这个字段
	Name      string    `json:"-"` //不要解析这个字段
	ImgType   string    `json:"image_type"`
}

// AllPageInfo Slice
type AllPageInfo []SinglePageInfo

func (s AllPageInfo) Len() int {
	return len(s)
}

// Less 按时间或URL，将图片排序
func (s AllPageInfo) Less(i, j int) (less bool) {
	//如何定义 s[i] < s[j]  根据文件名
	numI, err1 := getNumberFromString(s[i].Name)
	if err1 != nil {
		less = strings.Compare(s[i].Name, s[j].Name) > 0
		return less
	}
	numJ, err2 := getNumberFromString(s[j].Name)
	if err2 != nil {
		less = strings.Compare(s[i].Name, s[j].Name) > 0
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

func (b *Book) SetFileID() {
	fileAbaPath, err := filepath.Abs(b.FilePath)
	//fmt.Println("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	if err != nil {
		fmt.Println(err, fileAbaPath)
	}
	//fmt.Println(md5s(fileAbaPath))
	b.FileID = md5string(fileAbaPath)
}

// SetArchiveBookName  绑定到Book结构体的方法
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
	extractNum := 0
	Percent := 0
	tempPercent := 0
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
	for i := 0; i < count; i++ {
		extractNum++
		if b.AllPageNum != 0 {
			Percent = int((float32(extractNum) / float32(b.AllPageNum)) * 100)
			if tempPercent != Percent {
				if (Percent%20) == 0 || Percent == 10 {
					//fmt.Print(strconv.Itoa(Percent) + "% ")
				}
			}
			tempPercent = Percent
		}
	}
	// finish bar
	bar.Finish()
	log.Println(locale.GetString("check_image_completed"))
}

func SetImageType(p *SinglePageInfo) {
	err := p.GetImageSize()
	//log.Println(locale.GetString("check_image_ing"), p.LocalPath)
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
	img, err = imaging.Open(i.LocalPath)
	if err != nil {
		log.Printf(locale.GetString("check_image_error")+" %v\n", err)
	} else {
		i.Width = img.Bounds().Dx()
		i.Height = img.Bounds().Dy()
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
		if Config.CleanOnExit {
			fmt.Println("\r" + locale.GetString("start_clear_file"))
			if Config.CleanNotAll {
				deleteTempFilesByFileID(ReadingBook.FileID)
			} else {
				deleteAllTempFiles()
			}
		}
		os.Exit(0)
	}()
}
func InitReadingBook() (err error) {
	//准备解压，设置图片文件夹
	if ReadingBook.IsFolder {
		WebImagePath = ReadingBook.FilePath
		ReadingBook.ExtractComplete = true
		ReadingBook.ExtractNum = ReadingBook.AllPageNum
	} else {
		SetTempDir()
		WebImagePath = path.Join(RealExtractPath, ReadingBook.FileID) //extraFolder
		err = LsArchive(&ReadingBook)
		if err != nil {
			fmt.Println(locale.GetString("scan_archive_error"))
			return err
		}
		err = UnArchive(&ReadingBook)
		if err != nil {
			fmt.Println(locale.GetString("un_archive_error"))
			return err
		}
		ReadingBook.SetArchiveBookName(ReadingBook.FilePath) //设置书名
	}
	//服务器分析图片
	if Config.CheckImageInServer {
		ReadingBook.ScanAllImageGo() //扫描所有图片，取得分辨率信息，使用了协程
	}
	//服务器排序图片
	if Config.SortImage != "" {
		if Config.SortImage == "name" {
			ReadingBook.SortPages()
			fmt.Println(locale.GetString("COMI_SORT_BY_NAME"))
		}
		if Config.SortImage == "time" {
			ReadingBook.SortPages()
			fmt.Println(locale.GetString("COMI_SORT_BY_TIME"))
		}
		if Config.DebugMode {
			//判断是否已经排好顺序，将会打印true
			fmt.Println("IS Sorted?\t", sort.IsSorted(ReadingBook.PageInfo))
			//打印排序后的数据
			//litter.Dump(ReadingBook.PageInfo)
		}
	}
	return err
}

// SetTempDir 设置临时文件夹，退出时会被清理
func SetTempDir() {
	//手动设置的临时文件夹
	if Config.TempFolderSetting != "" && tools.ChickExists(Config.TempFolderSetting) && tools.ChickIsDir(Config.TempFolderSetting) {
		RealExtractPath = path.Join(Config.TempFolderSetting, "comigo_cache_files")
	} else {
		RealExtractPath = path.Join(os.TempDir(), "comigo_cache_files") //直接使用系统文件夹
	}
	err := os.MkdirAll(RealExtractPath, os.ModePerm)
	if err != nil {
		println(locale.GetString("temp_folder_error"))
	} else {
		fmt.Println(locale.GetString("temp_folder_path") + RealExtractPath)
	}
}

func deleteAllTempFiles() {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	if strings.Contains(RealExtractPath, "comigo_cache_files") { //判断文件夹前缀，避免删错文件
		err := os.RemoveAll(RealExtractPath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + RealExtractPath)
		} else {
			fmt.Println(locale.GetString("clear_temp_file_completed") + RealExtractPath)
		}
	}
}

func deleteTempFilesByFileID(UUID string) {
	fmt.Println(locale.GetString("clear_temp_file_start"))
	if strings.Contains(RealExtractPath, "comigo_cache_files") { //判断文件夹前缀，避免删错文件
		clearPath := path.Join(RealExtractPath, UUID)
		err := os.RemoveAll(clearPath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + clearPath)
		} else {
			fmt.Println(locale.GetString("clear_temp_file_completed") + clearPath)
		}
	}
}

////根据权限，清理文件可能失败
//func deleteOldTempFiles() {
//	tempDirUpperFolder := RealExtractPath
//	post := strings.LastIndex(RealExtractPath, "/") //Unix风格的路径分隔符
//	if post == -1 {
//		post = strings.LastIndex(RealExtractPath, "\\") //windows风格的分隔符
//	}
//	if post != -1 {
//		tempDirUpperFolder = string([]rune(RealExtractPath)[:post]) //为了防止中文字符被错误截断，先转换成rune，再转回来
//		fmt.Println(locale.GetString("temp_folder_path"), tempDirUpperFolder)
//	}
//	files, err := ioutil.ReadDir(tempDirUpperFolder)
//	if err != nil {
//		fmt.Println(err)
//	}
//	for _, fi := range files {
//		if fi.IsDir() {
//			oldTempDir := tempDirUpperFolder + "/" + fi.Name()
//			if strings.Contains(oldTempDir, "comic_cache_A8cG") { //判断文件夹前缀，避免删错文件
//				err := os.RemoveAll(oldTempDir)
//				if err != nil {
//					fmt.Println(locale.GetString("clear_temp_file_error") + oldTempDir)
//				} else {
//					fmt.Println(locale.GetString("clear_temp_file_completed") + oldTempDir)
//				}
//			}
//		}
//	}
//}
