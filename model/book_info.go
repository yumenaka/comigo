package model

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/logger"
)

// BookInfo 与 Book 唯一的区别是没有 AllPageInfo，而是封面图 URL，减小 JSON 文件的大小
type BookInfo struct {
	Author          string          `json:"author"`         // 作者
	BookID          string          `json:"id"`             // 根据 FilePath 生成的唯一 ID
	BookStorePath   string          `json:"-"`              // 在哪个子书库
	ChildBookNum    int             `json:"child_book_num"` // 子书籍数量
	Cover           MediaFileInfo   `json:"cover"`          // 封面图
	Deleted         bool            `json:"deleted"`        // 源文件是否已删除
	Depth           int             `json:"depth"`          // 书籍深度
	ExtractPath     string          `json:"-"`              // 解压路径，7z 用，JSON 不解析
	ExtractNum      int             `json:"-"`              // 文件解压数
	FileSize        int64           `json:"file_size"`      // 文件大小
	FilePath        string          `json:"-"`              // 文件绝对路径，JSON 不解析
	ISBN            string          `json:"isbn"`           // ISBN
	InitComplete    bool            `json:"-"`              // 是否解压完成
	Modified        time.Time       `json:"modified_time"`  // 修改时间
	NonUTF8Zip      bool            `json:"-"`              // 是否为特殊编码 zip
	PageCount       int             `json:"page_count"`     // 总页数
	ParentFolder    string          `json:"parent_folder"`  // 父文件夹
	Press           string          `json:"press"`          // 出版社
	PublishedAt     string          `json:"published_at"`   // 出版日期
	ReadPercent     float64         `json:"read_percent"`   // 阅读进度
	Title           string          `json:"title"`          // 书名
	Type            SupportFileType `json:"type"`           // 书籍类型
	ZipTextEncoding string          `json:"-"`              // zip 文件编码
}

