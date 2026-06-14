package model

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// BookInfo 与 Book 唯一的区别是没有 AllPageInfo，而是封面图 URL，减小 JSON 文件的大小
type BookInfo struct {
	// ===== 基本标识 =====
	BookID string          `json:"id"`     // 根据 BookPath 生成的唯一 ID
	Title  string          `json:"title"`  // 书名
	Author string          `json:"author"` // 作者
	Type   SupportFileType `json:"type"`   // 书籍类型

	// ===== 文件路径 =====
	BookPath     string `json:"-"` // 文件绝对路径，JSON 不解析
	ParentFolder string `json:"-"` // 父文件夹名称，JSON 不解析
	StoreUrl     string `json:"-"` // 在哪个子书库, JSON 不解析

	// ===== 远程存储 =====
	IsRemote        bool   `json:"is_remote"`              // 是否为远程书籍（WebDAV、Comigo 等）
	RemoteURL       string `json:"-"`                      // 远程存储的基础 URL，JSON 不解析
	RemoteBookID    string `json:"-"`                      // 远端 Comigo 中的原始 BookID，JSON 不解析
	RemoteStoreKey  string `json:"remote_store,omitempty"` // 远端 Comigo 书库公开标识，用于前端链接参数
	RemoteShelfKey  string `json:"-"`                      // 远端 Comigo 顶级书库内部标识，JSON 不解析
	RemoteShelfName string `json:"-"`                      // 远端 Comigo 顶级书库显示名，JSON 不解析

	// ===== 文件属性 =====
	FileSize  int64     `json:"file_size"`     // 文件大小
	Modified  time.Time `json:"modified_time"` // 修改时间
	PageCount int       `json:"page_count"`    // 总页数
	Cover     PageInfo  `json:"cover"`         // 封面图

	// ===== 出版信息 =====
	ISBN        string `json:"isbn"`         // ISBN
	Press       string `json:"press"`        // 出版社
	PublishedAt string `json:"published_at"` // 出版日期

	// ===== 书组相关 =====
	ChildBooksNum int      `json:"child_books_num"` // 子书籍数量，只统计直接的子书籍
	ChildBooksID  []string `json:"child_books_id"`  // 子书籍 BookID
	Depth         int      `json:"depth"`           // 书籍深度

	// ===== 压缩包相关 =====
	ExtractPath     string `json:"extract_path"`      // 解压路径，7z 用，JSON 不解析
	ExtractNum      int    `json:"extract_num"`       // 文件解压数
	NonUTF8Zip      bool   `json:"non_utf_8_zip"`     // 是否为特殊编码 zip
	ZipTextEncoding string `json:"zip_text_encoding"` // zip 文件编码

	// ===== 状态标记 =====
	InitComplete bool `json:"init_complete"` // 是否初始化完成（todo：7z解压）
	BookComplete bool `json:"book_complete"` // 书籍是否阅读完成
	Deleted      bool `json:"deleted"`       // 源文件是否已删除

	// ===== 元数据 =====
	CreatedByVersion string `json:"created_by_version"` // 生成数据的 Comigo 版本
}

// GetAllChildBooksNum 递归获取所有子书籍的数量
func (b *BookInfo) GetAllChildBooksNum() int {
	total := 0
	for _, childID := range b.ChildBooksID {
		childBook, err := IStore.GetBook(childID)
		if err != nil {
			continue
		}
		if childBook.Type == TypeBooksGroup {
			total += childBook.GetAllChildBooksNum()
		} else {
			total++
		}
	}
	return total
}

