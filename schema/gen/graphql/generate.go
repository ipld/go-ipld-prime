package gengraphql

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/warpfork/go-wish"
)

func Generate(pth string, ts schema.TypeSystem) {
	withFile(filepath.Join(pth, "schema.graphql"), func(f io.Writer) {
		EmitFileHeader(f)
		for _, typ := range ts.GetTypes() {
			switch t2 := typ.(type) {
			case *schema.TypeBool:
				EmitScalar(t2, f)
			case *schema.TypeInt:
				EmitScalar(t2, f)
			case *schema.TypeFloat:
				EmitScalar(t2, f)
			case *schema.TypeString:
				EmitScalar(t2, f)
			case *schema.TypeBytes:
				EmitScalar(t2, f)
			case *schema.TypeLink:
			case *schema.TypeStruct:
				EmitStruct(t2, f)
			case *schema.TypeMap:
				EmitMap(t2, f)
			case *schema.TypeList:
				EmitList(t2, f)
			case *schema.TypeUnion:
				EmitUnion(t2, f)
			default:
				panic("unknown type" + t2.Name())
			}
		}
		EmitFileCompletion(f, ts)
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

var graphQLScalars = []string{"Int", "Float", "Boolean", "ID"}

// EmitScalar defines a scalar type for custom scalars in the type system.
func EmitScalar(t schema.Type, w io.Writer) {
	for _, bi := range graphQLScalars {
		if bi == t.Name().String() {
			return
		}
	}
	writeTemplate(`
	scalar {{ .Name }}
	`, w, t)
}

func EmitStruct(t *schema.TypeStruct, w io.Writer) {
	if len(t.Fields()) == 0 {
		writeTemplate(`
		scalar {{ .Name }}
		`, w, t)

		return
	}
	writeTemplate(`
	type {{ .Name }} {
		{{- range $field := .Fields }}
		{{ $field.Name }}: {{ $field.Type | TypeSymbol }}{{if not $field.IsMaybe}}!{{end}}
		{{- end}}
	}
	`, w, t)
}

func EmitUnion(t *schema.TypeUnion, w io.Writer) {
	writeTemplate(`
		{{ $n := .Name }}
		union {{ .Name }} = {{- range $i,$kind := .Members}}{{if $i}} | {{end}}{{if $kind | IsComplex }}{{ $kind | TypeSymbolNoRecurse }}{{else}}Wrapped_{{$n}}_{{ $kind | TypeSymbolNoRecurse }}{{end}}{{- end}}
		{{- range $kind := .Members}}{{if $kind | IsComplex }}{{else}}
		type Wrapped_{{$n}}_{{ $kind | TypeSymbolNoRecurse }} {
			{{ $kind | TypeSymbolNoRecurse }}: {{ $kind | TypeSymbolNoRecurse }}!
		}
		{{end}}{{- end}}
	`, w, t)
}

func EmitList(t *schema.TypeList, w io.Writer) {
	writeTemplate(`
	type {{ .Name }} {
		{{ .ValueType.Name }}: {{ . | TypeSymbol }}
	}
	`, w, t)
}

func EmitMap(t *schema.TypeMap, w io.Writer) {
	t.KeyType()
	writeTemplate(`
	type {{ .Name }} {
		At(key: {{ .KeyType | TypeSymbol }}): [{{ .ValueType | TypeSymbolNoRecurse }}]! 
		All: [{{.Name}}__Record!]
	}
	type {{ .Name }}__Record {
		Key: {{ .KeyType | TypeSymbol }}!
		Value: {{ .ValueType | TypeSymbol }}
	}
 	`, w, t)
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

func writeTemplate(tmpl string, w io.Writer, data interface{}) {
	f := template.FuncMap{
		"TypeSymbol":          func(t schema.Type) string { return graphQLType(t, true) },
		"TypeSymbolNoRecurse": func(t schema.Type) string { return graphQLType(t, false) },
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
	}
	t := template.Must(template.New("").Funcs(f).Parse(wish.Dedent(tmpl)))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
