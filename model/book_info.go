package model

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// BookInfo 与 Book 唯一的区别是没有 AllPageInfo，而是封面图 URL，减小 JSON 文件的大小
type BookInfo struct {
	Author          string          `json:"author"`          // 作者
	BookID          string          `json:"id"`              // 根据 FilePath 生成的唯一 ID
	BookStorePath   string          `json:"-"`               // 在哪个子书库
	ChildBooksNum   int             `json:"child_books_num"` // 子书籍数量
	ChildBooksID    []string        `json:"child_books_id"`  // 子书籍BookID
	Cover           MediaFileInfo   `json:"cover"`           // 封面图
	Deleted         bool            `json:"deleted"`         // 源文件是否已删除
	Depth           int             `json:"depth"`           // 书籍深度
	ExtractPath     string          `json:"-"`               // 解压路径，7z 用，JSON 不解析
	ExtractNum      int             `json:"-"`               // 文件解压数
	FileSize        int64           `json:"file_size"`       // 文件大小
	FilePath        string          `json:"-"`               // 文件绝对路径，JSON 不解析
	ISBN            string          `json:"isbn"`            // ISBN
	InitComplete    bool            `json:"-"`               // 是否解压完成
	Modified        time.Time       `json:"modified_time"`   // 修改时间
	NonUTF8Zip      bool            `json:"-"`               // 是否为特殊编码 zip
	PageCount       int             `json:"page_count"`      // 总页数
	LastReadPage    int             `json:"last_read_page"`  // 阅读进度
	BookMarks       BookMarks       `json:"book_marks"`      // 书签列表
	ParentFolder    string          `json:"parent_folder"`   // 父文件夹
	Press           string          `json:"press"`           // 出版社
	PublishedAt     string          `json:"published_at"`    // 出版日期
	Title           string          `json:"title"`           // 书名
	Type            SupportFileType `json:"type"`            // 书籍类型
	ZipTextEncoding string          `json:"-"`               // zip 文件编码
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
	b62 := base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(tempStr))))
	// 生成短的 BookID，并避免冲突
	fullID := b62
	minLength := 7
	if len(fullID) <= minLength {
		logger.Infof("Cannot shorten ID: %s", fullID)
		b.BookID = fullID
	}
	shortID := fullID[:minLength]
	add := 0
	for {
		conflict := false
		for _, b := range IStore.ListBooks() {
			if b.BookID == shortID {
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
	b.BookID = shortID
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

// SetAuthor 设置作者
func (b *BookInfo) SetAuthor() *BookInfo {
	b.Author = tools.GetAuthor(b.Title)
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

var (
	// 只删除结尾处的常见扩展名（忽略大小写）
	reExt = regexp.MustCompile(`\.(?i)(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)$`)

	// 去除各种括号及其内容（非贪婪）
	reRound         = regexp.MustCompile(`\([^()]*?\)`)  // 匹配 ()
	reSquare        = regexp.MustCompile(`\[[^\[\]]*?]`) // 匹配 []
	reChineseRound  = regexp.MustCompile(`（[^（）]*?）`)    // 匹配 （）
	reChineseSquare = regexp.MustCompile(`【[^【】]*?】`)    // 匹配 【】

	// 如果只想移除开头的 domain 就保留 ^；想全局替换就去掉 ^
	domainReg = regexp.MustCompile(`^(((ht|f)tps?)://)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}/?`)

	// 去除开头的所有空白
	reLeadingSpace = regexp.MustCompile(`^\s+`)
	// 去除结尾的所有空白
	reTrailingSpace = regexp.MustCompile(`\s+$`)

	// 去除开头的一连串标点符号 (移除括号)
	reLeadingPunct = regexp.MustCompile(`^[\-` + "`" + `~!@#$^&*=|{}':;'@#￥……&*——|{}‘；：”“'。，、？]+`)
)

// ShortName 返回简短的标题（文件名）
func (b *BookInfo) ShortName() string {
	shortTitle := b.Title

	// 1. 移除常见文件扩展名 (忽略大小写)
	shortTitle = reExt.ReplaceAllString(shortTitle, "")

	// 2. 顺序移除所有括号及内部描述
	shortTitle = reRound.ReplaceAllString(shortTitle, "")         // 移除 ()
	shortTitle = reSquare.ReplaceAllString(shortTitle, "")        // 移除 []
	shortTitle = reChineseRound.ReplaceAllString(shortTitle, "")  // 移除 （）
	shortTitle = reChineseSquare.ReplaceAllString(shortTitle, "") // 移除 【】

	// 3. 移除域名
	shortTitle = domainReg.ReplaceAllString(shortTitle, "")

	// 4. 去除开头空格
	shortTitle = reLeadingSpace.ReplaceAllString(shortTitle, "")

	// 5. 去除结尾空格
	shortTitle = reTrailingSpace.ReplaceAllString(shortTitle, "")

	// 6. 去除开头标点
	shortTitle = reLeadingPunct.ReplaceAllString(shortTitle, "")

	// 7. 再次去除首尾空格（以防上述操作后留下空格）
	shortTitle = reLeadingSpace.ReplaceAllString(shortTitle, "")
	shortTitle = reTrailingSpace.ReplaceAllString(shortTitle, "")

	// 转成 rune，便于按字符截取
	runes := []rune(shortTitle)
	originalRunes := []rune(b.Title) // 原始标题的 runes

	// 如果简化后标题过短 (<2个字符)
	if len(runes) < 2 {
		// 但原标题很长 (>15个字符)，则截取原标题前15个字符 + ...
		if len(originalRunes) > 15 {
			cutLen := 15
			// 如果原标题本身不足15，则取原标题长度
			if len(originalRunes) < cutLen {
				cutLen = len(originalRunes)
			}
			return string(originalRunes[:cutLen]) + "…"
		}
		// 如果原标题不长，或者简化后长度为0但原标题不为0，返回原标题（或原标题截断）
		if len(originalRunes) > 0 {
			cutLen := 15
			if len(originalRunes) < cutLen {
				cutLen = len(originalRunes)
				return string(originalRunes) // 如果原标题 <= 15, 直接返回原标题
			}
			return string(originalRunes[:cutLen]) + "…" // 返回截断的原标题
		}
		// 如果原标题也是空的，返回空字符串
		return ""
	}

	// [简化标题] 如果简化后长度 <= 15，直接返回
	if len(runes) <= 15 {
		return shortTitle
	}

	// [简化不完全] 超过 15 则截断加省略号
	return string(runes[:15]) + "…"
}

// GetCover 获取封面
func (b *BookInfo) GetCover() MediaFileInfo {
	switch b.Type {
	// 书籍类型为书组的时候，遍历所有子书籍，然后获取第一个子书籍的封面
	case TypeBooksGroup:
		bookGroup, err := IStore.GetBook(b.BookID)
		if err != nil {
			logger.Infof("Error getting book group: %s", err)
			return MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		for _, childID := range bookGroup.ChildBooksID {
			book, err := IStore.GetBook(childID)
			if err != nil {
				return MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
			}
			// 递归调用
			return book.GetCover()
		}
	case TypeDir, TypeZip, TypeRar, TypeCbz, TypeCbr, TypeTar, TypeEpub:
		tempBook, err := IStore.GetBook(b.BookID)
		if err != nil || len(tempBook.Images) == 0 {
			return MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		return tempBook.GuestCover()
	case TypePDF:
		return MediaFileInfo{Name: "1.jpg", Url: "/api/get_file?id=" + b.BookID + "&filename=" + "1.jpg"}
	case TypeVideo:
		return MediaFileInfo{Name: "video.png", Url: "/images/video.png"}
	case TypeAudio:
		return MediaFileInfo{Name: "audio.png", Url: "/images/audio.png"}
	case TypeUnknownFile:
		return MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
	}
	return MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
}
