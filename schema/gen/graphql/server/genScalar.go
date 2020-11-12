package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

var graphQLScalars = []string{"Int", "Float", "Boolean", "String", "ID"}

func isBuiltInScalar(t schema.Type) bool {
	for _, bi := range graphQLScalars {
		if bi == t.Name().String() {
			return true
		}
	}
	if t.Name().String() == "Bool" {
		return true
	}
	return false
}

// EmitScalar defines a scalar type for custom scalars in the type system.
func EmitScalar(t schema.Type, w io.Writer, c *config) {
	if isBuiltInScalar(t) {
		return
	}
	writeMethod(graphQLName(t)+`__serialize`, `
	func {{ . | LocalName}}__serialize(value interface{}) interface{} {
		switch value := value.(type) {
		case ipld.Node:
			{{ if (eq .Kind.String "Int") }}
			i, err := value.AsInt()
			if err != nil {
				return err
			}
			return i
			{{ else if (eq .Kind.String "Bytes") }}
			b, err := value.AsBytes()
			if err != nil {
				return err
			}
			return b
			{{ else if (eq .Kind.String "String") }}
			s, err := value.AsString()
			if err != nil {
				return err
			}
			return s
			{{ else }}
			return value.As{{ .Kind }}()
			{{ end }}
		default:
			return nil
		}
	}
	`, w, t, c)

	writeMethod(graphQLName(t)+`__parse`, `
	func {{ . | LocalName}}__parse(value interface{}) interface{} {
		builder := {{ TypePackage}}.Type.{{ .Name }}__Repr.NewBuilder()
		switch v2 := value.(type) {
		case string:
			builder.AssignString(v2)
		case *string:
			builder.AssignString(*v2)
		default:
			return nil
		}
		return builder.Build()
	}
	`, w, t, c)

	writeMethod(graphQLName(t)+`__parseLiteral`, `
	func {{. |LocalName}}__parseLiteral(valueAST ast.Value) interface{} {
		builder := {{ TypePackage}}.Type.{{ .Name }}__Repr.NewBuilder()
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			builder.AssignString(valueAST.Value)
		default:
			return nil
		}
		return builder.Build()
	}
	`, w, t, c)

	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "{{ .Name }}",
		Description: "{{ .Name }}",
		Serialize: {{ . | LocalName}}__serialize,
		ParseValue: {{ . | LocalName}}__parse,
		ParseLiteral: {{. | LocalName}}__parseLiteral,
	})
	`, w, t, c)
}
