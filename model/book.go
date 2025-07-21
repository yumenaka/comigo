package model

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/util/logger"
)

// GetBookGroupByBookID 通过数组 ID 获取书组信息
func GetBookGroupByBookID(id string) (*BookGroup, error) {
	if value, ok := mapBookGroup.Load(id); ok {
		return value.(*BookGroup), nil
	}
	return nil, errors.New("GetBookGroupByBookID：cannot find book group, id=" + id)
}

// GetBookGroupIDByBookID 通过子书籍 ID 获取所属书组 ID
func GetBookGroupIDByBookID(id string) (string, error) {
	for _, value := range mapBookGroup.Range {
		group := value.(*BookGroup)
		for _, v := range group.ChildBook.Range {
			b := v.(*BookInfo)
			if b.BookID == id {
				return group.BookID, nil
			}
		}
	}
	return "", errors.New("cannot find group, id=" + id)
}

// ClearAllBookData  清空所有书籍与虚拟书组数据
func ClearAllBookData() {
	ClearBookData()
	ClearBookGroupData()
}

// ClearBookData 清空书籍数据
func ClearBookData() {
	mapBooks.Clear()
}

// ClearBookGroupData  清空书组相关数据
func ClearBookGroupData() {
	// Clear 会删除所有条目，从而产生一个空的 Map。
	mapBookGroup.Clear()
	MainStore.SubStores.Clear()
}

// ResetBookGroupData 重置虚拟书库
func ResetBookGroupData() {
	ClearBookGroupData()
	if err := MainStore.AnalyzeStore(); err != nil {
		logger.Infof("Error initializing main folder: %s", err)
	}
}

// Book 定义书籍结构
type Book struct {
	BookInfo
	Pages Pages `json:"pages"`
}

// GuestCover 猜测书籍的封面
func (b *Book) GuestCover() (cover MediaFileInfo) {
	// 封面图片的命名规则
	for i := range b.Pages.Images {
		// 先转换为小写
		filenameLower := strings.ToLower(b.Pages.Images[i].Name)
		// 再去掉后缀名
		filenameWithoutExt := strings.TrimSuffix(filenameLower, path.Ext(filenameLower))
		// 再去掉前置的0 ，例如00001 -> 1, 0 -> ""
		filenameTrimmed := strings.TrimLeft(filenameWithoutExt, "0")
		// 对原始不带前导0的文件名包含 "cover" 的检查
		// 检查文件名（去除后缀和前导0）是否包含 "cover" 或等于 "" (原为 "0") 或 "1"
		if strings.Contains(filenameWithoutExt, "cover") || filenameTrimmed == "" || filenameTrimmed == "1" {
			cover = b.Pages.Images[i] // 获取实际元素的指针
			return cover              // 找到封面，停止循环
		}
	}
	// 如果通过名称规则没有找到封面，并且书至少有一页，则使用第一页作为封面
	if len(b.Pages.Images) > 0 {
		cover = b.Pages.Images[0]
	}
	return cover // 返回找到的封面或空值
}

// CheckBookExist 查看内存中是否已经有了这本书，有了就返回 true，让调用者跳过
func CheckBookExist(filePath string, bookType SupportFileType) bool {
	if bookType == TypeDir || bookType == TypeBooksGroup {
		return false
	}
	for _, value := range mapBooks.Range {
		realBook := value.(*Book)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Infof("Error getting absolute path: %v", err)
			continue
		}
		if realBook.FilePath == absFilePath && realBook.Type == bookType {
			return true
		}
	}
	return false
}

// CheckAllBookFileExist 检查内存中的书的源文件是否存在，不存在就删掉
func CheckAllBookFileExist() {
	logger.Infof("Checking book files exist...")
	var deletedBooks []string
	// 遍历所有书籍
	for _, value := range mapBooks.Range {
		book := value.(*Book)
		if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
			deletedBooks = append(deletedBooks, book.FilePath)
			DeleteBookByID(book.BookID)
		}
	}
	// 删除不存在的书组
	if len(deletedBooks) > 0 {
		ResetBookGroupData()
	}
}

// NewBook 初始化 Book，设置文件路径、书名、BookID 等
func NewBook(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	if CheckBookExist(filePath, bookType) {
		return nil, errors.New("skip: " + filePath)
	}
	// 初始化书籍
	book := &Book{
		BookInfo: BookInfo{
			Modified:      modified,
			FileSize:      fileSize,
			InitComplete:  false,
			Depth:         depth,
			BookStorePath: storePath,
			Type:          bookType,
		},
	}
	// 设置文件路径、书名、BookID
	book.setFilePath(filePath).setParentFolder(filePath).setTitle(filePath).SetAuthor().initBookID()
	return book, nil
}

