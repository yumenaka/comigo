package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age").
			Positive(), //只能取正数
		field.String("name").
			Unique(),
		field.Time("created_at").
			Default(time.Now),
		field.String("username").
			Unique(), //字段可以使用 Unique 方法定义为唯一字段。 注意：唯一字段不能有默认值。
		field.String("password"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
