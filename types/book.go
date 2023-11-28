package types

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/yumenaka/comi/logger"
	"image"
	"log"
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
	"github.com/yumenaka/comi/util"
)

var (
	mapBooks     = make(map[string]*Book)      //实际存在的书，通过扫描生成
	mapBookGroup = make(map[string]*BookGroup) //通过分析路径与深度生成的书组。不备份，也不存储到数据库。key是BookID
	MainFolder   = Folder{
		SubFolders: make(map[string]*subFolder),
		SortBy:     "name",
	}
)

// Book 定义书籍，BooID不应该重复，根据文件路径生成
type Book struct {
	BookInfo
	Pages     Pages            `json:"pages"`       //storm:"inline" 内联字段，结构体嵌套时使用
	ChildBook map[string]*Book `json:"child_book" ` //key：BookID
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
	Blurhash          string    `json:"-"`        //`json:"blurhash"` //blurhash占位符。需要扫描图片生成（util.GetImageDataBlurHash）
	Height            int       `json:"-"`        //暂时用不着 这个字段不解析`json:"height"`   //blurhash用，图片高
	Width             int       `json:"-"`        //暂时用不着 这个字段不解析`json:"width"`    //blurhash用，图片宽
	ModeTime          time.Time `json:"-"`        //这个字段不解析
	FileSize          int64     `json:"-"`        //这个字段不解析
	RealImageFilePATH string    `json:"-"`        //这个字段不解析  书籍为文件夹的时候，实际图片的路径
	ImgType           string    `json:"-"`        //这个字段不解析
}

func NewImageInfo(pageNum int, nameInArchive string, url string, fileSize int64) *ImageInfo {
	return &ImageInfo{PageNum: pageNum, NameInArchive: nameInArchive, Url: url, FileSize: fileSize}
}

// CheckBookExist 查看内存中是否已经有了这本书,有了就报错，让调用者跳过
func CheckBookExist(filePath string, bookType SupportFileType, storePath string) bool {
	//如果是文件夹，就不用检查了
	if bookType == TypeDir || bookType == TypeBooksGroup {
		return false
	}

	//实际存在的书，通过扫描生成
	for _, realBook := range mapBooks {
		fileAbaPath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Info(err, fileAbaPath)
			if realBook.FilePath == filePath && realBook.ParentFolder == storePath && realBook.Type == bookType {
				return true
			}
		} else {
			if realBook.FilePath == fileAbaPath && realBook.Type == bookType {
				return true
			}
		}
	}
	return false
}

// New  初始化Book，设置文件路径、书名、BookID等等
func New(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	if CheckBookExist(filePath, bookType, storePath) {
		return nil, errors.New("skip:" + filePath)
	}
	//初始化书籍
	var b = Book{
		BookInfo: BookInfo{
			Author:        "",
			Modified:      modified,
			FileSize:      fileSize,
			InitComplete:  false,
			Depth:         depth,
			BookStorePath: storePath,
			Type:          bookType},
	}
	//设置属性：
	//FilePath，转换为绝对路径
	b.setFilePath(filePath)
	b.setTitle(filePath)
	b.Author, _ = util.GetAuthor(b.Title)
	//设置属性：父文件夹
	b.setParentFolder(filePath)
	b.setBookID()
	return &b, nil
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

// 初始化Book时，设置页数
func (b *Book) setPageNum() {
	b.PageCount = len(b.Pages.Images)
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
		if b.GetPageCount() < minPageNum {
			continue
		}
		err = AddBook(b, basePath, minPageNum)
		if err != nil {
			return err
		}
	}
	return err
}

