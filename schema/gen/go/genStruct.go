package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type structGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapTraits
	PkgName string
	Type    schema.TypeStruct
}

// --- native content and specializations --->

func (g structGenerator) EmitNativeType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			{{- range $field := .Type.Fields}}
			{{ $field | FieldSymbolLower }} _{{ $field.Type | TypeSymbol }}{{if $field.IsMaybe }}__Maybe{{end}}
			{{- end}}
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNativeAccessors(w io.Writer) {
	doTemplate(`
		{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
		{{- range $field := .Type.Fields }}
		func (n _{{ $type | TypeSymbol }}) Field{{ $field | FieldSymbolUpper }}()	{{ if $field.IsMaybe }}Maybe{{end}}{{ $field.Type | TypeSymbol }} {
			return &n.{{ $field | FieldSymbolLower }}
		}
		{{- end}}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNativeBuilder(w io.Writer) {
	// Unclear what, if anything, goes here.
}

func (g structGenerator) EmitNativeMaybe(w io.Writer) {
	// TODO maybes need a lot of questions answered
}

// --- type info --->

func (g structGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g structGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	// Perhaps surprisingly, the way to get the representation node pointer
	//  does not actually depend on what the representation strategy is.
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Representation() ipld.Node {
			return (*_{{ .Type | TypeSymbol }}__Repr)(n)
		}
	`, w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g structGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
	// We do, however, want some constants for our fields;
	//  they'll make iterators able to work faster.  So let's emit those.
	doTemplate(`
		var (
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} = _String{"{{ $field.Name }}"}
			{{- end }}
		)
	`, w, g.AdjCfg, g)

}

func (g structGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
		var _ schema.TypedNode = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupString(key string) (ipld.Node, error) {
			switch key {
			{{- range $field := .Type.Fields }}
			case "{{ $field.Name }}":
				{{- if $field.IsOptional }}
				if n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Absent {
					return ipld.Undef, nil
				}
				{{- end}}
				{{- if $field.IsNullable }}
				if n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Null {
					return ipld.Null, nil
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				return n.{{ $field | FieldSymbolLower }}.v, nil
				{{- else}}
				return &n.{{ $field | FieldSymbolLower }}, nil
				{{- end}}
			{{- end}}
			default:
				return nil, schema.ErrNoSuchField{Type: nil /*TODO*/, FieldName: key}
			}
		}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Lookup(key ipld.Node) (ipld.Node, error) {
			ks, err := key.AsString()
			if err != nil {
				return nil, err
			}
			return n.LookupString(ks)
		}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	// Note that the typed iterator will report absent fields.
	//  The representation iterator (if has one) however will skip those.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) MapIterator() ipld.MapIterator {
			return &_{{ .Type | TypeSymbol }}__MapItr{n, 0}
		}

		type _{{ .Type | TypeSymbol }}__MapItr struct {
			n {{ .Type | TypeSymbol }}
			idx  int
		}

		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Next() (k ipld.Node, v ipld.Node, _ error) {
			if itr.idx >= {{ len .Type.Fields }} {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			switch itr.idx {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				k = &fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}
				{{- if $field.IsOptional }}
				if itr.n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Absent {
					v = ipld.Undef
					break
				}
				{{- end}}
				{{- if $field.IsNullable }}
				if itr.n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Null {
					v = ipld.Null
					break
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				v = itr.n.{{ $field | FieldSymbolLower}}.v
				{{- else}}
				v = &itr.n.{{ $field | FieldSymbolLower}}
				{{- end}}
			{{- end}}
			default:
				panic("unreachable")
			}
			itr.idx++
			return
		}
		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
		}

	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Length() int {
			return {{ len .Type.Fields }}
		}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeMethodStyle(w io.Writer) {
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}

func (g structGenerator) EmitNodeStyleType(w io.Writer) {
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Style struct{}

		func (_{{ .Type | TypeSymbol }}__Style) NewBuilder() ipld.NodeBuilder {
			var nb _{{ .Type | TypeSymbol }}__Builder
			nb.Reset()
			return &nb
		}
	`, w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g structGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return structBuilderGenerator{
		g.AdjCfg,
		mixins.MapAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type structBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapAssemblerTraits
	PkgName string
	Type    schema.TypeStruct
}

func (g structBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			if nb.state != maState_finished {
				panic("invalid state: assembler for {{ .PkgName }}.{{ .Type.Name }} must be 'finished' before Build can be called!")
			}
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{w: &w, state: maState_initial}}
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			state maState
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			panic("todo structBuilderGenerator BeginMap")
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			panic("todo structBuilderGenerator AssignNode")
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// TODO key assembler goes here.  or in a small helper method for org purposes, whatever.
	for _, field := range g.Type.Fields() {
		g.emitFieldValueAssembler(field, w)
	}
}
func (g structBuilderGenerator) emitFieldValueAssembler(f schema.StructField, w io.Writer) {
	// TODO for Any, this should do a whole Thing;
	// TODO for any specific type, we should be able to tersely create a new type that embeds its assembler and wraps the one method that's valid for finishing its kind.
	doTemplate(`
		// todo child assembler for field {{ .Name }}
	`, w, g.AdjCfg, f)
}
