package gengraphqlserver

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/warpfork/go-wish"
)

type config struct {
	schemaPkg      string
	initDirectives *bytes.Buffer
}

func Generate(pth string, pkg string, ts schema.TypeSystem, tsPkgName, tsPkgPath string) {
	c := config{
		schemaPkg:      tsPkgName,
		initDirectives: bytes.NewBuffer(nil),
	}
	withFile(filepath.Join(pth, "schema.go"), func(f io.Writer) {
		EmitFileHeader(f, pkg, tsPkgPath, &c)
		for _, typ := range ts.GetTypes() {
			switch t2 := typ.(type) {
			case *schema.TypeBool:
				EmitScalar(t2, f, &c)
			case *schema.TypeInt:
				EmitScalar(t2, f, &c)
			case *schema.TypeFloat:
				EmitScalar(t2, f, &c)
			case *schema.TypeString:
				EmitScalar(t2, f, &c)
			case *schema.TypeBytes:
				EmitScalar(t2, f, &c)
			case *schema.TypeLink:
			case *schema.TypeStruct:
				EmitStruct(t2, f, &c)
			case *schema.TypeMap:
				EmitMap(t2, f, &c)
			case *schema.TypeList:
				EmitList(t2, f, &c)
			case *schema.TypeUnion:
				EmitUnion(t2, f, &c)
			default:
				panic("unknown type" + t2.Name())
			}
		}
		EmitFileCompletion(f, ts, &c)
	})
}

func withFile(filename string, fn func(io.Writer)) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fn(f)
}

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
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewScalar(graphql.ScalarConfig{
		Name:        "{{ .Name }}",
		Description: "{{ .Name }}",
		Serialize: func(value interface{}) interface{} {
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
		},
		ParseValue: func(value interface{}) interface{} {
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
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			builder := {{ TypePackage}}.Type.{{ .Name }}__Repr.NewBuilder()
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				builder.AssignString(valueAST.Value)
			default:
				return nil
			}
			return builder.Build()
		},
	})
	`, w, t, c)
}

func EmitStruct(t *schema.TypeStruct, w io.Writer, c *config) {
	if len(t.Fields()) == 0 {
		writeTemplate(`
		var {{ . | LocalName }} = graphql.NewScalar(graphql.ScalarConfig{
			Name: "{{ .Name }}",
		})
		`, w, t, c)

		return
	}
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewObject(graphql.ObjectConfig{
		Name: "{{ .Name }}",
		Fields: graphql.Fields{
			{{- range $field := .Fields }}
			"{{ .Name }}": &graphql.Field{
				Type: {{ .Type | LocalName }},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ $ | TypeSymbol }})
					if !ok {
						return nil, fmt.Errorf("not node")
					}
					{{ if $field.IsMaybe }}
					f := ts.Field{{ $field.Name }}()
					if f.Exists() {
						{{ if $field.Type | IsBuiltIn }}
						return f.Must().As{{ $field.Type | TypeRepr }}()
						{{ else if (eq $field.Type.Kind.String "Link") }}
						return "IS a link", nil
						{{ else }}
						return f.Must(), nil
						{{ end }}
					} else {
						return nil, nil
					}
					{{ else if $field.Type | IsBuiltIn }}
					return ts.Field{{ $field.Name }}().As{{ $field.Type | TypeRepr }}()
					{{ else if (and (eq $field.Type.Kind.String "Link") (ne ($field.Type | LinkTarget) "graphql.ID")) }}
					targetCid := ts.Field{{ $field.Name }}().Link()
					if cl, ok := targetCid.(cidlink.Link); ok {
						v := p.Context.Value(nodeLoaderCtxKey)
						if v == nil {
							return cl.Cid, nil
						}
						loader, ok := v.(func(context.Context, cidlink.Link, ipld.NodeBuilder) error)
						if !ok {
							return nil, fmt.Errorf("invalid Loader provided")
						}

						builder := {{ TypePackage }}.Type.{{ $field.Type | LinkTarget }}__Repr.NewBuilder()
						if err := loader(p.Context, cl, builder); err != nil {
							return nil, err
						}
						return builder.Build(), nil
					}
					return nil, fmt.Errorf("Invalid link")
					{{ else }}
					return ts.Field{{ $field.Name }}(), nil
					{{ end }}
				},
			},			
			{{- end}}
		},
	})
	`, w, t, c)
}

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
				case {{ TypePackage}}.Type.{{ $kind.Name }}__Repr:
					return {{ $kind | LocalName }}
					{{- end}}
				{{end}}
				}				
			}
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
						return nil, fmt.Errorf("not node")
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
				return nil, fmt.Errorf("not node")
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
				return nil, fmt.Errorf("not node")
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

