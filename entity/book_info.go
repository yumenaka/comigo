package entity

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comi/util"
	"github.com/yumenaka/comi/util/logger"
)

// BookInfo 与Book唯一的区别是没有AllPageInfo,而是封面图URL 减小 json文件的大小
type BookInfo struct {
	Title           string          `json:"title"`          // 书名
	BookID          string          `json:"id"`             // 根据FilePath生成的唯一ID
	BookStorePath   string          `json:"-"   `           // 在哪个子书库
	Type            SupportFileType `json:"type"`           // 书籍类型
	Depth           int             `json:"depth"`          // 书籍深度
	ChildBookNum    int             `json:"child_book_num"` // 子书籍数量
	PageCount       int             `json:"page_count"`     // 总页数
	Cover           ImageInfo       `json:"cover"`          // 封面图
	ParentFolder    string          `json:"parent_folder"`  // 父文件夹
	Author          string          `json:"author"`         // 作者
	ISBN            string          `json:"isbn"`           // ISBN
	Press           string          `json:"press"`          // 出版社
	PublishedAt     string          `json:"published_at"`   // 出版日期
	FilePath        string          `json:"-"`              // 文件绝对路径，json不解析
	FileSize        int64           `json:"file_size"`      // 文件大小
	Modified        time.Time       `json:"modified_time"`  // 修改时间
	ReadPercent     float64         `json:"read_percent"`   // 阅读进度
	ExtractPath     string          `json:"-"`              // 解压路径，7z用，json不解析
	ExtractNum      int             `json:"-"`              // 文件解压数    extract_num
	InitComplete    bool            `json:"-"`              // 是否解压完成  extract_complete
	NonUTF8Zip      bool            `json:"-"`              // 是否为特殊编码zip
	ZipTextEncoding string          `json:"-"`              // zip文件编码
}

// NewBaseInfo 模拟构造函数
func NewBaseInfo(b *Book) *BookInfo {
	// 需要单独先执行这个，来设定封面
	pageCount := b.GetPageCount()
	return &BookInfo{
		Title:        b.Title,
		Author:       b.Author,
		Depth:        b.Depth,
		ISBN:         b.ISBN,
		FilePath:     b.GetFilePath(),
		ExtractPath:  b.ExtractPath,
		PageCount:    pageCount,
		Type:         b.Type,
		ChildBookNum: b.ChildBookNum,
		FileSize:     b.FileSize,
		Modified:     b.Modified,
		BookID:       b.BookID,
		ExtractNum:   b.ExtractNum,
		InitComplete: b.InitComplete,
		ReadPercent:  b.ReadPercent,
		NonUTF8Zip:   b.NonUTF8Zip,
		Cover:        b.Cover,
		ParentFolder: b.ParentFolder,
	}
}

// initBookID  根据路径的MD5，初始化书籍ID。
func (b *BookInfo) initBookID() *BookInfo {
	// logger.Infof("文件绝对路径："+fileAbaPath, "路径的md5："+md5string(fileAbaPath))
	fileAbaPath, err := filepath.Abs(b.FilePath)
	if err != nil {
		logger.Info(err, fileAbaPath)
	}
	tempStr := b.FilePath + strconv.Itoa(b.ChildBookNum) + strconv.Itoa(int(b.FileSize)) + string(b.Type) + b.ParentFolder + b.BookStorePath
	b62 := base62.EncodeToString([]byte(md5string(md5string(tempStr))))
	b.BookID = getShortBookID(b62, 7)
	return b
}

// SetClover 设置封面
func (b *BookInfo) SetClover(c ImageInfo) *BookInfo {
	b.Cover = c
	return b
}

// 初始化Book时，设置FilePath
func (b *BookInfo) setFilePath(path string) *BookInfo {
	fileAbaPath, err := filepath.Abs(path)
	if err != nil {
		// 因为权限问题，无法取得绝对路径的情况下，用相对路径
		logger.Info(err, fileAbaPath)
		b.FilePath = path
	} else {
		b.FilePath = fileAbaPath
	}
	return b
}

func (b *BookInfo) setParentFolder(filePath string) *BookInfo {
	// 取得文件所在文件夹的路径
	// 如果类型是文件夹，同时最后一个字符是路径分隔符的话，就多取一次dir，移除多余的Unix路径分隔符或windows分隔符
	if b.Type == TypeDir {
		if filePath[len(filePath)-1] == '/' || filePath[len(filePath)-1] == '\\' {
			filePath = filepath.Dir(filePath)
		}
	}
	folder := filepath.Dir(filePath)
	post := strings.LastIndex(folder, "/") // Unix路径分隔符
	if post == -1 {
		post = strings.LastIndex(folder, "\\") // windows分隔符
	}
	if post != -1 {
		// FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
		p := folder[post:]
		p = strings.ReplaceAll(p, "\\", "")
		p = strings.ReplaceAll(p, "/", "")
		b.ParentFolder = p
	}
	return b
}

func (b *BookInfo) setAuthor() *BookInfo {
	b.Author = util.GetAuthor(b.Title)
	return b
}

