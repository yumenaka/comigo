package schema

import (
	"time"

	"entgo.io/ent"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Book 定义书籍，BooID不应该重复，根据文件路径生成
type Book struct {
	Name            string           `json:"name" storm:"index"`                              //书名 //storm:"index" 索引字段
	BookID          string           `json:"id"   storm:"id"`                                 //根据FilePath计算 //storm会搜索id或ID做为主键
	FilePath        string           `json:"-" storm:"filepath" storm:"index" storm:"unique"` //storm:"index" 索引字段 storm:"unique" 唯一字段
	BookStorePath   string           `json:"-"    storm:"index"`                              //在哪个子书库
	Type            SupportFileType  `json:"book_type" storm:"index"`                         //可以是书籍组(book_group)、文件夹(dir)、文件后缀( .zip .rar .pdf .mp4)等
	ChildBookNum    int              `json:"child_book_num" storm:"index"`                    //子书籍的数量
	ChildBook       map[string]*Book `json:"child_book" `                                     //key：BookID
	Depth           int              `json:"depth" storm:"index"`                             //文件深度
	ParentFolder    string           `json:"parent_folder" storm:"index"`                     //所在父文件夹
	AllPageNum      int              `json:"all_page_num" storm:"index"`                      //storm:"index" 索引字段
	FileSize        int64            `json:"file_size" storm:"index"`                         //storm:"index" 索引字段
	Cover           SinglePageInfo   `json:"cover" storm:"inline"`                            //storm:"inline" 内联字段，结构体嵌套时使用
	Pages           AllPageInfo      `json:"pages" storm:"inline"`                            //storm:"inline" 内联字段，结构体嵌套时使用
	Author          []string         `json:"-"`                                               //json不解析，启用可改为`json:"author"`
	ISBN            string           `json:"-"`                                               //json不解析，启用可改为`json:"isbn"`
	Press           string           `json:"-"`                                               //json不解析，启用可改为`json:"press"`        //出版社
	PublishedAt     string           `json:"-"`                                               //json不解析，启用可改为`json:"published_at"` //出版日期
	ExtractPath     string           `json:"-"`                                               //json不解析
	Modified        time.Time        `json:"-"`                                               //json不解析，启用可改为`json:"modified_time"`
	ExtractNum      int              `json:"-"`                                               //json不解析，启用可改为`json:"extract_num"`
	InitComplete    bool             `json:"-"`                                               //json不解析，启用可改为`json:"extract_complete"`
	ReadPercent     float64          `json:"-"`                                               //json不解析，启用可改为`json:"read_percent"`
	NonUTF8Zip      bool             `json:"-"`                                               //json不解析，启用可改为    `json:"non_utf8_zip"`
	ZipTextEncoding string           `json:"-"`                                               //json不解析，启用可改为   `json:"zip_text_encoding"`
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return nil
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return nil
}
