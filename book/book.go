package book

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bbrks/go-blurhash"
	"github.com/cheggaaa/pb/v3"
	"github.com/disintegration/imaging"
	"github.com/jxskiss/base62"
	"github.com/xxjwxc/gowp/workpool"

	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/tools"
)

var (
	mapBooks      = make(map[string]*Book) //实际存在的书，通过扫描生成
	mapBookGroups = make(map[string]*Book) //通过分析路径与深度生成的书组
	Stores        = Bookstores{
		mapBookstores: make(map[string]*singleBookstore),
		SortBy:        "name",
	}
)

// Book 定义书籍，BooID不应该重复，根据文件路径生成
type Book struct {
	Name            string           `json:"name"` //书名
	BookID          string           `json:"id"`   //根据FilePath+BookType+修改时间+filesize等等计算，bookID应该唯一
	FilePath        string           `json:"-" storm:"filepath"`
	BookStorePath   string           `json:"-"   `           //在哪个子书库
	Type            SupportFileType  `json:"book_type"`      //可以是书籍组(book_group)、文件夹(dir)、文件后缀( .zip .rar .pdf .mp4)等
	ChildBookNum    int              `json:"child_book_num"` //子书籍的数量
	ChildBook       map[string]*Book `json:"child_book" `    //key：BookID
	Depth           int              `json:"depth"`          //文件深度
	ParentFolder    string           `json:"parent_folder"`  //所在父文件夹
	AllPageNum      int              `json:"all_page_num"`   //storm:"index" 索引字段
	FileSize        int64            `json:"file_size"`      //storm:"index" 索引字段
	Cover           ImageInfo        `json:"cover"`          //storm:"inline" 内联字段，结构体嵌套时使用
	Pages           AllPageInfo      `json:"pages"`          //storm:"inline" 内联字段，结构体嵌套时使用
	Author          []string         `json:"-"`              //json不解析，启用可改为`json:"author"`
	ISBN            string           `json:"-"`              //json不解析，启用可改为`json:"isbn"`
	Press           string           `json:"-"`              //json不解析，启用可改为`json:"press"`        //出版社
	PublishedAt     string           `json:"-"`              //json不解析，启用可改为`json:"published_at"` //出版日期
	ExtractPath     string           `json:"-"`              //json不解析
	Modified        time.Time        `json:"-"`              //json不解析，启用可改为`json:"modified_time"`
	ExtractNum      int              `json:"-"`              //json不解析，启用可改为`json:"extract_num"`
	InitComplete    bool             `json:"-"`              //json不解析，启用可改为`json:"extract_complete"`
	ReadPercent     float64          `json:"-"`              //json不解析，启用可改为`json:"read_percent"`
	NonUTF8Zip      bool             `json:"-"`              //json不解析，启用可改为    `json:"non_utf8_zip"`
	ZipTextEncoding string           `json:"-"`              //json不解析，启用可改为   `json:"zip_text_encoding"`
}

type SupportFileType string

// 书籍类型
const (
	TypeDir         SupportFileType = "dir"
	TypeZip         SupportFileType = ".zip"
	TypeRar         SupportFileType = ".rar"
	TypeBooksGroup  SupportFileType = "book_group"
	TypeCbz         SupportFileType = ".cbz"
	TypeCbr         SupportFileType = ".cbr"
	TypeTar         SupportFileType = ".tar"
	TypeEpub        SupportFileType = ".epub"
	TypePDF         SupportFileType = ".pdf"
	TypeVideo       SupportFileType = "video"
	TypeAudio       SupportFileType = "audio"
	TypeUnknownFile SupportFileType = "unknown"
)

