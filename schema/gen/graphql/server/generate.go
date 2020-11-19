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
	overrides      map[string]struct{}
}

func Generate(pth string, pkg string, ts schema.TypeSystem, tsPkgName, tsPkgPath string) {
	c := config{
		schemaPkg:      tsPkgName,
		initDirectives: bytes.NewBuffer(nil),
		overrides:      GetPreExistingMethods(pth),
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

func followingLinks(t schema.Type, linkName, into string, c *config) string {
	target := linkTarget(t)
	return fmt.Sprintf(`
	if cl, ok := %s.(cidlink.Link); ok {
		v := p.Context.Value(nodeLoaderCtxKey)
		if v == nil {
			return cl.Cid, nil
		}
		loader, ok := v.(func(context.Context, cidlink.Link, ipld.NodeBuilder) (ipld.Node, error))
		if !ok {
			return nil, errInvalidLoader
		}

		builder := %s.Type.%s__Repr.NewBuilder()
		n, err := loader(p.Context, cl, builder);
		if err != nil {
			return nil, err
		}
		%s = n
	} else {
		return nil, errInvalidLink
	}
	`, linkName, c.schemaPkg, target, into)
}

func linkTarget(t schema.Type) string {
	switch t2 := t.(type) {
	case *schema.TypeLink:
		if t2.HasReferencedType() {
			return t2.ReferencedType().Name().String()
		}
	}
	return "graphql.ID"
}

func writeTemplate(tmpl string, w io.Writer, data interface{}, c *config) {
	f := template.FuncMap{
		"TypePackage":         func() string { return c.schemaPkg },
		"TypeSymbol":          func(t schema.Type) string { return graphQLType(t, true) },
		"TypeRepr":            func(t schema.Type) string { return t.Kind().String() },
		"LocalName":           graphQLName,
		"TypeSymbolNoRecurse": func(t schema.Type) string { return graphQLType(t, false) },
		"IsBuiltIn":           isBuiltInScalar,
		"LinkTarget":          linkTarget,
		"ReturnTarget": func(name string, t schema.Type) string {
			return `
			var node ipld.Node
			` + followingLinks(t, name, "node", c) + `
			return node, nil
			`
		},
		"IntoTarget": func(name, into string, t schema.Type) string { return followingLinks(t, name, into, c) },
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

func writeMethod(name string, tmpl string, w io.Writer, data interface{}, c *config) {
	if _, ok := c.overrides[name]; ok {
		return
	}
	writeTemplate(tmpl, w, data, c)
}
