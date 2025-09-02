package model

import (
	"errors"
	"path"
	"strings"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
)

// Book 定义书籍结构
type Book struct {
	BookInfo
	Pages `json:"pages"`
}

// GetBookInfo 创建新的 BookInfo 实例
func (b *Book) GetBookInfo() *BookInfo {
	return &BookInfo{
		Author:          b.Author,
		BookID:          b.BookID,
		Cover:           b.GuestCover(), // 使用 GuestCover 方法获取封面
		BookStorePath:   b.BookStorePath,
		ChildBooksNum:   b.ChildBooksNum,
		ChildBooksID:    b.ChildBooksID,
		Deleted:         b.Deleted,
		Depth:           b.Depth,
		ExtractPath:     b.ExtractPath,
		ExtractNum:      b.ExtractNum,
		FilePath:        b.GetFilePath(),
		FileSize:        b.FileSize,
		ISBN:            b.ISBN,
		InitComplete:    b.InitComplete,
		Modified:        b.Modified,
		NonUTF8Zip:      b.NonUTF8Zip,
		PageCount:       b.GetPageCount(),
		ParentFolder:    b.ParentFolder,
		Press:           b.Press,
		PublishedAt:     b.PublishedAt,
		ReadPercent:     b.ReadPercent,
		Type:            b.Type,
		Title:           b.Title,
		ZipTextEncoding: b.ZipTextEncoding,
	}
}

// GuestCover 猜测书籍的封面
func (b *Book) GuestCover() (cover MediaFileInfo) {
	// 封面图片的命名规则
	for i := range b.Images {
		// 先转换为小写
		filenameLower := strings.ToLower(b.Images[i].Name)
		// 再去掉后缀名
		filenameWithoutExt := strings.TrimSuffix(filenameLower, path.Ext(filenameLower))
		// 再去掉前置的0 ，例如00001 -> 1, 0 -> ""
		filenameTrimmed := strings.TrimLeft(filenameWithoutExt, "0")
		// 对原始不带前导0的文件名包含 "cover" 的检查
		// 检查文件名（去除后缀和前导0）是否包含 "cover" 或等于 "" (原为 "0") 或 "1"
		if strings.Contains(filenameWithoutExt, "cover") || filenameTrimmed == "" || filenameTrimmed == "1" {
			cover = b.Images[i] // 获取实际元素的指针
			return cover        // 找到封面，停止循环
		}
	}
	// 如果通过名称规则没有找到封面，并且书至少有一页，则使用第一页作为封面
	if len(b.Images) > 0 {
		cover = b.Images[0]
	}
	return cover // 返回找到的封面或空值
}

// NewBook 初始化 Book，设置文件路径、书名、BookID 等
func NewBook(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*Book, error) {
	if MainStoreGroup.CheckRawFileExist(filePath, bookType) {
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

// NewBookInfo   初始化BookGroup，设置文件路径、书名、BookID等等
func NewBookInfo(filePath string, modified time.Time, fileSize int64, storePath string, depth int, bookType SupportFileType) (*BookInfo, error) {
	// 初始化书籍
	bookInfo := BookInfo{
		Modified:      modified,
		FileSize:      fileSize,
		InitComplete:  false,
		Depth:         depth,
		BookStorePath: storePath,
		Type:          bookType,
	}
	// 设置属性：
	bookInfo.setTitle(filePath).setFilePath(filePath).SetAuthor().setParentFolder(filePath).initBookID()
	return &bookInfo, nil
}

// setPageNum 设置书籍的页数
func (b *Book) setPageNum() {
	b.PageCount = len(b.Images)
}

// SortPages 对页面进行排序
func (b *Book) SortPages(s string) {
	if b.Type == TypeEpub && s == "default" {
		return
	}
	if s != "" {
		b.SortImages(s)
	}
}

// SortPagesByImageList 根据给定的文件列表排序页面（用于 EPUB）
func (b *Book) SortPagesByImageList(imageList []string) {
	if len(imageList) == 0 {
		return
	}
	imageInfos := b.Images
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
	b.Images = reSortList
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