// setPageNum 设置书籍的页数
func (b *Book) setPageNum() {
	b.PageCount = len(b.Pages.Images)
}

// RestoreDatabaseBooks 从数据库中读取的书籍信息，放到内存中
func RestoreDatabaseBooks(list []*Book) {
	for _, b := range list {
		if b.Type == TypeZip || b.Type == TypeRar || b.Type == TypeCbz || b.Type == TypeCbr || b.Type == TypeTar || b.Type == TypeEpub {
			mapBooks.Store(b.BookID, b)
		}
	}
}

// SortPages 对页面进行排序
func (b *Book) SortPages(s string) {
	if b.Type == TypeEpub && s == "default" {
		return
	}
	if s != "" {
		b.Pages.SortBy = s
		sort.Sort(b.Pages)
	}
}

// SortPagesByImageList 根据给定的文件列表排序页面（用于 EPUB）
func (b *Book) SortPagesByImageList(imageList []string) {
	if len(imageList) == 0 {
		return
	}
	imageInfos := b.Pages.Images
	var reSortList []MediaFileInfo
	for _, imgName := range imageList {
		for _, imgInfo := range imageInfos {
			if imgInfo.Name == imgName {
				reSortList = append(reSortList, imgInfo)
				break
			}
		}
	}
	if len(reSortList) == 0 {
		logger.Infof(locale.GetString("epub_cannot_resort")+"%s", b.FilePath)
		return
	}
	// 添加不在列表中的图片
	for _, imgInfo := range imageInfos {
		found := false
		for _, imgName := range imageList {
			if imgInfo.Name == imgName {
				found = true
				break
			}
		}
		if !found {
			reSortList = append(reSortList, imgInfo)
		}
	}
	b.Pages.Images = reSortList
}

// getShortBookID 生成短的 BookID，避免冲突
func getShortBookID(fullID string, minLength int) string {
	if len(fullID) <= minLength {
		logger.Infof("Cannot shorten ID: %s", fullID)
		return fullID
	}
	shortID := fullID[:minLength]
	add := 0
	for {
		conflict := false
		for _, value := range mapBooks.Range {
			b := value.(*Book)
			if b.BookID == shortID {
				conflict = true
				break
			}
		}
		for _, value := range mapBookGroup.Range {
			group := value.(*BookGroup)
			if group.BookID == shortID {
				conflict = true
				break
			}
		}
		if !conflict {
			break
		}
		add++
		if minLength+add > len(fullID) {
			break
		}
		shortID = fullID[:minLength+add]
	}
	return shortID
}

// GetBookID 获取书籍的 ID
func (b *Book) GetBookID() string {
	if b.BookID == "" {
		logger.Infof("BookID 未初始化，可能存在错误")
		b.initBookID()
	}
	return b.BookID
}

// GetAuthor 获取作者信息
func (b *Book) GetAuthor() string {
	return b.Author
}

// GetPageCount 获取页数
func (b *Book) GetPageCount() int {
	if !b.InitComplete {
		b.setPageNum()
		b.InitComplete = true
	}
	return b.PageCount
}

// GetFilePath 获取文件路径
func (b *Book) GetFilePath() string {
	return b.FilePath
}

// analyzePageImages 解析漫画的分辨率与类型
func analyzePageImages(p *MediaFileInfo, bookPath string) {
	err := p.analyzeImage(bookPath)
	if err != nil {
		logger.Infof(locale.GetString("check_image_error") + err.Error())
		return
	}
	if p.Width == 0 && p.Height == 0 {
		p.ImgType = "Unknown"
		return
	}
	if p.Width > p.Height {
		p.ImgType = "DoublePage"
	} else {
		p.ImgType = "SinglePage"
	}
}

// ClearTempFilesALL 清理所有缓存的临时图片
func ClearTempFilesALL(debug bool, cacheFilePath string) {
	for _, value := range mapBooks.Range {
		book := value.(*Book)
		clearTempFilesOne(debug, cacheFilePath, book)
	}
}

// clearTempFilesOne 清理某一本书的缓存
func clearTempFilesOne(debug bool, cacheFilePath string, book *Book) {
	cachePath := path.Join(cacheFilePath, book.GetBookID())
	err := os.RemoveAll(cachePath)
	if err != nil {
		logger.Infof("Error clearing temp files: %s", cachePath)
	} else if debug {
		logger.Infof("Cleared temp files: %s", cachePath)
	}
}
