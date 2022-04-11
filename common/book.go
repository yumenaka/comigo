package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"path"
	"path/filepath"
	"sort"
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

// Book 定义书籍，BooID不应该重复，根据文件路径生成
type Book struct {
	Name            string           `json:"name" storm:"index"`                              //书名 //storm:"index" 索引字段
	BookID          string           `json:"id"   storm:"id"`                                 //根据FilePath计算 //storm会搜索id或ID做为主键
	FilePath        string           `json:"-" storm:"filepath" storm:"index" storm:"unique"` //storm:"index" 索引字段 storm:"unique" 唯一字段
	BookStorePath   string           `json:"-"    storm:"index"`                              //在哪个子书库
	Type            BookType         `json:"book_type" storm:"index"`                         //可以是书籍组(book_group)、文件夹(dir)、文件后缀( .zip .rar .pdf .mp4)等
	ChildBookNum    int              `json:"child_book_num" storm:"index"`                    //子书籍的数量
	ChildBook       map[string]*Book `json:"child_book" `                                     //key：BookID
	Depth           int              `json:"depth" storm:"index"`                             //文件深度
	ParentFolder    string           `json:"parent_folder" storm:"index"`                     //所在父文件夹
	AllPageNum      int              `json:"all_page_num" storm:"index"`                      //storm:"index" 索引字段
	FileSize        int64            `json:"file_size" storm:"index"`                         //storm:"index" 索引字段
	Cover           SinglePageInfo   `json:"cover" storm:"inline"`                            //storm:"inline" 内联字段，结构体嵌套时使用
	Pages           AllPageInfo      `json:"pages" storm:"inline"`                            //storm:"inline" 内联字段，结构体嵌套时使用
	Author          []string         `json:"-"`                                               //json不解析，启用可改为`json:"author"`
	ISBN            string           `json:"-"`                                               //json不解析，启用可改为`json:"isbn"`
	Press           string           `json:"-"`                                               //json不解析，启用可改为`json:"press"`        //出版社
	PublishedAt     string           `json:"-"`                                               //json不解析，启用可改为`json:"published_at"` //出版日期
	ExtractPath     string           `json:"-"`                                               //json不解析
	Modified        time.Time        `json:"-"`                                               //json不解析，启用可改为`json:"modified_time"`
	ExtractNum      int              `json:"-"`                                               //json不解析，启用可改为`json:"extract_num"`
	ExtractComplete bool             `json:"-"`                                               //json不解析，启用可改为`json:"extract_complete"`
	ReadPercent     float64          `json:"-"`                                               //json不解析，启用可改为`json:"read_percent"`
	NonUTF8Zip      bool             `json:"-"`                                               //json不解析，启用可改为    `json:"non_utf8_zip"`
	ZipTextEncoding string           `json:"-"`                                               //json不解析，启用可改为   `json:"zip_text_encoding"`
}

type BookType string

//书籍类型
const (
	BookTypeDir         BookType = "dir"
	BookTypeZip         BookType = ".zip"
	BookTypeRar         BookType = ".rar"
	BookTypeBooksGroup  BookType = "book_group"
	BookTypeCbz         BookType = ".cbz"
	BookTypeCbr         BookType = ".cbr"
	BookTypeEpub        BookType = ".epub"
	BookTypePDF         BookType = ".pdf"
	BookTypeUnknownFile BookType = "unknown"
)

