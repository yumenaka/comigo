package model

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// Book 定义书籍结构
type Book struct {
	BookInfo  // 嵌入 BookInfo 结构体
	BookMarks // 书签信息
	PageInfos // 书籍内所有页面的信息
}

// GuestCover 猜测书籍的封面
func (b *Book) GuestCover() (cover PageInfo) {
	// 封面图片的命名规则
	for i := range b.PageInfos {
		// 先转换为小写
		filenameLower := strings.ToLower(b.PageInfos[i].Name)
		// 再去掉后缀名
		filenameWithoutExt := strings.TrimSuffix(filenameLower, path.Ext(filenameLower))
		// 再去掉前置的0 ，例如00001 -> 1, 0 -> ""
		filenameTrimmed := strings.TrimLeft(filenameWithoutExt, "0")
		// 对原始不带前导0的文件名包含 "cover" 的检查
		// 检查文件名（去除后缀和前导0）是否包含 "cover" 或等于 "" (原为 "0") 或 "1"
		if strings.Contains(filenameWithoutExt, "cover") || filenameTrimmed == "" || filenameTrimmed == "1" {
			cover = b.PageInfos[i] // 获取实际元素的指针
			return cover           // 找到封面，停止循环
		}
	}
	// 如果通过名称规则没有找到封面，并且书至少有一页，则使用第一页作为封面
	if len(b.PageInfos) > 0 {
		cover = b.PageInfos[0]
	}
	return cover // 返回找到的封面或空值
}

// NewBook 初始化 Book，设置文件路径、书名、BookID 等
func NewBook(bookPath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	// 初始化书籍
	book := &Book{
		BookInfo: BookInfo{
			Modified:     modified,
			FileSize:     fileSize,
			InitComplete: false,
			Depth:        depth,
			StoreUrl:     storePath,
			Type:         bookType,
		},
	}
	// 设置文件路径、书名、BookID
	_, err := book.setFilePath(bookPath).setParentFolder(bookPath).setTitle(bookPath).setAuthor().initBookID(bookPath)
	return book, err
}

// setPageNum 设置书籍的页数
func (b *Book) setPageNum() {
	b.PageCount = len(b.PageInfos)
}

// SortPages 对页面进行排序
func (b *Book) SortPages(s string) {
	if b.Type == TypeEpub && s == "default" {
		return
	}
	if s != "" {
		b.PageInfos.SortImages(s)
	}
}

// SortPagesByImageList 根据给定的文件列表排序页面（用于 EPUB）
func (b *Book) SortPagesByImageList(imageList []string) {
	if len(imageList) == 0 {
		return
	}
	imageInfos := b.PageInfos
	var reSortList []PageInfo
	for _, imgName := range imageList {
		for _, imgInfo := range imageInfos {
			if imgInfo.Name == imgName {
				reSortList = append(reSortList, imgInfo)
				break
			}
		}
	}
	if len(reSortList) == 0 {
		logger.Infof(locale.GetString("epub_cannot_resort")+"%s", b.BookPath)
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
	b.PageInfos = reSortList
}

// GetStoreID 获取编码后的书库ID，StoreID是书库URL路径的 base62 编码
func (b *Book) GetStoreID() string {
	return base62.EncodeToString([]byte(b.StoreUrl))
}

// GetPageCount 获取页数
func (b *Book) GetPageCount() int {
	if !b.InitComplete {
		b.setPageNum()
		b.InitComplete = true
	}
	return b.PageCount
}

// GetAllBooksNumber  获取书籍总数，不包括 BookGroup 类型
func GetAllBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.Type == TypeBooksGroup {
			continue // 跳过书组类型
		}
		count++
	}
	return count
}

// ClearBookNotExist  检查内存中的书的源文件是否存在，不存在就删掉
func ClearBookNotExist() {
	logger.Infof("Checking book files exist...")
	var deletedBooks []string
	// 遍历所有书籍
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, book := range allBooks {
		// 如果父文件夹存在，但书籍文件不存在，也说明这本书被删除了
		if _, err := os.Stat(book.BookPath); os.IsNotExist(err) {
			deletedBooks = append(deletedBooks, book.BookPath)
			err := IStore.DeleteBook(book.BookID)
			if err != nil {
				logger.Infof("Error deleting book %s: %s", book.BookID, err)
			}
		}
	}
	// 重新生成书组
	if len(deletedBooks) > 0 {
		if err := IStore.GenerateBookGroup(); err != nil {
			logger.Infof("Error initializing main folder: %s", err)
		}
	}
}
