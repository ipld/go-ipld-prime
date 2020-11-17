package gengraphqlserver

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
)

func genMapCommon() string {
	return `
	func resolve_map_at(p graphql.ResolveParams) (interface{}, error) {
		ts, ok := p.Source.(ipld.Node)
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
	}
	`
}

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
				Resolve: resolve_map_at,
			},
			"Keys": &graphql.Field{
				Type: graphql.NewList({{ .KeyType | LocalName}}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.(ipld.Node)
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
					ts, ok := p.Source.(ipld.Node)
					if !ok {
						return nil, errNotNode
					}
					it := ts.MapIterator()
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
			"All": &graphql.Field{
				Type: graphql.NewList({{ . | LocalName }}__entry),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.(ipld.Node)
					if !ok {
						return nil, errNotNode
					}
					it := ts.MapIterator()
					children := make([][]ipld.Node, 0)

					for !it.Done() {
						k, v, err := it.Next()
						if err != nil {
							return nil, err
						}
						children = append(children, []ipld.Node{k, v})
					}
					return children, nil
				},
			},
		},	
	})
	var {{ . | LocalName }}__entry = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{ .Name }}_Entry",
		Fields: graphql.Fields{
			"Key": &graphql.Field{
				Type: {{ .KeyType | LocalName }},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					kv, ok := p.Source.([]ipld.Node)
					if !ok {
						return nil, errNotNode
					}
					return kv[0], nil
				},
			},
			"Value": &graphql.Field{
				Type: {{ .ValueType | LocalName }},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					kv, ok := p.Source.([]ipld.Node)
					if !ok {
						return nil, errNotNode
					}
					return kv[1], nil
				},
			},
		},
	})
 	`, w, t, c)
}
