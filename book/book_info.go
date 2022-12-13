package book

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
	Type            SupportFileType      `json:"book_type"`
	Depth           int                  `json:"depth"`
	ChildBookNum    int                  `json:"child_book_num"` //子书籍的数量
	ChildBook       map[string]*BookInfo `json:"child_book"`     //key：BookID
	AllPageNum      int                  `json:"all_page_num"`
	Cover           ImageInfo            `json:"cover"`
	ParentFolder    string               `json:"parent_folder"` //所在父文件夹
	Author          []string             `json:"-"`             //暂时用不着 这个字段不解析 `json:"author"`
	ISBN            string               `json:"-"`             //暂时用不着 这个字段不解析 `json:"isbn"`
	FilePath        string               `json:"-"`             //这个字段不解析
	ExtractPath     string               `json:"-"`             //这个字段不解析
	FileSize        int64                `json:"-"`             //暂不解析，启用可改为`json:"file_size"`
	Modified        time.Time            `json:"-"`             //暂时用不着 这个字段不解析 `json:"modified_time"`
	ExtractNum      int                  `json:"-"`             //暂时用不着 这个字段不解析 `json:"extract_num"`
	InitComplete    bool                 `json:"-"`             //暂时用不着 这个字段不解析 `json:"extract_complete"`
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
		ChildBook:    getChildInfoMap(b.ChildBook),
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
		//如果一样大（可能是大小默认为0的文件夹），先比较子书籍数量
		if s.BookInfos[i].FileSize == s.BookInfos[j].FileSize && s.BookInfos[i].ChildBookNum != s.BookInfos[j].ChildBookNum {
			return !(s.BookInfos[i].ChildBookNum < s.BookInfos[j].ChildBookNum) //
		}
		//子书籍数量也一样，就比较书名（免得结果都不一样）
		if s.BookInfos[i].FileSize == s.BookInfos[j].FileSize && s.BookInfos[i].ChildBookNum == s.BookInfos[j].ChildBookNum {
			return !tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name) //
		}
		//一般情况，比较文件大小
		return !tools.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time":
		return !tools.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author":
		return tools.Compare(s.BookInfos[i].Author[0], s.BookInfos[j].Author[0])
	//如何定义 Images[i] < Images[j] 反向
	case "filename_reverse":
		return !tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name) //(使用了第三方库、比较自然语言字符串)
	case "filesize_reverse":
		//如果一样大（很可能是大小默认为0的文件夹），先比较子书籍数量
		if s.BookInfos[i].FileSize == s.BookInfos[j].FileSize && s.BookInfos[i].ChildBookNum != s.BookInfos[j].ChildBookNum {
			return s.BookInfos[i].ChildBookNum < s.BookInfos[j].ChildBookNum //
		}
		//子书籍数量也一样的话，就比较书名（免得结果都不一样）
		if s.BookInfos[i].FileSize == s.BookInfos[j].FileSize && s.BookInfos[i].ChildBookNum == s.BookInfos[j].ChildBookNum {
			return tools.Compare(s.BookInfos[i].Name, s.BookInfos[j].Name) //
		}
		//一般情况，比较文件大小
		return tools.Compare(strconv.Itoa(int(s.BookInfos[i].FileSize)), strconv.Itoa(int(s.BookInfos[j].FileSize)))
	case "modify_time_reverse":
		return tools.Compare(s.BookInfos[i].Modified.String(), s.BookInfos[j].Modified.String())
	case "author_reverse":
		return !tools.Compare(s.BookInfos[i].Author[0], s.BookInfos[j].Author[0])
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
