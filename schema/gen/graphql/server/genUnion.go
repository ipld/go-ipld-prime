package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func EmitUnion(t *schema.TypeUnion, w io.Writer, c *config) {
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewUnion(graphql.UnionConfig{
		Name: "{{ .Name }}",
		Types: []*graphql.Object{
			{{- range $kind := .Members}}
			{{if $kind | IsStruct}}
			{{ $kind | LocalName }},
			{{else}}
			union__{{$.Name}}__{{$kind.Name}},
			{{end}}
			{{- end}}
		},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if node, ok := p.Value.(ipld.Node); ok {
				switch node.Prototype() {
					{{- range $kind := .Members}}
				{{if $kind | IsStruct}}
				case {{ TypePackage}}.Type.{{ $kind.Name }}:
					fallthrough
				case {{ TypePackage}}.Type.{{ $kind.Name }}__Repr:
					return {{ $kind | LocalName }}
				{{else}}
				case {{ TypePackage}}.Type.{{ $kind.Name }}:
					fallthrough
				case {{ TypePackage}}.Type.{{ $kind.Name }}__Repr:
					return union__{{$.Name}}__{{$kind.Name}}
				{{end}}
				{{- end}}
				}				
			}
			fmt.Printf("Actual type %T: %v not in union\n", p.Value, p.Value)
			return nil
		},
	})

	{{- range $kind := .Members}}
	{{if $kind | IsStruct}}
	{{else}}
	var union__{{$.Name}}__{{$kind.Name}} = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{$.Name}}.{{$kind.Name}}",
		Description: "Synthetic union member wrapper",
		Fields: graphql.Fields{
			{{ if $kind | IsBuiltIn }}
			"": &graphql.Field{
				Type: {{ $kind | LocalName}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ $kind | TypeSymbolNoRecurse }})
					if !ok {
						return nil, errNotNode
					}
					return ts.As{{ $kind | TypeRepr }}()
				},
			},
			{{ end }}
		},
	})
	{{end}}
	{{- end}}
	`, w, t, c)

	// types which may involve type reference cycles defer to init block
	// to make golang compiler happy.
	writeTemplate(`
	{{- range $kind := .Members}}
	{{if (eq "Map" $kind.Kind.String) }}
	union__{{$.Name}}__{{$kind.Name}}.AddFieldConfig("", &graphql.Field{
		Type: {{ $kind | LocalName}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ts, ok := p.Source.({{ TypePackage }}.{{ $kind | TypeSymbolNoRecurse }})
			if !ok {
				return nil, errNotNode
			}
			mi := ts.MapIterator()
			items := make(map[string]interface{})
			for !mi.Done() {
				k, v, err := mi.Next()
				if err != nil {
					return nil, err
				}
				// TODO: key type may not be string.
				ks, err := k.AsString()
				if err != nil {
					return nil, err
				}
				items[ks] = v
			}
			return items, nil
		},
	})
	{{else if (eq "List" $kind.Kind.String) }}
	union__{{$.Name}}__{{$kind.Name}}.AddFieldConfig("", &graphql.Field{
		Type: {{ $kind | LocalName}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ts, ok := p.Source.({{ TypePackage }}.{{ $kind | TypeSymbolNoRecurse }})
			if !ok {
				return nil, errNotNode
			}
			li := ts.ListIterator()
			items := make([]ipld.Node, 0)
			for !li.Done() {
				_, v, err := li.Next()
				if err != nil {
					return nil, err
				}
				items = append(items, v)
			}
			return items, nil
		},
	})
	{{end}}
	{{- end}}
	`, c.initDirectives, t, c)
}
