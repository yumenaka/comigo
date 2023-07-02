package book

import (
	"sort"
	"strconv"
	"time"

	"github.com/yumenaka/comi/tools"
)

// BookInfo 与Book唯一的区别是没有AllPageInfo,而是封面图URL 减小 json文件的大小
type BookInfo struct {
	Name            string          `json:"name"`           //书名
	BookID          string          `json:"id"`             //根据FilePath生成的唯一ID
	BookStorePath   string          `json:"-"   `           //在哪个子书库
	Type            SupportFileType `json:"book_type"`      //书籍类型
	Depth           int             `json:"depth"`          //书籍深度
	ChildBookNum    int             `json:"child_book_num"` //子书籍数量
	AllPageNum      int             `json:"all_page_num"`   //总页数
	Cover           ImageInfo       `json:"cover"`          //封面图
	ParentFolder    string          `json:"parent_folder"`  //父文件夹
	Author          string          `json:"author"`         //作者
	ISBN            string          `json:"isbn"`           //ISBN
	Press           string          `json:"press"`          //出版社
	PublishedAt     string          `json:"published_at"`   //出版日期
	FilePath        string          `json:"-"`              //文件绝对路径，json不解析
	FileSize        int64           `json:"file_size"`      //文件大小
	Modified        time.Time       `json:"modified_time"`  //修改时间
	ExtractPath     string          `json:"-"`              //解压路径，7z用，json不解析
	ExtractNum      int             `json:"-"`              //文件解压数    extract_num
	InitComplete    bool            `json:"-"`              //是否解压完成  extract_complete
	ReadPercent     float64         `json:"read_percent"`   //阅读进度
	NonUTF8Zip      bool            `json:"-"`              //是否为特殊编码zip
	ZipTextEncoding string          `json:"-"`              //zip文件编码
}

func getChildInfoMap(ChildBookMap map[string]*Book) (ChildInfoMap map[string]*BookInfo) {
	ChildInfoMap = make(map[string]*BookInfo)
	for key, book := range ChildBookMap {
		ChildInfoMap[key] = NewBookInfo(book)
	}
	return ChildInfoMap
}

// NewBookInfo BookInfo的模拟构造函数
func NewBookInfo(b *Book) *BookInfo {
	//需要单独先执行这个，来设定封面
	allPageNum := b.GetAllPageNum()
	return &BookInfo{
		Name:         b.Name,
		Author:       b.Author,
		Depth:        b.Depth,
		ISBN:         b.ISBN,
		FilePath:     b.GetFilePath(),
		ExtractPath:  b.ExtractPath,
		AllPageNum:   allPageNum,
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

// BookInfoList Slice
type BookInfoList struct {
	BookInfos []BookInfo
	sortBy    string
}

func (s BookInfoList) Len() int {
	return len(s.BookInfos)
}

// Less 按时间或URL，将图片排序
func (s BookInfoList) Less(i, j int) (less bool) {
	//如何定义 Images[i] < Images[j]
	//根据文件名
	switch s.sortBy {
	case "filename":
		return tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name) //(使用了第三方库、比较自然语言字符串)
	case "filesize":
		//两本之中有一本是书籍组。同样是书籍组，比较子书籍数。
		if s.BookInfos[i].Type == TypeBooksGroup || s.BookInfos[j].Type == TypeBooksGroup {
			if s.BookInfos[i].Type == TypeBooksGroup && s.BookInfos[j].Type == TypeBooksGroup {
				return s.BookInfos[i].ChildBookNum > s.BookInfos[j].ChildBookNum
			}
			if s.BookInfos[i].Type != TypeBooksGroup || s.BookInfos[j].Type != TypeBooksGroup {
				return s.BookInfos[i].Type == TypeBooksGroup
			}
		}
		//两本之中有一本是文件夹。同样是文件夹，比较页数。
		if s.BookInfos[i].Type == TypeDir || s.BookInfos[j].Type == TypeDir {
			if s.BookInfos[i].Type == TypeDir && s.BookInfos[j].Type == TypeDir {
				return s.BookInfos[i].AllPageNum > s.BookInfos[j].AllPageNum
			}
			if s.BookInfos[i].Type != TypeDir || s.BookInfos[j].Type != TypeDir {
				return s.BookInfos[i].Type == TypeDir
			}
		}
		//一般情况，比较文件大小
		return !tools.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time":
		return !tools.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author":
		return tools.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
	//如何定义 Images[i] < Images[j] 反向
	case "filename_reverse":
		return !tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name) //(使用了第三方库、比较自然语言字符串)
	case "filesize_reverse":
		//两本之中有一本是书籍组。同样是书籍组，比较子书籍数。
		if s.BookInfos[i].Type == TypeBooksGroup || s.BookInfos[j].Type == TypeBooksGroup {
			if s.BookInfos[i].Type == TypeBooksGroup && s.BookInfos[j].Type == TypeBooksGroup {
				return !(s.BookInfos[i].ChildBookNum > s.BookInfos[j].ChildBookNum)
			}
			if s.BookInfos[i].Type != TypeBooksGroup || s.BookInfos[j].Type != TypeBooksGroup {
				return !(s.BookInfos[i].Type == TypeBooksGroup)
			}
		}
		//两本之中有一本是文件夹。同样是文件夹，比较页数。
		if s.BookInfos[i].Type == TypeDir || s.BookInfos[j].Type == TypeDir {
			if s.BookInfos[i].Type == TypeDir && s.BookInfos[j].Type == TypeDir {
				return !(s.BookInfos[i].AllPageNum > s.BookInfos[j].AllPageNum)
			}
			if s.BookInfos[i].Type != TypeDir || s.BookInfos[j].Type != TypeDir {
				return !(s.BookInfos[i].Type == TypeDir)
			}
		}
		//一般情况，比较文件大小
		return tools.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time_reverse":
		return tools.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author_reverse":
		return !tools.Compare(s.BookInfos[i].Author, s.BookInfos[j].Author)
	default:
		return tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name)
	}
}

func (s BookInfoList) Swap(i, j int) {
	s.BookInfos[i], s.BookInfos[j] = s.BookInfos[j], s.BookInfos[i]
}

// SortBooks 上面三个函数定义好了，终于可以使用sort包排序了
func (s *BookInfoList) SortBooks(by string) {
	if by != "" {
		s.sortBy = by
		sort.Sort(s)
	}
}
