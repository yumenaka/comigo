package common

import (
	"sort"
	"strconv"
	"time"

	"github.com/yumenaka/comi/tools"
)

// BookInfo 与Book唯一的区别是没有AllPageInfo,而是封面图URL
type BookInfo struct {
	Name            string               `json:"name"`
	BookID          string               `json:"id"` //根据FilePath计算
	Type            BookType             `json:"book_type"`
	Depth           int                  `json:"depth"`
	ChildBookNum    int                  `json:"child_book_num"` //子书籍的数量
	ChildBook       map[string]*BookInfo `json:"child_book"`     //key：BookID
	AllPageNum      int                  `json:"all_page_num"`
	Cover           SinglePageInfo       `json:"cover"`
	ParentFolder    string               `json:"parent_folder"` //所在父文件夹
	Author          []string             `json:"-"`             //暂时用不着 这个字段不解析 `json:"author"`
	ISBN            string               `json:"-"`             //暂时用不着 这个字段不解析 `json:"isbn"`
	FilePath        string               `json:"-"`             //这个字段不解析
	ExtractPath     string               `json:"-"`             //这个字段不解析
	FileSize        int64                `json:"-"`             //暂不解析，启用可改为`json:"file_size"`
	Modified        time.Time            `json:"-"`             //暂时用不着 这个字段不解析 `json:"modified_time"`
	ExtractNum      int                  `json:"-"`             //暂时用不着 这个字段不解析 `json:"extract_num"`
	ExtractComplete bool                 `json:"-"`             //暂时用不着 这个字段不解析 `json:"extract_complete"`
	ReadPercent     float64              `json:"-"`             //暂不解析，启用可改为`json:"read_percent"`
	NonUTF8Zip      bool                 `json:"-"`             //暂时用不着 这个字段不解析 `json:"non_utf_8_zip"`
	ZipTextEncoding string               `json:"-"`             //暂时用不着 这个字段不解析 `json:"zip_text_encoding"`
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
	return &BookInfo{
		Name:            b.Name,
		Author:          b.Author,
		Depth:           b.Depth,
		ISBN:            b.ISBN,
		FilePath:        b.GetFilePath(),
		ExtractPath:     b.ExtractPath,
		AllPageNum:      b.GetAllPageNum(),
		Type:            b.Type,
		ChildBookNum:    b.ChildBookNum,
		ChildBook:       getChildInfoMap(b.ChildBook),
		FileSize:        b.FileSize,
		Modified:        b.Modified,
		BookID:          b.BookID,
		ExtractNum:      b.ExtractNum,
		ExtractComplete: b.ExtractComplete,
		ReadPercent:     b.ReadPercent,
		NonUTF8Zip:      b.NonUTF8Zip,
		Cover:           b.Cover,
		ParentFolder:    b.ParentFolder,
	}
}

// BookInfoList Slice
type BookInfoList struct {
	BookInfos []BookInfo
	SortBy    string
}

func (s BookInfoList) Len() int {
	return len(s.BookInfos)
}

// Less 按时间或URL，将图片排序
func (s BookInfoList) Less(i, j int) (less bool) {
	//如何定义 s[i] < s[j]  根据文件名(第三方库、自然语言字符串)
	if s.SortBy == "name" {
		less = tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name)
	} else if s.SortBy == "file_size" {
		less = tools.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	} else if s.SortBy == "time" {
		less = tools.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	} else if s.SortBy == "author" {
		less = tools.Compare(s.BookInfos[i].Author[0], s.BookInfos[j].Author[0])
	} else {
		less = tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name)
	}
	return less
}

func (s BookInfoList) Swap(i, j int) {
	s.BookInfos[i], s.BookInfos[j] = s.BookInfos[j], s.BookInfos[i]
}

// SortBooks 上面三个函数定义好了，终于可以使用sort包排序了
func (s *BookInfoList) SortBooks() {
	sort.Sort(s)
}