func (b *BookInfo) setTitle(filePath string) *BookInfo {
	b.Title = filePath
	// 设置属性：书籍名，取文件后缀(可能为 .zip .rar .pdf .mp4等等)。
	if b.Type != TypeBooksGroup { // 不是书籍组(book_group)。
		post := strings.LastIndex(filePath, "/") // Unix路径分隔符
		if post == -1 {
			post = strings.LastIndex(filePath, "\\") // windows分隔符
		}
		if post != -1 {
			// FilePath = string([]rune(FilePath)[post:]) //为了防止中文字符被错误截断，先转换成rune，再转回来
			filename := filePath[post:]
			filename = strings.ReplaceAll(filename, "\\", "")
			filename = strings.ReplaceAll(filename, "/", "")
			b.Title = filename
		}
	}
	return b
}

func (b *BookInfo) ShortTitle() string {
	shortTitle := b.Title
	// 使用 Go 的正则表达式替换掉一些字符串
	re1 := regexp.MustCompile(`[\[\(（【][A-Za-z0-9_\-×\s+\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}]+`)
	shortTitle = re1.ReplaceAllString(shortTitle, "")

	re2 := regexp.MustCompile(`[\]）】\)]`)
	shortTitle = re2.ReplaceAllString(shortTitle, "")

	re3 := regexp.MustCompile(`\.(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)`)
	shortTitle = re3.ReplaceAllString(shortTitle, "")

	domainReg := regexp.MustCompile(`^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}\/?`)
	shortTitle = domainReg.ReplaceAllString(shortTitle, "")

	re4 := regexp.MustCompile(`^[\s]`)
	shortTitle = re4.ReplaceAllString(shortTitle, "")

	re5 := regexp.MustCompile(`^[\-` + "`" + `~!@#$^&*()=|{}':;'@#￥……&*（）——|{}‘；：”“'。，、？]`)
	shortTitle = re5.ReplaceAllString(shortTitle, "")
	// rune 切片
	if len([]rune(shortTitle)) <= 15 {
		return shortTitle
	}
	return string([]rune(shortTitle[:15])) + "…"
}

func GetBookInfoListByDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	// 首先加上所有真实的书籍
	mapBooks.Range(func(_, value interface{}) bool {
		b := value.(*Book)
		if b.Depth == depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		return true
	})

	// 接下来还要加上扫描生成出来的书籍组
	MainFolder.SubFolders.Range(func(_, value interface{}) bool {
		bs := value.(*subFolder)
		bs.BookGroupMap.Range(func(key, value interface{}) bool {
			group := value.(*BookInfo)
			if group.Depth == depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
			return true
		})
		return true
	})

	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByDepth")
}

func GetBookInfoListByMaxDepth(depth int, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	// 首先加上所有真实的书籍
	mapBooks.Range(func(_, value interface{}) bool {
		b := value.(*Book)
		if b.Depth <= depth {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		return true
	})

	// 扫描生成的书籍组
	MainFolder.SubFolders.Range(func(_, value interface{}) bool {
		bs := value.(*subFolder)
		bs.BookGroupMap.Range(func(key, value interface{}) bool {
			group := value.(*BookInfo)
			if group.Depth <= depth {
				infoList.BookInfos = append(infoList.BookInfos, *group)
			}
			return true
		})
		return true
	})

	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByMaxDepth")
}

func TopOfShelfInfo(sortBy string) (*BookInfoList, error) {
	// 顶层书组
	topGroupBookNum := 0
	// 扫描生成的书籍组
	var top0 BookInfoList
	MainFolder.SubFolders.Range(func(_, value interface{}) bool {
		bs := value.(*subFolder)
		bs.BookGroupMap.Range(func(key, value interface{}) bool {
			group := value.(*BookInfo)
			if group.Depth == 0 {
				topGroupBookNum++
				top0.BookInfos = append(top0.BookInfos, *group)
			}
			return true
		})
		return true
	})
	// 如果顶层书组数量大于1，说明有多个顶层书库。只显示顶层书库。
	if topGroupBookNum > 1 {
		top0.SortBooks(sortBy)
		return &top0, nil
	}
	// 如果只有一个顶层书库，显示真实的书籍
	var infoList BookInfoList
	mapBooks.Range(func(_, value interface{}) bool {
		b := value.(*Book)
		if b.Depth == 1 {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		return true
	})
	mapBookGroup.Range(func(_, value interface{}) bool {
		group := value.(*BookGroup)
		if group.BookInfo.Depth == 1 {
			infoList.BookInfos = append(infoList.BookInfos, group.BookInfo)
		}
		return true
	})
	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("error:can not found bookshelf. GetBookInfoListByMaxDepth")
}

func GetBookInfoListByID(BookID string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	group, ok := mapBookGroup.Load(BookID)
	if ok {
		tempGroup := group.(*BookGroup)
		// 首先加上所有真实的书籍
		tempGroup.ChildBook.Range(func(key, value interface{}) bool {
			b := value.(*BookInfo)
			infoList.BookInfos = append(infoList.BookInfos, *b)
			return true
		})

		if len(infoList.BookInfos) > 0 {
			infoList.SortBooks(sortBy)
			return &infoList, nil
		}
	}
	return nil, errors.New("can not found bookshelf")
}

func GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error) {
	var infoList BookInfoList
	mapBooks.Range(func(_, value interface{}) bool {
		b := value.(*Book)
		if b.ParentFolder == parentFolder {
			info := NewBaseInfo(b)
			infoList.BookInfos = append(infoList.BookInfos, *info)
		}
		return true
	})

	if len(infoList.BookInfos) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	return nil, errors.New("can not found book,parentFolder=" + parentFolder)
}
