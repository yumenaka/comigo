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
		field.String("name").
			MaxLen(50). // 限制长度
			Unique().Comment("用户称呼"),
		field.Time("created_at").
			Default(time.Now).Comment("创建时间"),
		field.String("username").Comment("用户名").
			MaxLen(50). // 限制长度
			Unique(),   // 字段可以使用 Unique 方法定义为唯一字段。 注意：唯一字段不能有默认值。
		field.String("password").Comment("登录密码"),
		field.Time("last_login").
			Default(time.Now).Comment("最后登录时间"),
		field.Int("age").
			Positive(), // 只能取正数
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