// ImageInfo 单张书页
type ImageInfo struct {
	PageNum           int       `json:"-"`        //这个字段不解析
	NameInArchive     string    `json:"filename"` //用于解压的压缩文件内文件路径，或图片名，为了适应特殊字符，经过一次转义
	Url               string    `json:"url"`      //远程用户读取图片的URL，为了适应特殊字符，经过一次转义
	Blurhash          string    `json:"-"`        //`json:"blurhash"` //blurhash占位符。需要扫描图片生成（tools.GetImageDataBlurHash）
	Height            int       `json:"-"`        //暂时用不着 这个字段不解析`json:"height"`   //blurhash用，图片的高
	Width             int       `json:"-"`        //暂时用不着 这个字段不解析`json:"width"`    //blurhash用，图片的宽
	ModeTime          time.Time `json:"-"`        //这个字段不解析
	FileSize          int64     `json:"-"`        //这个字段不解析
	RealImageFilePATH string    `json:"-"`        //这个字段不解析  书籍为文件夹的时候，实际图片的路径
	ImgType           string    `json:"-"`        //这个字段不解析
}

func NewImageInfo(pageNum int, nameInArchive string, url string, fileSize int64) *ImageInfo {
	return &ImageInfo{PageNum: pageNum, NameInArchive: nameInArchive, Url: url, FileSize: fileSize}
}

// New  初始化Book，设置文件路径、书名、BookID等等
func New(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	//查看内存中是否已经有了这本书,有了就报错，让调用者跳过
	for _, realBook := range mapBooks {
		fileAbaPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Println(err, fileAbaPath)
			if realBook.FilePath == filePath && realBook.ParentFolder == storePath {
				return nil, errors.New("Duplicate books:" + filePath)
			}
		} else {
			if realBook.FilePath == fileAbaPath {
				return nil, errors.New("Duplicate books:" + fileAbaPath)
			}
		}
	}
	for _, groupBook := range mapBookGroups {
		fileAbaPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Println(err, fileAbaPath)
			if groupBook.FilePath == filePath && groupBook.ParentFolder == storePath {
				return nil, errors.New("Duplicate books:" + filePath)
			}
		} else {
			if groupBook.FilePath == fileAbaPath {
				return nil, errors.New("Duplicate books:" + fileAbaPath)
			}
		}
	}
	//初始化书籍
	var b = Book{
		Author:        []string{""},
		Modified:      modified,
		FileSize:      fileSize,
		InitComplete:  false,
		Depth:         depth,
		BookStorePath: storePath,
		Type:          bookType,
	}
	//设置属性：
	//FilePath，转换为绝对路径
	b.setFilePath(filePath)
	b.setName(filePath)
	//设置属性：父文件夹
	b.setParentFolder(filePath)
	b.setBookID()
	return &b, nil
}

// 初始化Book时，设置FilePath
func (b *Book) setFilePath(path string) {
	fileAbaPath, err := filepath.Abs(path)
	if err != nil {
		//因为权限问题，无法取得绝对路径的情况下，用相对路径
		fmt.Println(err, fileAbaPath)
		b.FilePath = path
	} else {
		b.FilePath = fileAbaPath
	}
}

// GetBookTypeByFilename 初始化Book时，取得BookType
func GetBookTypeByFilename(filename string) SupportFileType {
	//获取文件后缀
	switch strings.ToLower(path.Ext(filename)) {
	case ".zip":
		return TypeZip
	case ".rar":
		return TypeRar
	case ".cbz":
		return TypeCbz
	case ".cbr":
		return TypeCbr
	case ".epub":
		return TypeEpub
	case ".tar":
		return TypeTar
	case ".pdf":
		return TypePDF
	case ".mp4", ".m4v", ".flv", ".avi", ".webm":
		return TypeVideo
	case ".mp3", ".wav", ".wma", ".ogg":
		return TypeAudio
	default:
		return TypeUnknownFile
	}
}

func (b *Book) setParentFolder(filePath string) {
	//取得文件所在文件夹的路径
	//如果类型是文件夹，同时最后一个字符是路径分隔符的话，就多取一次dir，移除多余的Unix路径分隔符或windows分隔符
	if b.Type == TypeDir {
		if filePath[len(filePath)-1] == '/' || filePath[len(filePath)-1] == '\\' {
			filePath = filepath.Dir(filePath)
		}
	}
	folder := filepath.Dir(filePath)
	post := strings.LastIndex(folder, "/") //Unix路径分隔符
	if post == -1 {
		post = strings.LastIndex(folder, "\\") //windows分隔符
	}
	if post != -1 {
		//FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
		p := folder[post:]
		p = strings.ReplaceAll(p, "\\", "")
		p = strings.ReplaceAll(p, "/", "")
		b.ParentFolder = p
	}
}

