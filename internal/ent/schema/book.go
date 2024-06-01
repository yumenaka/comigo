package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Book 定义书籍，BookID不应该重复，根据文件路径生成
type Book struct {
	ent.Schema
}

// Fields 每次添加或修改 fields 和 edges后, 都需要在项目的根目录执行 go generate ./ent 命令重新生成文件
// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.String("Title").
			MaxLen(1024). //限制长度
			Comment("书名"),
		field.String("BookID").
			Unique().Comment("书籍ID"), //字段可以使用 Unique 方法定义为唯一字段。 注意：唯一字段不能有默认值。
		field.String("Owner").
			Default("admin").
			Comment("拥有者"),
		field.String("FilePath").Comment("文件路径"),
		field.String("BookStorePath").Comment("书库路径"),
		field.String("Type").Comment("书籍类型"),
		field.Int("ChildBookNum").NonNegative(),
		field.Int("Depth").NonNegative(),
		field.String("ParentFolder"),
		field.Int("PageCount").
			NonNegative(). //内置校验器，非负数
			Comment("总页数"),
		field.Int64("FileSize"),
		field.String("Authors"),
		field.String("ISBN"),
		field.String("Press"),
		field.String("PublishedAt"),
		field.String("ExtractPath"),
		field.Time("Modified").
			Default(time.Now). //设置默认值
			Comment("创建时间"),
		field.Int("ExtractNum"),
		field.Bool("InitComplete"),
		field.Float("ReadPercent"),
		field.Bool("NonUTF8Zip"),
		field.String("ZipTextEncoding"),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("PageInfos", SinglePageInfo.Type), //Type是一种虚拟方法，用于Edge(关系)声明。
	}
}
