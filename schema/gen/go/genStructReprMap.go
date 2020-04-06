package gengo

import (
	"io"
	"strconv"

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
				return {{if not (MaybeUsesPtr $field.Type) }}&{{end}}n.{{ $field | FieldSymbolLower }}.v, nil
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
	// The 'idx' int is what field we'll yield next.
	// Note that this iterator doesn't mention fields that are absent.
	//  This makes things a bit trickier -- especially the 'Done' predicate,
	//   since it may have to do lookahead if there's any optionals at the end of the structure!
	//  It also means 'idx' can jump ahead by more than one per Next call in order to skip over absent fields.
	// TODO : support for implicits is still future work.

	// First: Count how many trailing fields are optional.
	//  The 'Done' predicate gets more complex when in the trailing optionals.
	fields := g.Type.Fields()
	fieldCount := len(fields)
	beginTrailingOptionalField := fieldCount
	for i := fieldCount - 1; i >= 0; i-- {
		if !fields[i].IsOptional() {
			break
		}
		beginTrailingOptionalField = i
	}
	haveTrailingOptionals := beginTrailingOptionalField < fieldCount

	// Now: finally we can get on with the actual templating.
	// FIXME: this is still yielding type-level keys -- should handle rename directives.
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}__Repr) MapIterator() ipld.MapIterator {
			{{- if .HaveTrailingOptionals }}
			end := {{ len .Type.Fields }}`+
		func() string { // this next part was too silly in templates due to lack of reverse ranging.
			v := "\n"
			for i := fieldCount - 1; i >= beginTrailingOptionalField; i-- {
				v += "\t\t\tif n." + g.AdjCfg.FieldSymbolLower(fields[i]) + ".m == schema.Maybe_Absent {\n"
				v += "\t\t\t\tend = " + strconv.Itoa(i) + "\n"
				v += "\t\t\t} else {\n"
				v += "\t\t\t\tgoto done\n"
				v += "\t\t\t}\n"
			}
			return v
		}()+`done:
			return &_{{ .Type | TypeSymbol }}__ReprMapItr{n, 0, end}
			{{- else}}
			return &_{{ .Type | TypeSymbol }}__ReprMapItr{n, 0}
			{{- end}}
		}

		type _{{ .Type | TypeSymbol }}__ReprMapItr struct {
			n   *_{{ .Type | TypeSymbol }}__Repr
			idx int
			{{if .HaveTrailingOptionals }}end int{{end}}
		}

		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Next() (k ipld.Node, v ipld.Node, _ error) {
		advance:
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
					goto advance
				}
				{{- end}}
				{{- if $field.IsNullable }}
				if itr.n.{{ $field | FieldSymbolLower }}.m == schema.Maybe_Null {
					v = ipld.Null
					break
				}
				{{- end}}
				{{- if or $field.IsOptional $field.IsNullable }}
				v = {{if not (MaybeUsesPtr $field.Type) }}&{{end}}itr.n.{{ $field | FieldSymbolLower}}.v
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
		{{- if .HaveTrailingOptionals }}
		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Done() bool {
			return itr.idx >= itr.end
		}
		{{- else}}
		func (itr *_{{ .Type | TypeSymbol }}__ReprMapItr) Done() bool {
			return itr.idx >= {{ len .Type.Fields }}
		}
		{{- end}}
	`, w, g.AdjCfg, struct {
		Type                       schema.TypeStruct
		HaveTrailingOptionals      bool
		BeginTrailingOptionalField int
	}{
		g.Type,
		haveTrailingOptionals,
		beginTrailingOptionalField,
	})
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

func (g structReprMapReprGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return structReprMapReprBuilderGenerator{
		g.AdjCfg,
		mixins.MapAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Repr",
		},
		g.PkgName,
		g.Type,
	}
}

type structReprMapReprBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapAssemblerTraits
	PkgName string
	Type    schema.TypeStruct
}

func (g structReprMapReprBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprBuilder struct {
			_{{ .Type | TypeSymbol }}__ReprAssembler
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	doTemplate(`
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
	// Cover up some of the assembler methods that are prepared to handle null;
	//  this saves them from having to check for a nil 'fcb'.
	//  (The MapBuilder.Finish method still has to check for nil 'fcb', though;
	//   it's far too much duplicate code to make two MapBuilder types for that.)
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__ReprBuilder) AssignNull() error {
			return mixins.MapAssembler{"{{ .PkgName }}.{{ .TypeName }}"}.AssignNull()
		}
		func (nb *_{{ .Type | TypeSymbol }}__ReprBuilder) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				return nb.AssignString(v2)
			}
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'state' is what it says on the tin.
	// - 's' is a bitfield for what's been **s**et.
	// - 'f' is the **f**ocused field that will be assembled next.
	// - 'z' is used to denote a null (in case we're used in a context that's acceptable).  z for **z**ilch.
	// - 'fcb' is the **f**inish **c**all**b**ack, supplied by the parent if we're a child assembler.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprAssembler struct {
			w *_{{ .Type | TypeSymbol }}
			state maState
			s int
			f int
			z bool
			fcb func() error

			{{range $field := .Type.Fields -}}
			ca_{{ $field | FieldSymbolLower }} _{{ $field.Type | TypeSymbol }}__ReprAssembler
			{{end -}}
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	// We currently disregard sizeHint.  It's not relevant to us.
	//  We could check it strictly and emit errors; presently, we don't.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) BeginMap(int) (ipld.MapAssembler, error) {
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// All assemblers have the infrastructure and memory to potentially accept null...
	//  if used within some parent structure, the handling or the rejection of null must be supplied by the 'fcb';
	//  or if used in a NodeBuilder (where null isn't valid) this method gets overriden.
	//  We don't need a nil check for 'fcb' because all parent assemblers use it, and root builders override this method.
	//  We don't pass any args to 'fcb' because we assume it comes from something that can already see this whole struct.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNull() error {
			if na.state != maState_initial {
				panic("misuse")
			}
			na.z = true
			na.state = maState_finished
			return nil
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNode(v ipld.Node) error {
			panic("todo structReprMapReprBuilderGenerator AssignNode")
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// TODO key assembler goes here.  or in a small helper method for org purposes, whatever.
	for _, field := range g.Type.Fields() {
		g.emitFieldValueAssembler(field, w)
	}
}
func (g structReprMapReprBuilderGenerator) emitFieldValueAssembler(f schema.StructField, w io.Writer) {
	// TODO for Any, this should do a whole Thing;
	// TODO for any specific type, we should be able to tersely create a new type that embeds its assembler and wraps the one method that's valid for finishing its kind.
	doTemplate(`
		// todo child assembler for field {{ .Name }}
	`, w, g.AdjCfg, f)
}
