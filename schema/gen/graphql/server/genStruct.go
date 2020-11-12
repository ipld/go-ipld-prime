package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func EmitStruct(t *schema.TypeStruct, w io.Writer, c *config) {
	if len(t.Fields()) == 0 {
		writeTemplate(`
		var {{ . | LocalName }} = graphql.NewObject(graphql.ObjectConfig{
			Name: "{{ .Name }}",
			Fields: graphql.Fields{
				"__Exists": &graphql.Field{
					Type: graphql.Boolean,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return true, nil
					},
				},
			},
		})
		`, w, t, c)

		return
	}

	for _, field := range t.Fields() {
		writeMethod(t.Name().String()+`__`+field.Name()+`__resolve`, `
		func {{ .t.Name }}__{{ .field.Name}}__resolve(p graphql.ResolveParams) (interface{}, error) {
			ts, ok := p.Source.({{ TypePackage }}.{{ .t | TypeSymbol }})
			if !ok {
				return nil, errNotNode
			}
			{{ if .field.IsMaybe }}
			f := ts.Field{{ .field.Name }}()
			if f.Exists() {
				{{ if .field.Type | IsBuiltIn }}
				return f.Must().As{{ .field.Type | TypeRepr }}()
				{{ else if (eq .field.Type.Kind.String "Link") }}
				return "IS a link", nil
				{{ else }}
				return f.Must(), nil
				{{ end }}
			} else {
				return nil, nil
			}
			{{ else if .field.Type | IsBuiltIn }}
			return ts.Field{{ .field.Name }}().As{{ .field.Type | TypeRepr }}()
			{{ else if (and (eq .field.Type.Kind.String "Link") (ne (.field.Type | LinkTarget) "graphql.ID")) }}
			targetCid := ts.Field{{ .field.Name }}().Link()
			{{ .field.Type | ReturnTarget "targetCid" }}
			{{ else }}
			return ts.Field{{ .field.Name }}(), nil
			{{ end }}
		}
		`, w, map[string]interface{}{
			"t":     t,
			"field": field,
		}, c)
	}

	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{ .Name }}",
		Fields: graphql.Fields{
			{{- range $field := .Fields }}
			"{{ .Name }}": &graphql.Field{
				{{ if .IsMaybe }}
				Type: {{ .Type | LocalName }},
				{{ else }}
				Type: graphql.NewNonNull({{ .Type | LocalName }}),
				{{ end }}
				Resolve: {{ $.Name }}__{{ $field.Name }}__resolve,
			},			
			{{- end}}
		},
	})
	`, w, t, c)
}