//SinglePageInfo 单张书页
type SinglePageInfo struct {
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

//NewBook  初始化Book，设置文件路径、书名、BookID等等
func NewBook(filePath string, modified time.Time, fileSize int64, storePath string, depth int) *Book {
	var b = Book{
		Author:          []string{""},
		Modified:        modified,
		FileSize:        fileSize,
		ExtractComplete: false,
		Depth:           depth,
		BookStorePath:   storePath,
	}
	//设置属性：
	//FilePath，转换为绝对路径
	b.setFilePath(filePath)
	//书籍类型
	b.Type = getBookTypeByFilename(filePath)
	b.setName(filePath)
	//设置属性：父文件夹
	b.setParentFolder(filePath)
	b.setBookID()
	return &b
}

//初始化Book时，设置FilePath
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

//初始化Book时，取得BookType
func getBookTypeByFilename(filename string) BookType {
	//获取文件后缀
	switch strings.ToLower(path.Ext(filename)) {
	case ".zip":
		return BookTypeZip
	case ".rar":
		return BookTypeRar
	case ".cbz":
		return BookTypeCbz
	case ".cbr":
		return BookTypeCbr
	case ".epub":
		return BookTypeEpub
	case ".pdf":
		return BookTypePDF
	default:
		return BookTypeUnknownFile
	}
}

func (b *Book) setParentFolder(filePath string) {
	//取得文件所在文件夹的路径
	//如果类型是文件夹，同时最后一个字符是路径分隔符的话，就多取一次dir，移除多余的Unix路径分隔符或windows分隔符
	if b.Type == BookTypeDir {
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
	if b.Type != BookTypeBooksGroup { //不是书籍组(book_group)。
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

//初始化Book时，设置页数
func (b *Book) setPageNum() {
	b.AllPageNum = len(b.Pages)
}

//初始化Book时， 设置封面信息
func (b *Book) setClover() {
	if len(b.Pages) >= 1 {
		b.Cover = b.Pages[0]
	}
}

// AddBooks 添加一组书
func AddBooks(list []*Book, basePath string) (err error) {
	for _, b := range list {
		err = AddBook(b, basePath)
		if err != nil {
			return err
		}
	}
	return err
}

// AddBook 添加一本书
func AddBook(b *Book, basePath string) error {
	if b.BookID == "" {
		return errors.New("add book Error：empty BookID")
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

func GetAllBookInfoList(sort string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		info := NewBookInfo(b)
		infoList.BookInfos = append(infoList.BookInfos, *info)
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBy = sort
		infoList.SortBooks()
		return &infoList, nil
	}
	return nil, errors.New("can not found bookshelf")
}

func GetBookInfoListByDepth(depth int, sort string) (*BookInfoList, error) {
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
		infoList.SortBy = sort
		infoList.SortBooks()
		return &infoList, nil
	}
	return nil, errors.New("can not found bookshelf")
}

func GetBookInfoListByMaxDepth(depth int, sort string) (*BookInfoList, error) {
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
		infoList.SortBy = sort
		infoList.SortBooks()
		return &infoList, nil
	}
	return nil, errors.New("can not found bookshelf")
}

func GetBookInfoListByBookGroupBookID(BookID string, sort string) (*BookInfoList, error) {
	var infoList BookInfoList
	book := mapBookGroups[BookID]
	if book != nil {
		//首先加上所有真实的书籍
		for _, b := range book.ChildBook {
			info := NewBookInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		if len(infoList.BookInfos) > 0 {
			infoList.SortBy = sort
			infoList.SortBooks()
			return &infoList, nil
		}
	}
	return nil, errors.New("can not found bookshelf")
}

// GetBookByID 获取特定书籍，复制一份数据
// TODO: 只获取、不改变原始数据。
func GetBookByID(id string, sort bool) (*Book, error) {
	//根据id查找
	b, ok := mapBooks[id]
	if ok {
		if sort {
			b.SortPages()
		}
		return b, nil
	}
	//为了调试方便，支持模糊查找，可以使用UUID的开头来查找书籍，当然这样有可能出错
	for _, b := range mapBooks {
		if strings.HasPrefix(b.BookID, id) {
			if sort {
				b.SortPages()
			}
			return b, nil
		}
	}
	return nil, errors.New("can not found book,id=" + id)
}

// GetBookByAuthor 获取同一作者的书籍。
// TODO: 只获取、不改变原始数据。
func GetBookByAuthor(author string, sort bool) ([]*Book, error) {
	var bookList []*Book
	for _, b := range mapBooks {
		if len(b.Author) == 0 {
			continue
		}
		if b.Author[0] == author {
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
	b62 := base62.EncodeToString([]byte(md5string(b.FilePath)))
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

func (b *Book) GetAllPageNum() int {
	//设置页数
	b.setPageNum()
	if b.Cover.Url == "" {
		b.setClover()
	}
	return b.AllPageNum
}

func (b *Book) GetFilePath() string {
	return b.FilePath
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
		analyzePageImages(&b.Pages[i], b.FilePath)
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
			analyzePageImages(&b.Pages[ii], b.FilePath)
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