// initBookID 根据路径和书籍属性初始化稳定 BookID。
// 行为保持原则：
// - 已存在同路径、同类型书籍时，仍返回“书籍数据已存在”的错误。
// - 完整 ID 的生成字段和短 ID 扩展顺序不变，保证重扫时 ID 前缀稳定。
// - 只把重复读取 Store 的部分换成一次快照，不改变外部可见的成功/失败结果。
func (b *BookInfo) initBookID(bookPath string) (*BookInfo, error) {
	// 只读取一次当前书库快照。旧实现会先为了重复路径检测读一次，
	// 后续每尝试一个短 ID 又重新 ListBooks；大书库下这个成本会被放大。
	allBooks, err := IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}

	// existingIDs 用于后面的短 ID 冲突判断。这里只保存 ID，不保存整本书，
	// 避免后续每扩展一位 shortID 都重新遍历 allBooks。
	existingIDs := make(map[string]struct{}, len(allBooks))

	// 保持旧逻辑的比较边界：重复书籍检测使用传入 bookPath 的 filepath.Abs 结果，
	// 而不是 b.BookPath。这样不会改变远程路径、相对路径失败等边界场景的现有结果。
	absBookPath, absErr := filepath.Abs(bookPath)
	if absErr != nil {
		logger.Infof(locale.GetString("log_error_getting_absolute_path"), absErr)
	}

	for _, existingBook := range allBooks {
		existingIDs[existingBook.BookID] = struct{}{}

		// 与旧实现一致：遍历到第一个同路径、同类型书籍时立刻返回，
		// 错误信息中也继续使用这个已存在书籍的 BookID。
		if absErr == nil && existingBook.BookPath == absBookPath && existingBook.Type == b.Type {
			return nil, fmt.Errorf(locale.GetString("log_book_data_already_exists"), existingBook.BookID, bookPath)
		}
	}

	// 生成完整 BookID 的字符串。字段顺序保持旧规则不变：
	// BookPath + FileSize + Type + ParentFolder + StoreUrl。
	idSource := b.BookPath + strconv.Itoa(int(b.FileSize)) + string(b.Type) + b.ParentFolder + b.StoreUrl
	// 两次 MD5 加密，然后转为 base62 编码。
	// 为什么选择 Base62?
	// 1. 人类可读，可以目视或简单的 regexp 进行验证。
	// 2. 仅包含字母数字符号，不包含特殊字符。
	// 3. 可以通过在任何文本编辑器和浏览器地址栏中双击鼠标来完全选择。
	// 4. 紧凑，生成的字符串比 Base32 短。
	fullID := base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(idSource))))

	// 生成短 BookID，并沿用旧策略避免冲突：
	// 从 7 位开始尝试；如果该前缀已被占用，就逐位扩展，直到没有冲突或已经扩展到 fullID 全长。
	minLength := 7
	if len(fullID) <= minLength {
		logger.Infof(locale.GetString("log_cannot_shorten_id"), fullID)
		b.BookID = fullID
		return b, nil
	}
	shortID := fullID[:minLength]
	add := 0
	for {
		if _, conflict := existingIDs[shortID]; !conflict {
			break
		}
		add++
		if minLength+add > len(fullID) {
			break
		}
		shortID = fullID[:minLength+add]
	}
	b.BookID = shortID
	return b, nil
}

// setFilePath 初始化 Book 时，设置 BookPath
func (b *BookInfo) setFilePath(path string) *BookInfo {
	// 如果 StoreUrl 是远程 URL（包含 ://），说明这是远程书籍，直接使用原始路径
	// 因为远程路径不应该使用 filepath.Abs 转换为本地绝对路径
	if strings.Contains(b.StoreUrl, "://") {
		b.BookPath = path
		return b
	}

	// 本地书籍：转换为绝对路径
	fileAbsPath, err := filepath.Abs(path)
	if err != nil {
		// 因为权限问题，无法取得绝对路径的情况下，用相对路径
		logger.Info(err, fileAbsPath)
		b.BookPath = path
	} else {
		b.BookPath = fileAbsPath
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

// ShortName 返回简短的标题（文件名）。
// 实际字符串清洗放在 tools.ShortName，避免 model 层承担无状态文本规则。
func (b *BookInfo) ShortName() string {
	return tools.ShortName(b.Title)
}

// GetCover 获取封面
func (b *BookInfo) GetCover() PageInfo {
	switch b.Type {
	// 书籍类型为书组的时候，遍历所有子书籍，然后获取第一个子书籍的封面
	case TypeBooksGroup:
		bookGroup, err := IStore.GetBook(b.BookID)
		if err != nil {
			logger.Infof(locale.GetString("log_error_getting_book_group"), err)
			return PageInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		for _, childID := range bookGroup.ChildBooksID {
			book, err := IStore.GetBook(childID)
			if err != nil {
				return PageInfo{Name: "unknown.png", Url: "/images/unknown.png"}
			}
			// 递归调用
			return book.GetCover()
		}
	case TypeDir, TypeZip, TypeRar, TypeCbz, TypeCbr, TypeTar, TypeEpub:
		tempBook, err := IStore.GetBook(b.BookID)
		if err != nil || len(tempBook.PageInfos) == 0 {
			return PageInfo{Name: "unknown.png", Url: "/images/unknown.png"}
		}
		return tempBook.GuestCover()
	case TypePDF:
		return PageInfo{Name: "1.jpg", Url: "/api/get-file?id=" + b.BookID + "&filename=" + "1.jpg"}
	case TypeVideo:
		return PageInfo{Name: "video.png", Url: "/images/video.png"}
	case TypeAudio:
		return PageInfo{Name: "audio.png", Url: "/images/audio.png"}
	case TypeHTML:
		return PageInfo{Name: "html.png", Url: "/images/html.png"}
	case TypeUnknownFile:
		return PageInfo{Name: "unknown.png", Url: "/images/unknown.png"}
	}
	return PageInfo{Name: "unknown.png", Url: "/images/unknown.png"}
}
