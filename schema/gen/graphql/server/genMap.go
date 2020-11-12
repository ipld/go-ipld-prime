package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func EmitMap(t *schema.TypeMap, w io.Writer, c *config) {
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{ .Name }}",
		Fields: graphql.Fields{
			"At": &graphql.Field{
				Type: {{ .ValueType | LocalName }},
				Args: graphql.FieldConfigArgument{
					"key": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull({{ .KeyType | LocalName }}),
					},
				},	
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbol }})
					if !ok {
						return nil, errNotNode
					}
					arg := p.Args["key"]

					switch ta := arg.(type) {
					case ipld.Node:
						return ts.LookupByNode(ta)
					case string:
						return ts.LookupByString(ta)
					default:
						return nil, fmt.Errorf("unknown key type: %T", arg)
					}
				},
			},
			"Keys": &graphql.Field{
				Type: graphql.NewList({{ .KeyType | LocalName}}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbol }})
					if !ok {
						return nil, errNotNode
					}
					it := ts.MapIterator()
					children := make([]ipld.Node, 0)

					for !it.Done() {
						node, _, err := it.Next()
						if err != nil {
							return nil, err
						}
						children = append(children, node)
					}
					return children, nil
				},
			},
			"Values": &graphql.Field{
				Type: graphql.NewList({{ .ValueType | LocalName}}),
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
		},	
	})
 	`, w, t, c)
}
