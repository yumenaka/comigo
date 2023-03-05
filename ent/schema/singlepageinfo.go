package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// SinglePageInfo holds the schema definition for the SinglePageInfo entity.
type SinglePageInfo struct {
	ent.Schema
}

// Fields of the SinglePageInfo.
func (SinglePageInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("BookID"),
		field.Int("PageNum"),
		field.String("NameInArchive"),
		field.String("Url"),
		field.String("BlurHash"),
		field.Int("Height"),
		field.Int("Width"),
		field.Time("ModeTime").Default(time.Now),
		field.Int64("FileSize"),
		field.String("RealImageFilePATH"),
		field.String("ImgType"),
	}
}

// Edges of the SinglePageInfo.
func (SinglePageInfo) Edges() []ent.Edge {
	return nil
	// TODO: 如何在这里加上这个关系？ https://entgo.io/zh/docs/tutorial-todo-crud
	//return []ent.Edge{
	//	edge.From("BookID", Book.Type).
	//		Ref("Pages").
	//		Unique(),
	//}
}
