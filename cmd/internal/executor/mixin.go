package executor

const (
	mixinTemplate = `package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Default(time.Now().Local),                     // 创建时间
		field.Time("updated_at").Default(time.Now().Local).UpdateDefault(time.Now().Local), // 更新时间
	}
}`
)