func EmitList(t *schema.TypeList, w io.Writer, c *config) {
	writeTemplate(`
	var {{ . | LocalName }} = graphql.NewList({{ .ValueType | LocalName }})
	`, w, t, c)
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
						Type: {{ .KeyType | LocalName }},
					},
				},	
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ts, ok := p.Source.({{ TypePackage }}.{{ . | TypeSymbol }})
					if !ok {
						return nil, fmt.Errorf("unexpected node type")
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
		},	
	})
 	`, w, t, c)
}

func graphQLType(t schema.Type, allowRecurse bool) string {
	switch t2 := t.(type) {
	case *schema.TypeList:
		if !allowRecurse {
			// multiple layers of lists can't be directly expressed.
			return string(t.Name())
		}
		if t2.ValueIsNullable() {
			return "[" + graphQLType(t2.ValueType(), false) + "]"
		}
		return "[" + graphQLType(t2.ValueType(), false) + "!]"
	case *schema.TypeLink:
		if t2.HasReferencedType() {
			return graphQLType(t2.ReferencedType(), allowRecurse)
		}
		return "ID"
	}

	return string(t.Name())
}

func graphQLName(t schema.Type) string {
	n := t.Name().String()
	if isBuiltInScalar(t) {
		if n == "Bool" {
			n = "Boolean"
		}
		return "graphql." + n
	}
	if t.Kind() == schema.Kind_Link {
		sl, ok := t.(*schema.TypeLink)
		if !ok {
			return fmt.Sprintf("t is link but err: %v\n", t)
		}
		if sl.HasReferencedType() {
			return graphQLName(sl.ReferencedType())
		}
		return "graphql.ID"
	}
	return n + "__type"
}

func writeTemplate(tmpl string, w io.Writer, data interface{}, c *config) {
	f := template.FuncMap{
		"TypePackage":         func() string { return c.schemaPkg },
		"TypeSymbol":          func(t schema.Type) string { return graphQLType(t, true) },
		"TypeRepr":            func(t schema.Type) string { return t.Kind().String() },
		"LocalName":           graphQLName,
		"TypeSymbolNoRecurse": func(t schema.Type) string { return graphQLType(t, false) },
		"IsBuiltIn":           isBuiltInScalar,
		"LinkTarget": func(t schema.Type) string {
			switch t2 := t.(type) {
			case *schema.TypeLink:
				if t2.HasReferencedType() {
					return t2.ReferencedType().Name().String()
				}
			}
			return "graphql.ID"
		},
		"IsComplex": func(t schema.Type) bool {
			switch t2 := t.(type) {
			case *schema.TypeStruct:
				return true
			case *schema.TypeUnion:
				return true
			case *schema.TypeMap:
				return true
			case *schema.TypeList:
				return !t2.IsAnonymous()
			}
			return false
		},
		"IsStruct": func(t schema.Type) bool {
			switch t.(type) {
			case *schema.TypeStruct:
				return true
			}
			return false
		},
	}
	t := template.Must(template.New("").Funcs(f).Parse(wish.Dedent(tmpl)))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