// RestoreDatabaseBooks 从数据库中读取的书籍信息，放到内存中
func RestoreDatabaseBooks(list []*Book) (err error) {
	for _, b := range list {
		if b.Type == TypeZip || b.Type == TypeRar || b.Type == TypeCbz || b.Type == TypeCbr || b.Type == TypeTar || b.Type == TypeEpub {
			mapBooks[b.BookID] = b
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
	if b.GetPageCount() < minPageNum {
		return errors.New("add book Error：minPageNum = " + strconv.Itoa(b.GetPageCount()))
	}
	if _, ok := MainFolder.SubFolders[basePath]; !ok {
		if err := MainFolder.AddSubFolder(basePath); err != nil {
			logger.Info(err)
		}
	}
	mapBooks[b.BookID] = b
	return MainFolder.AddBookToSubFolder(basePath, &b.BookInfo)
}

// DeleteBookByID 删除一本书
func DeleteBookByID(bookID string) {
	delete(mapBooks, bookID) //如果key存在在删除此数据；如果不存在，delete不进行操作，也不会报错
}

// GetBooksNumber 获取书籍总数，当然不包括BookGroup
func GetBooksNumber() int {
	return len(mapBooks)
}

func GetAllBookInfoList(sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	//首先加上所有真实的书籍
	for _, b := range mapBooks {
		info := NewBaseInfo(b)
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

func GetArchiveBooks() []*Book {
	var list []*Book
	//所有真实书籍
	for _, b := range mapBooks {
		if b.Type == TypeZip || b.Type == TypeRar || b.Type == TypeCbz || b.Type == TypeCbr || b.Type == TypeTar || b.Type == TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBookByID 获取特定书籍，复制一份数据
func GetBookByID(id string, sortBy string) (*Book, error) {
	//根据id查找
	b, ok := mapBooks[id]
	if ok {
		b.SortPages(sortBy)
		return b, nil
	}
	g, ok := mapBookGroup[id]
	if ok {
		temp := Book{
			BookInfo: g.BookInfo,
		}
		return &temp, nil
	}
	return nil, errors.New("can not found book,id=" + id)
}

func GetBookGroupIDByBookID(id string) (string, error) {
	//根据id查找
	for _, group := range mapBookGroup {
		for _, b := range group.ChildBook {
			if b.BookID == id {
				return group.BookID, nil
			}
		}
	}
	return "", errors.New("can not found group,id=" + id)
}

// GetBookByAuthor 获取同一作者的书籍。
func GetBookByAuthor(author string, sortBy string) ([]*Book, error) {
	var bookList []*Book
	for _, b := range mapBooks {
		if b.Author == author {
			b.SortPages(sortBy)
			bookList = append(bookList, b)
		}
	}
	if len(bookList) > 0 {
		return bookList, nil
	}
	return nil, errors.New("can not found book,author=" + author)
}

type Pages struct {
	Images []ImageInfo `json:"images"`
	SortBy string      `json:"sort_by"`
}

func (s Pages) Len() int {
	return len(s.Images)
}

// Less 按时间或URL，将图片排序
func (s Pages) Less(i, j int) (less bool) {
	//如何定义 Images[i] < Images[j]
	switch s.SortBy {
	case "filename": //根据文件名(第三方库、自然语言字符串)
		return util.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	case "filesize": //根据文件大小
		return s.Images[i].FileSize < s.Images[j].FileSize
	case "modify_time": //根据修改时间
		return s.Images[i].ModeTime.Before(s.Images[j].ModeTime) // Images[i] 的修改时间，是否比 Images[j] 晚
	// 如何定义 Images[i] < Images[j](反向)
	case "filename_reverse": //根据文件名(反向)
		return !util.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	case "filesize_reverse": //根据文件大小(反向)
		return !(s.Images[i].FileSize < s.Images[j].FileSize)
	case "modify_time_reverse": //根据修改时间(反向)
		return !s.Images[i].ModeTime.Before(s.Images[j].ModeTime) // Images[i] 的修改时间，是否比 Images[j] 晚
	default: //默认根据文件名
		return util.Compare(s.Images[i].NameInArchive, s.Images[j].NameInArchive)
	}
}

func (s Pages) Swap(i, j int) {
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
	}
	b.setClover() //重新排序后重新设置封面
}

// SortPagesByImageList 根据一个既定的文件列表，重新对页面排序。用于epub文件。
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
		logger.Info(locale.GetString("EPUB_CANNOT_RESORT"), b.FilePath)
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
func (b *BookInfo) setBookID() {
	//logger.Info("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	fileAbaPath, err := filepath.Abs(b.FilePath)
	if err != nil {
		logger.Info(err, fileAbaPath)
	}
	tempStr := b.FilePath + strconv.Itoa(b.ChildBookNum) + strconv.Itoa(int(b.FileSize)) + string(b.Type) + b.ParentFolder + b.BookStorePath
	b62 := base62.EncodeToString([]byte(md5string(md5string(tempStr))))
	b.BookID = getShortBookID(b62, 7)
}

func getShortBookID(fullID string, minLength int) string {
	if len(fullID) <= minLength {
		logger.Info("can not short ID:" + fullID)
		return fullID
	}
	shortID := fullID[0:minLength]
	notFound := true
	add := 0
	pass := false
	for notFound {
		pass = true
		for _, book := range mapBooks {
			if shortID == book.BookID {
				add++
				shortID = fullID[0 : minLength+add]
				pass = false
			}
		}
		for _, group := range mapBookGroup {
			if shortID == group.BookID {
				add++
				shortID = fullID[0 : minLength+add]
				pass = false
			}
		}
		if pass {
			notFound = false
			return shortID
		}
	}
	return fullID
}

// GetBookID  根据路径的MD5，生成书籍ID
func (b *Book) GetBookID() string {
	//防止未初始化，最好不要用到
	if b.BookID == "" {
		logger.Info("BookID未初始化，一定是哪里写错了")
		b.setBookID()
	}
	return b.BookID
}

// GetAuthor  获取作者信息
func (b *Book) GetAuthor() string {
	return b.Author
}

func (b *Book) GetPageCount() int {
	b.setClover()
	if !b.InitComplete {
		//设置页数
		b.setPageNum()
		b.InitComplete = true
	}
	return b.PageCount
}

func (b *Book) GetFilePath() string {
	return b.FilePath
}

func (b *Book) GetName() string { //绑定到Book结构体的方法
	return b.Title
}

// ScanAllImage 服务器端分析分辨率、漫画单双页，只适合已解压文件
func (b *Book) ScanAllImage() {
	log.Println(locale.GetString("check_image_start"))
	// Console progress bar
	bar := pb.StartNew(b.GetPageCount())
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
	bar := pb.StartNew(b.GetPageCount())
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
		logger.Info(err)
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
	//logger.Info(locale.GetString("clear_temp_file_start"))
	for _, tempBook := range mapBooks {
		clearTempFilesOne(debug, cacheFilePath, tempBook)
	}
}

// 清空某一本压缩漫画的解压缓存
func clearTempFilesOne(debug bool, cacheFilePath string, book *Book) {
	//logger.Info(locale.GetString("clear_temp_file_start"))
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
			logger.Info(locale.GetString("clear_temp_file_error") + cachePath)
		} else {
			if debug {
				logger.Info(locale.GetString("clear_temp_file_completed") + cachePath)
			}
		}
	}
}
