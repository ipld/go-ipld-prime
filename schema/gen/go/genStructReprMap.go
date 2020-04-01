package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &stringReprStringGenerator{}

func NewStructReprMapGenerator(pkgName string, typ schema.TypeStruct, adjCfg *AdjunctCfg) TypeGenerator {
	return structReprMapGenerator{
		structGenerator{
			adjCfg,
			mixins.MapTraits{
				pkgName,
				string(typ.Name()),
				adjCfg.TypeSymbol(typ),
			},
			pkgName,
			typ,
		},
	}
}

type structReprMapGenerator struct {
	structGenerator
}

func (g structReprMapGenerator) GetRepresentationNodeGen() NodeGenerator {
	return structReprMapReprGenerator{
		g.AdjCfg,
		mixins.MapTraits{
			g.PkgName,
			string(g.Type.Name()) + ".Repr",
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Repr",
		},
		g.PkgName,
		g.Type,
	}
}

type structReprMapReprGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapTraits
	PkgName string
	Type    schema.TypeStruct
}

func (g structReprMapReprGenerator) EmitNodeType(w io.Writer) {
	// The type is structurally the same, but will have a different set of methods.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Repr _{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = &_{{ .Type | TypeSymbol }}__Repr{}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeMethodLookupString(w io.Writer) {
	// Similar to the type-level method, except any undef fields also return ErrNotExists.
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) LookupString(key string) (ipld.Node, error) {
			switch key {
			{{- range $field := .Type.Fields }}
			case "{{ $field | $field.Parent.RepresentationStrategy.GetFieldKey }}":
				{{- if $field.IsOptional }}
				if n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Absent {
					return ipld.Undef, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
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

func (g structReprMapReprGenerator) EmitNodeMethodLookup(w io.Writer) {
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) Lookup(key ipld.Node) (ipld.Node, error) {
			ks, err := key.AsString()
			if err != nil {
				return nil, err
			}
			return n.LookupString(ks)
		}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	// Note that this iterator doesn't mention fields that are absent.
	// FIXME significant issues here, continue doesn't work, and also continue shouldn't drive off a cliff if it's encountered at the end.
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) MapIterator() ipld.MapIterator {
			return &_{{ .Type | TypeSymbol }}__ReprMapItr{n, 0}
		}

		type _{{ .Type | TypeSymbol }}__ReprMapItr struct {
			n {{ .Type | TypeSymbol }}
			idx  int
		}

		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Next() (k ipld.Node, v ipld.Node, _ error) {
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
					itr.idx++
					continue
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
		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
		}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeMethodLength(w io.Writer) {
	// This is fun: it has to count down for any unset optional fields.
	// TODO : support for implicits is still future work.
	doTemplate(`
		func (rn *_{{ .Type | TypeSymbol }}__Repr) Length() int {
			l := {{ len .Type.Fields }}
			{{- range $field := .Type.Fields }}
			{{- if $field.IsOptional }}
			if rn.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Absent {
				l--
			}
			{{- end}}
			{{- end}}
			return l
		}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeMethodStyle(w io.Writer) {
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		func (_{{ .Type | TypeSymbol }}__Repr) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__ReprStyle{}
		}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeStyleType(w io.Writer) {
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprStyle struct{}

		func (_{{ .Type | TypeSymbol }}__ReprStyle) NewBuilder() ipld.NodeBuilder {
			var nb _{{ .Type | TypeSymbol }}__ReprBuilder
			nb.Reset()
			return &nb
		}
	`, w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g structReprMapReprGenerator) EmitNodeBuilder(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprBuilder struct {
			_{{ .Type | TypeSymbol }}__ReprAssembler
		}

		func (nb *_{{ .Type | TypeSymbol }}__ReprBuilder) Build() ipld.Node {
			if nb.state != maState_finished {
				panic("invalid state: assembler for {{ .PkgName }}.{{ .Type.Name }}.Repr must be 'finished' before Build can be called!")
			}
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__ReprBuilder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			*nb = _{{ .Type | TypeSymbol }}__ReprBuilder{_{{ .Type | TypeSymbol }}__ReprAssembler{w: &w, state: maState_initial}}
		}
	`, w, g.AdjCfg, g)
}

func (g structReprMapReprGenerator) EmitNodeAssembler(w io.Writer) {
	// FIXME this is getting egregious; it's high time to break EmitNodeAssembler down into a generator with more reusable parts.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprAssembler struct {
			w *_{{ .Type | TypeSymbol }}
			state maState
		}

		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			panic("todo structassembler reprmap beginmap")
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.BeginList(0)
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNull() error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignNull()
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignBool(bool) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignBool(false)
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignInt(int) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignInt(0)
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignFloat(float64) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignFloat(0)
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignString(v string) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignString("")
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignBytes([]byte) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignBytes(nil)
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) AssignLink(ipld.Link) error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .Type.Name }}.Repr"}.AssignLink(nil)
		}
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNode(v ipld.Node) error {
			panic("todo structassembler assignNode")
		}
		func (_{{ .Type | TypeSymbol }}__ReprAssembler) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__ReprStyle{}
		}
	`, w, g.AdjCfg, g)
}
