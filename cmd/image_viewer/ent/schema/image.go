package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Image 模型定义（图片文件）
type Image struct {
	ent.Schema
}

// Fields 定义 Image 的字段
func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.String("path").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.Int64("size"),
		field.Time("mod_time"),
		field.Time("create_time"),
	}
}

// Edges 定义 Image 与 Directory 的关系
func (Image) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("directory", Directory.Type).Ref("images").Unique().Required(),
	}
}
