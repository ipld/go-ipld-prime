package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func EmitList(t *schema.TypeList, w io.Writer, c *config) {
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{ .Name }}",
		Fields: graphql.Fields{
			"At": &graphql.Field{
				Type: {{ .ValueType | LocalName }},
				Args: graphql.FieldConfigArgument{
					"key": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbolNoRecurse }})
					if !ok {
						return nil, errNotNode
					}

					arg := p.Args["key"]
					var out ipld.Node
					var err error
					switch ta := arg.(type) {
					case ipld.Node:
						out, err = ts.LookupByNode(ta)
					case int:
						out, err = ts.LookupByIndex(ta)
					default:
						return nil, fmt.Errorf("unknown key type: %T", arg)
					}
					{{ if (and (eq .ValueType.Kind.String "Link") (ne (.ValueType | LinkTarget) "graphql.ID")) }}
					if err != nil {
						return nil, err
					}
					targetCid, err := out.AsLink()
					if err != nil {
						return nil, err
					}
					{{ .ValueType | ReturnTarget "targetCid" }}
					{{ else }}
					return out, err
					{{ end }}
				},
			},
			"All": &graphql.Field{
				Type: graphql.NewList({{ .ValueType | LocalName }}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbolNoRecurse }})
					if !ok {
						return nil, errNotNode
					}
					it := ts.ListIterator()
					children := make([]ipld.Node, 0)
					for !it.Done() {
						_, node, err := it.Next()
						if err != nil {
							return nil, err
						}
						{{ if (and (eq .ValueType.Kind.String "Link") (ne (.ValueType | LinkTarget) "graphql.ID")) }}
						targetCid, err := node.AsLink()
						if err != nil {
							return nil, err
						}
						{{ .ValueType | IntoTarget "targetCid" "node" }}
						{{ end }}
						children = append(children, node)
					}
					return children, nil	
				},
			},
			"Range": &graphql.Field{
				Type: graphql.NewList({{ .ValueType | LocalName }}),
				Args: graphql.FieldConfigArgument{
					"skip": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"take": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbolNoRecurse }})
					if !ok {
						return nil, errNotNode
					}
					it := ts.ListIterator()
					children := make([]ipld.Node, 0)

					for !it.Done() {
						_, node, err := it.Next()
						if err != nil {
							return nil, err
						}
						{{ if (and (eq .ValueType.Kind.String "Link") (ne (.ValueType | LinkTarget) "graphql.ID")) }}
						targetCid, err := node.AsLink()
						if err != nil {
							return nil, err
						}
						{{ .ValueType | IntoTarget "targetCid" "node" }}
						{{ end }}
						children = append(children, node)
					}
					return children, nil	
				},
			},
			"Count": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbolNoRecurse }})
					if !ok {
						return nil, errNotNode
					}
					return ts.Length(), nil
				},
			},
		},
	})	
	`, w, t, c)
}
