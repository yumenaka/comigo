package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Directory 模型定义（目录）
type Directory struct {
	ent.Schema
}

// Fields 定义 Directory 的字段
func (Directory) Fields() []ent.Field {
	return []ent.Field{
		field.String("path").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
	}
}

// Edges 定义 Directory 的关系（父子目录，自身关联；以及与 Image 的关系）
func (Directory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Directory.Type).From("parent").Unique(), // 子目录列表，唯一父目录
		edge.To("images", Image.Type),                               // 关联的图片列表
	}
}