func (b *Book) setName(filePath string) {
	b.Name = filePath
	//设置属性：书籍名，取文件后缀(可能为 .zip .rar .pdf .mp4等等)。
	if b.Type != TypeBooksGroup { //不是书籍组(book_group)。
		post := strings.LastIndex(filePath, "/") //Unix路径分隔符
		if post == -1 {
			post = strings.LastIndex(filePath, "\\") //windows分隔符
		}
		if post != -1 {
			//FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
			name := filePath[post:]
			name = strings.ReplaceAll(name, "\\", "")
			name = strings.ReplaceAll(name, "/", "")
			b.Name = name
		}
	}
}

// 初始化Book时，设置页数
func (b *Book) setPageNum() {
	b.AllPageNum = len(b.Pages.Images)
}

// 初始化Book时， 设置封面信息
func (b *Book) setClover() {
	if len(b.Pages.Images) >= 1 {
		b.Cover = b.Pages.Images[0]
	}
}

// AddBooks 添加一组书
func AddBooks(list []*Book, basePath string, minPageNum int) (err error) {
	for _, b := range list {
		if b.GetAllPageNum() < minPageNum {
			continue
		}
		err = AddBook(b, basePath, minPageNum)
		if err != nil {
			return err
		}
	}
	return err
}

// AddBook 添加一本书
func AddBook(b *Book, basePath string, minPageNum int) error {
	//没有初始化BookID
	if b.BookID == "" {
		return errors.New("add book Error：empty BookID")
	}
	//页数不符合要求
	if b.GetAllPageNum() < minPageNum {
		return errors.New("add book Error：minPageNum = " + strconv.Itoa(b.GetAllPageNum()))
	}
	if _, ok := Stores.mapBookstores[basePath]; !ok {
		if err := Stores.NewSingleBookstore(basePath); err != nil {
			fmt.Println(err)
		}
	}
	mapBooks[b.BookID] = b
	return Stores.AddBookToStores(basePath, b)
}

// DeleteBookByID 删除一本书
func DeleteBookByID(bookID string) {
	delete(mapBooks, bookID) //如果key存在在删除此数据；如果不存在，delete不进行操作，也不会报错
}

// GetBooksNumber 获取书籍总数，当然不包括BookGroup
func GetBooksNumber() int {
	return len(mapBooks)
}

// GetRandomBook 随机选取一本书
func GetRandomBook() *Book {
	if len(mapBooks) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano()) //随机种子，否则每回都会一样
	randNum := rand.Intn(100) % len(mapBooks)
	start := 0
	for _, b := range mapBooks {
		if randNum == start {
			return b
		}
		start++
	}
	return nil
}

func GetAllBookInfoList(sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		info := NewBookInfo(b)
		infoList.BookInfos = append(infoList.BookInfos, *info)
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetAllBookInfoList")
}

func GetAllBookList() []*Book {
	var list []*Book
	//加上所有真实书籍
	for _, b := range mapBooks {
		list = append(list, b)
	}
	return list
}

func GetBookInfoListByDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		if b.Depth == depth {
			info := NewBookInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	//接下来还要加上扫描生成出来的书籍组
	for _, bs := range Stores.mapBookstores {
		for _, b := range bs.BookGroupMap {
			if b.Depth == depth {
				info := NewBookInfo(b)
				infoList.BookInfos = append(infoList.BookInfos, *info)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByDepth")
}

func GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		if b.Depth <= depth {
			info := NewBookInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	//接下来还要加上扫描生成出来的书籍组
	for _, bs := range Stores.mapBookstores {
		for _, b := range bs.BookGroupMap {
			if b.Depth <= depth {
				info := NewBookInfo(b)
				infoList.BookInfos = append(infoList.BookInfos, *info)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByMaxDepth")
}

func GetBookInfoListByBookGroupBookID(BookID string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	book := mapBookGroups[BookID]
	if book != nil {
		//首先加上所有真实的书籍
		for _, b := range book.ChildBook {
			info := NewBookInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		if len(infoList.BookInfos) > 0 {
			infoList.SortBooks(sortBy)
			return &infoList, nil
		}
	}
	return nil, errors.New("can not found bookshelf")
}

// GetBookByID 获取特定书籍，复制一份数据
func GetBookByID(id string, sortBy string) (*Book, error) {
	//根据id查找
	b, ok := mapBooks[id]
	if ok {
		b.SortPages(sortBy)
		return b, nil
	}
	//为了调试方便，支持模糊查找，可以使用UUID的开头来查找书籍，当然这样有可能出错
	for _, b := range mapBooks {
		if strings.HasPrefix(b.BookID, id) {
			b.SortPages(sortBy)
			return b, nil
		}
	}
	return nil, errors.New("can not found book,id=" + id)
}

// GetBookByAuthor 获取同一作者的书籍。
func GetBookByAuthor(author string, sortBy string) ([]*Book, error) {
	var bookList []*Book
	for _, b := range mapBooks {
		if len(b.Author) == 0 {
			continue
		}
		if b.Author[0] == author {
			b.SortPages(sortBy)
			bookList = append(bookList, b)
		}
	}
	if len(bookList) > 0 {
		return bookList, nil
	}
	return nil, errors.New("can not found book,author=" + author)
}

// AllPageInfo
type AllPageInfo struct {
	Images []ImageInfo `json:"images"`
	SortBy string      `json:"sort_by"`
}

func (s AllPageInfo) Len() int {
	return len(s.Images)
}

// Less 按时间或URL，将图片排序
func (s AllPageInfo) Less(i, j int) (less bool) {
	//如何定义 Images[i] < Images[j]
	switch s.SortBy {
	case "filename": //根据文件名(第三方库、自然语言字符串)
		return tools.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	case "filesize": //根据文件大小
		return s.Images[i].FileSize < s.Images[j].FileSize
	case "modify_time": //根据修改时间
		return s.Images[i].ModeTime.Before(s.Images[j].ModeTime) // Images[i] 的修改时间，是否比 Images[j] 晚
	// 如何定义 Images[i] < Images[j](反向)
	case "filename_reverse": //根据文件名(反向)
		return !tools.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	case "filesize_reverse": //根据文件大小(反向)
		return !(s.Images[i].FileSize < s.Images[j].FileSize)
	case "modify_time_reverse": //根据修改时间(反向)
		return !s.Images[i].ModeTime.Before(s.Images[j].ModeTime) // Images[i] 的修改时间，是否比 Images[j] 晚
	default: //默认根据文件名
		return tools.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	}
}

func (s AllPageInfo) Swap(i, j int) {
	s.Images[i], s.Images[j] = s.Images[j], s.Images[i]
}

// SortPages 上面三个函数定义好了，终于可以使用sort包排序了
func (b *Book) SortPages(s string) {
	if b.Type == TypeEpub && s == "default" {
		return
	}
	if s != "" {
		b.Pages.SortBy = s
		sort.Sort(b.Pages)
		//fmt.Println("sort_by:" + s)
	}
	b.setClover() //重新排序后重新设置封面
}

// 根据一个既定的文件列表，重新对页面排序。用于epub文件。
func (b *Book) SortPagesByImageList(imageList []string) {
	if len(imageList) == 0 {
		return
	}
	imageInfos := b.Pages.Images
	//如果在有序表中，按照有序表的顺序重排
	var reSortList []ImageInfo
	for i := 0; i < len(imageList); i++ {
		checkSrc := imageList[i]
		for j := 0; j < len(imageInfos); j++ {
			if imageInfos[j].NameInArchive == checkSrc {
				reSortList = append(reSortList, imageInfos[j])
			}
		}
	}
	if len(reSortList) == 0 {
		fmt.Println("can not resort by epub metadata!")
		return
	}
	//不在表中的话，就不改变顺序，并加在有序表的后面
	for i := 0; i < len(imageInfos); i++ {
		checkName := imageInfos[i].NameInArchive
		find := false
		for j := 0; j < len(imageList); j++ {
			if imageList[j] == checkName {
				find = true
			}
		}
		if !find {
			reSortList = append(reSortList, imageInfos[i])
		}
	}
	b.Pages.Images = reSortList
	b.setClover() //重新排序后重新设置封面
}

func md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}

// setBookID  根据路径的MD5，生成书籍ID。初始化时调用。
func (b *Book) setBookID() {
	//fmt.Println("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	fileAbaPath, err := filepath.Abs(b.FilePath)
	if err != nil {
		fmt.Println(err, fileAbaPath)
	}
	tempStr := b.FilePath + strconv.Itoa(b.ChildBookNum) + strconv.Itoa(int(b.FileSize)) + string(b.Type) + b.Modified.String()
	b62 := base62.EncodeToString([]byte(md5string(tempStr)))
	b.BookID = getShortBookID(b62, 5)
}

func getShortBookID(fullID string, minLength int) string {
	if len(fullID) <= minLength {
		fmt.Println("can not short ID:" + fullID)
		return fullID
	}
	shortID := fullID
	//最短为5位，最长等于全长
	for i := minLength; i <= len(fullID); i++ {
		canUse := true
		for key := range mapBooks {
			if strings.HasPrefix(key, fullID[0:i]) {
				canUse = false
			}
		}
		for key := range mapBookGroups {
			if strings.HasPrefix(key, fullID[0:i]) {
				canUse = false
			}
		}
		if canUse {
			shortID = fullID[0:i]
			break
		}
	}
	return shortID
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

// GetAuthor  获取作者信息
func (b *Book) GetAuthor() string {
	//防止未初始化，最好不要用到
	if len(b.Author) == 0 {
		return ""
	}
	return b.Author[0]
}

func (b *Book) GetAllPageNum() int {
	b.setClover()
	if !b.InitComplete {
		//设置页数
		b.setPageNum()
		b.InitComplete = true
	}
	return b.AllPageNum
}

func (b *Book) GetFilePath() string {
	return b.FilePath
}

func (b *Book) GetName() string { //绑定到Book结构体的方法
	return b.Name
}

// ScanAllImage 服务器端分析分辨率、漫画单双页，只适合已解压文件
func (b *Book) ScanAllImage() {
	log.Println(locale.GetString("check_image_start"))
	// Console progress bar
	bar := pb.StartNew(b.GetAllPageNum())
	tmpl := `{{ red "With funcs:" }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}} {{string . "my_green_string" | green}} {{string . "my_blue_string" | blue}}`
	bar.SetTemplateString(tmpl)
	for i := 0; i < len(b.Pages.Images); i++ { //此处不能用range，因为会修改b.Pages.Images本身
		analyzePageImages(&b.Pages.Images[i], b.FilePath)
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
	for i := 0; i < len(b.Pages.Images); i++ { //此处不能用range，因为会修改b.Pages.Images本身
		//wg.Add(1)
		count++
		ii := i
		//并发处理，提升图片分析速度
		wp.Do(func() error {
			//defer wg.Done()
			analyzePageImages(&b.Pages.Images[ii], b.FilePath)
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

// analyzePageImages 解析漫画的分辨率与blurhash
func analyzePageImages(p *ImageInfo, bookPath string) {
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
func (i *ImageInfo) analyzeImage(bookPath string) (err error) {
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

// ClearTempFilesALL web加载时保存的临时图片，在在退出后清理
func ClearTempFilesALL(debug bool, cacheFilePath string) {
	//fmt.Println(locale.GetString("clear_temp_file_start"))
	for _, tempBook := range mapBooks {
		clearTempFilesOne(debug, cacheFilePath, tempBook)
	}
}

// 清空某一本压缩漫画的解压缓存
func clearTempFilesOne(debug bool, cacheFilePath string, book *Book) {
	//fmt.Println(locale.GetString("clear_temp_file_start"))
	haveThisBook := false
	for _, tempBook := range mapBooks {
		if tempBook.GetBookID() == book.GetBookID() {
			haveThisBook = true
		}
	}
	if haveThisBook {
		cachePath := path.Join(cacheFilePath, book.GetBookID())
		err := os.RemoveAll(cachePath)
		if err != nil {
			fmt.Println(locale.GetString("clear_temp_file_error") + cachePath)
		} else {
			if debug {
				fmt.Println(locale.GetString("clear_temp_file_completed") + cachePath)
			}
		}
	}
}
