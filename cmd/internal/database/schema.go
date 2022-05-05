package database

const (
	SchemaTemplate = `package schema

import (
	{{range .Imports }}
	{{ . }}
    {{ end }}
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	{{range .EntImports }}
	{{ . }}
    {{ end }}
)

// {{ .Name }} holds the schema definition for the {{ .Name }} entity.
type {{ .Name }} struct {
	ent.Schema
}

// Annotations of the {{ .Name }}.
func ({{ .Name }}) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "{{ .TableName }}" },
    }
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

// Mixin of the {{ .Name }}.
func ({{ .Name }}) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Indexes of the {{ .Name }}.
func ({{ .Name }}) Indexes() []ent.Index {
	return []ent.Index{
	{{range .Indexs }}
	    {{ . }}
    {{ end }}
	}
}`
)
