package executor

const (
	schemaTemplate = `package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// {{ .Name }} holds the schema definition for the {{ .Name }} entity.
type {{ .Name }} struct {
	ent.schema
}

// Fields of the {{ .Name }}.
func ({{ .Name }}) Fields() []ent.Field {
	return []ent.Field{
	{{range .Fields }}
	    {{ . }}
    {{ end }}
	}
}

// Edges of the {{ .Name }}.
func ({{ .Name }}) Edges() []ent.Edge {
	return nil
}

// mixin of the {{ .Name }}.
func ({{ .Name }}) mixin() []ent.mixin {
	return []ent.mixin{
		TimeMixin{},
	}
}

// Indexes of the {{ .Name }}.
func ({{ .Name }}) Indexes() []ent.Index {
	return []ent.Index{
	}
}`
)