// NewBaseInfo 创建新的 BookInfo 实例
func NewBaseInfo(b *Book) *BookInfo {
	return &BookInfo{
		Author:          b.Author,
		BookID:          b.BookID,
		BookStorePath:   b.BookStorePath,
		Cover:           b.Cover,
		ChildBookNum:    b.ChildBookNum,
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

// initBookID 根据路径的 MD5，初始化书籍 ID
func (b *BookInfo) initBookID() *BookInfo {
	// 生成 BookID 的字符串
	tempStr := b.FilePath + strconv.Itoa(int(b.FileSize)) + string(b.Type) + b.ParentFolder + b.BookStorePath
	// 两次 MD5 加密，然后转为 base62 编码
	// 为什么选择 Base62?
	// 1. 人类可读，可以目视或简单的 regexp 进行验证
	// 2. 仅包含字母数字符号，不包含特殊字符
	// 3. 可以通过在任何文本编辑器和浏览器地址栏中双击鼠标来完全选择
	// 4. 紧凑，生成的字符串比 Base32 短
	b62 := base62.EncodeToString([]byte(md5string(md5string(tempStr))))
	b.BookID = getShortBookID(b62, 7)
	return b
}

// SetCover 设置封面
func (b *BookInfo) SetCover(c MediaFileInfo) *BookInfo {
	b.Cover = c
	return b
}

// setFilePath 初始化 Book 时，设置 FilePath
func (b *BookInfo) setFilePath(path string) *BookInfo {
	fileAbsPath, err := filepath.Abs(path)
	if err != nil {
		// 因为权限问题，无法取得绝对路径的情况下，用相对路径
		logger.Info(err, fileAbsPath)
		b.FilePath = path
	} else {
		b.FilePath = fileAbsPath
	}
	return b
}

// setParentFolder 设置父文件夹
func (b *BookInfo) setParentFolder(filePath string) *BookInfo {
	dirPath := filePath
	if b.Type == TypeDir {
		dirPath = strings.TrimRight(filePath, "/\\")
	}
	parentDir := filepath.Dir(dirPath)
	b.ParentFolder = filepath.Base(parentDir)
	return b
}

// setAuthor 设置作者
func (b *BookInfo) setAuthor() *BookInfo {
	b.Author = util.GetAuthor(b.Title)
	return b
}

// setTitle 设置标题
func (b *BookInfo) setTitle(filePath string) *BookInfo {
	if b.Type != TypeBooksGroup {
		b.Title = filepath.Base(filePath)
	} else {
		b.Title = filePath
	}
	return b
}

// ShortTitle 返回简短的标题
func (b *BookInfo) ShortTitle() string {
	shortTitle := b.Title
	// 使用预编译正则表达式进行替换
	shortTitle = re1.ReplaceAllString(shortTitle, "")
	shortTitle = re2.ReplaceAllString(shortTitle, "")
	shortTitle = re3.ReplaceAllString(shortTitle, "")
	shortTitle = domainReg.ReplaceAllString(shortTitle, "")
	shortTitle = re4.ReplaceAllString(shortTitle, "")
	shortTitle = re5.ReplaceAllString(shortTitle, "")
	// [过度简化]如果标题长度小于 2，返回前 15 个字符
	if len([]rune(shortTitle)) < 2 && len([]rune(b.Title)) > 15 {
		return string([]rune(b.Title)[:15]) + "…"
	}
	// [简化标题]如果标题长度小于等于 15，返回标题
	if len([]rune(shortTitle)) <= 15 {
		return shortTitle
	}
	// [简化不完全]如果标题长度大于 15，返回前 15 个字符
	return string([]rune(shortTitle)[:15]) + "…"
}

// 预编译正则表达式
var (
	re1       = regexp.MustCompile(`[\[(（【][A-Za-z0-9_\-×\s+\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}]+`)
	re2       = regexp.MustCompile(`[]）】)]`)
	re3       = regexp.MustCompile(`\.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)`)
	domainReg = regexp.MustCompile(`^(((ht|f)tps?)://)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}/?`)
	re4       = regexp.MustCompile(`^\s`)
	re5       = regexp.MustCompile(`^[\-` + "`" + `~!@#$^&*()=|{}':;'@#￥……&*（）——|{}‘；：”“'。，、？]`)
)

// GetBookInfoListByDepth 根据深度获取书籍列表
func GetBookInfoListByDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	// 首先加上所有真实的书籍
	for _, bookValue := range mapBooks.Range {
		b := bookValue.(*Book)
		if b.Depth == depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	// 接下来还要加上扫描生成出来的书籍组
	for _, folderValue := range MainStore.SubStores.Range {
		bs := folderValue.(*subStore)
		for _, groupValue := range bs.BookGroupMap.Range {
			group := groupValue.(*BookInfo)
			if group.Depth == depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetBookInfoListByDepth")
}

// GetBookInfoListByMaxDepth 获取指定最大深度的书籍列表
func GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	// 首先加上所有真实的书籍
	for _, bookValue := range mapBooks.Range {
		b := bookValue.(*Book)
		if b.Depth <= depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	// 扫描生成的书籍组
	for _, folderValue := range MainStore.SubStores.Range {
		bs := folderValue.(*subStore)
		for _, groupValue := range bs.BookGroupMap.Range {
			group := groupValue.(*BookInfo)
			if group.Depth <= depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error: cannot find bookshelf in GetBookInfoListByMaxDepth")
}

// TopOfShelfInfo 获取顶层书架信息
func TopOfShelfInfo(sortBy string) (*BookInfoList, error) {
	//if len(*LocalStores) == 0 {
	//	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
	//}
	//if len(*LocalStores) > 1 {
	//	// 有多个书库
	//	var infoList BookInfoList
	//	for _, localPath := range *LocalStores {
	//		for _, groupValue := range mapBookGroup.Range {
	//			group := groupValue.(*BookGroup)
	//			if group.BookInfo.ParentFolder == localPath {
	//				infoList.BookInfos = append(infoList.BookInfos, group.BookInfo)
	//			}
	//		}
	//	}
	//	if len(infoList.BookInfos) > 0 {
	//		infoList.SortBooks(sortBy)
	//		return &infoList, nil
	//	}
	//	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
	//}
	// 显示顶层书库的书籍
	var infoList BookInfoList
	for _, bookValue := range mapBooks.Range {
		b := bookValue.(*Book)
		if b.Depth == 0 {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	for _, groupValue := range mapBookGroup.Range {
		group := groupValue.(*BookGroup)
		if group.BookInfo.Depth == 0 {
			infoList.BookInfos = append(infoList.BookInfos, group.BookInfo)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	// 没找到任何书
	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
}

// GetBookInfoListByID 根据 ID 获取书籍列表
func GetBookInfoListByID(BookID string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	groupValue, ok := mapBookGroup.Load(BookID)
	if ok {
		tempGroup := groupValue.(*BookGroup)
		for _, bookValue := range tempGroup.ChildBook.Range {
			b := bookValue.(*BookInfo)
			infoList.BookInfos = append(infoList.BookInfos, *b)
		}
		if len(infoList.BookInfos) > 0 {
			infoList.SortBooks(sortBy)
			return &infoList, nil
		}
	}
	return nil, errors.New("cannot find BookInfo，ID：" + BookID)
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	for _, bookValue := range mapBooks.Range {
		b := bookValue.(*Book)
		if b.ParentFolder == parentFolder {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
	}
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("cannot find book, parentFolder=" + parentFolder)
}
