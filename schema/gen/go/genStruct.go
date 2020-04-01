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
	// - 'w' is the "**w**ip" pointer.
	// - 'state' is what it says on the tin.
	// - 's' is a bitfield for what's been **s**et.
	// - 'f' is the **f**ocused field that will be assembled next.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			state maState
			s int
			f int
		}

		var (
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} = 1 << {{ $i }}
			{{- end}}
			fieldBits__{{ $type | TypeSymbol }}_sufficient = 0 {{- range $i, $field := .Type.Fields }}{{if not $field.IsOptional }} + 1 << {{ $i }}{{end}}{{end}}
		)
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	// We currently disregard sizeHint.  It's not relevant to us.
	//  We could check it strictly and emit errors; presently, we don't.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			if na.state != maState_initial {
				panic("misuse")
			}
			if v2, ok := v.(*_{{ .Type | TypeSymbol }}); ok {
				*na.w = *v2
				na.state = maState_finished
				return nil
			}
			if v.ReprKind() != ipld.ReprKind_Map {
				return ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .Type.Name }}", MethodName: "AssignNode", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: v.ReprKind()}
			}
			itr := v.MapIterator()
			for !itr.Done() {
				k, v, err := itr.Next()
				if err != nil {
					return err
				}
				if err := na.AssembleKey().AssignNode(k); err != nil {
					return err
				}
				if err := na.AssembleValue().AssignNode(v); err != nil {
					return err
				}
			}
			return na.Finish()
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	g.emitMapAssemblerMethods(w)
	g.emitKeyAssembler(w)
	for _, field := range g.Type.Fields() {
		g.emitFieldValueAssembler(field, w)
	}
}
func (g structBuilderGenerator) emitMapAssemblerMethods(w io.Writer) {
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
			if ma.state != maState_initial {
				panic("misuse")
			}
			switch k {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case "{{ $field.Name }}":
				if ma.s & fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} != 0 {
					return nil, ipld.ErrRepeatedMapKey{&fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}}
				}
				ma.s += fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}
				ma.state = maState_midValue
				ma.f = {{ $i }} // REVIEW not sure if this matters in this path; suspect not.
				return &relevantChildValueAssembler, nil
			{{- end}}
			default:
				return nil, ipld.ErrInvalidKey{TypeName:"{{ $type.Name }}", Key:&_String{k}}
			}
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleKey() ipld.NodeAssembler {
			return ma
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleValue() ipld.NodeAssembler {
			return &relevantChildValueAssembler
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) Finish() error {
			panic("todo structbuilder mapassembler finish")
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) KeyStyle() ipld.NodeStyle {
			return _String__Style{}
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) ValueStyle(k string) ipld.NodeStyle {
			panic("todo structbuilder mapassembler valuestyle")
		}
	`, w, g.AdjCfg, g)
}
func (g structBuilderGenerator) emitKeyAssembler(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__KeyAssembler _{{ .Type | TypeSymbol }}__Assembler
	`, w, g.AdjCfg, g)
	stubs := mixins.StringAssemblerTraits{
		g.PkgName,
		g.TypeName,
		"", // unexercised here.
		"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Key",
	}
	stubs.EmitNodeAssemblerMethodBeginMap(w)
	stubs.EmitNodeAssemblerMethodBeginList(w)
	stubs.EmitNodeAssemblerMethodAssignNull(w)
	stubs.EmitNodeAssemblerMethodAssignBool(w)
	stubs.EmitNodeAssemblerMethodAssignInt(w)
	stubs.EmitNodeAssemblerMethodAssignFloat(w)
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__KeyAssembler) AssignString(k string) error {
			if ka.state != maState_midKey {
				panic("misuse")
			}
			switch k {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case "{{ $field.Name }}":
				if ka.s & fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} != 0 {
					return ipld.ErrRepeatedMapKey{&fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}}
				}
				ka.s += fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}
				ka.state = maState_expectValue
				ka.f = {{ $i }}
			{{- end}}
			default:
				return ipld.ErrInvalidKey{TypeName:"{{ $type.Name }}", Key:&_String{k}}
			}
			return nil
		}
	`, w, g.AdjCfg, g)
	stubs.EmitNodeAssemblerMethodAssignBytes(w)
	stubs.EmitNodeAssemblerMethodAssignLink(w)
	stubs.EmitNodeAssemblerMethodAssignNode(w)

}
func (g structBuilderGenerator) emitFieldValueAssembler(f schema.StructField, w io.Writer) {
	// TODO for Any, this should do a whole Thing;
	// TODO for any specific type, we should be able to tersely create a new type that embeds its assembler and wraps the one method that's valid for finishing its kind.

	// yes, do lets use unexported 'func()' callbacks for finishers, and thus generate WAY less types.
	// adding a word of memory to every builder?  not really.  the childbuilders would anyway for '*p' instead.  toplevels grow, true; but they frankly don't matter: 'one' doesn't show up on an asymtote.
	// is it an indirect call that can't be inlined?  yeah.  but... good lord; there's enough of them in the area already.  this one won't tip anything over.
	// one per field anyway?  probably, at least on the first pass.  easier to write; minimal consequence; can optimize size later.

	// '_MaybeT__Assembler'?  probably, yes -- the 'w' pointer has to be a different type, or things would get dumb.  ugh, every valid assign method is near-dupe but not quite.  shikata ga nai?
	// will there be '_MaybeT__ReprAssembler'?  probably, yes -- how would there not be?
	// Wild alternative: every '_T__Assembler' has a 'null bool' in addition to the 'finishCb func()'.  If used a context where null isn't valid: check it during finish.
	//  Scalar builders used as a root can do this via an overriden finisher method (literally, the AssignNull) on the Builder; easy peasy.
	//  Child assemblers can do it during the 'finishCb'.
	//  It's important to still return the possible rejection of null in the AssignNull method.
	//   But I think this shakes out: finishCb gets called at the end of *any* finisher method -- **that includes AssignNull**.
	//  If you're handling a MaybeT, the 'finishCb' can set the null state even though the '_T__Assembler' can't see the '_MaybeT.m'!
	//  Similarly: '_T__ReprAssembler' sprouts 'z bool' and 'finishCb func()' and does double duty.

	doTemplate(`
		// todo child assembler for field {{ .Name }}
	`, w, g.AdjCfg, f)
}
