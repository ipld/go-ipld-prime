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
				{{- if $field.IsMaybe }}
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
				{{- if $field.IsMaybe }}
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
	// While for scalars we cover up all the assembler methods that are prepared to handle null,
	//  here, it's unfeasible: would result in deep bifrucations.
}
func (g structBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'state' is what it says on the tin.
	// - 's' is a bitfield for what's been **s**et.
	// - 'f' is the **f**ocused field that will be assembled next.
	// - 'z' is used to denote a null (in case we're used in a context that's acceptable).  z for **z**ilch.
	// - 'fcb' is the **f**inish **c**all**b**ack, supplied by the parent if we're a child assembler.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			state maState
			s int
			f int
			z bool
			fcb func() error

			{{range $field := .Type.Fields -}}
			ca_{{ $field | FieldSymbolLower }} _{{ $field.Type | TypeSymbol }}__Assembler
			{{end -}}
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
func (g structBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// Twist: All generated assemblers quietly accept null... and if used in a context they shouldn't,
	//  either this method gets overriden (this is the case for NodeBuilders),
	//  or the 'na.fcb' (short for "finish callback") returns the rejection.
	//  We *do* need a nil check for 'fcb' (unlike with scalars, there's just too much dupe code if we try to make builders inherit no methods).
	//  We don't pass any args to 'fcb' because we assume it comes from something that can already see this whole struct.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNull() error {
			if na.state != maState_initial {
				panic("misuse")
			}
			na.z = true
			na.state = maState_finished
			if na.fcb != nil {
				return na.fcb()
			} else {
				return nil
			}
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
				if na.fcb != nil {
					return na.fcb()
				} else {
					return nil
				}
			}
			if v.ReprKind() == ipld.ReprKind_Null {
				return na.AssignNull()
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
	// FUTURE: some of the setup of the child assemblers could probably be DRY'd up.
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
				{{- if $field.IsMaybe }}
				ma.ca_{{ $field | FieldSymbolLower }}.w = {{if not (MaybeUsesPtr $field.Type) }}&{{end}}ma.w.{{ $field | FieldSymbolLower }}.v
				{{- else}}
				ma.ca_{{ $field | FieldSymbolLower }}.w = &ma.w.{{ $field | FieldSymbolLower }}
				{{- end}}
				ma.ca_{{ $field | FieldSymbolLower }}.fcb = ma.fcb_{{ $field | FieldSymbolLower }}
				return &ma.ca_{{ $field | FieldSymbolLower }}, nil
			{{- end}}
			default:
				return nil, ipld.ErrInvalidKey{TypeName:"{{ $type.Name }}", Key:&_String{k}}
			}
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleKey() ipld.NodeAssembler {
			if ma.state != maState_initial {
				panic("misuse")
			}
			ma.state = maState_midKey
			return (*_{{ .Type | TypeSymbol }}__KeyAssembler)(ma)
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleValue() ipld.NodeAssembler {
			if ma.state != maState_expectValue {
				panic("misuse")
			}
			ma.state = maState_midValue
			switch ma.f {
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				{{- if $field.IsMaybe }}
				ma.ca_{{ $field | FieldSymbolLower }}.w = {{if not (MaybeUsesPtr $field.Type) }}&{{end}}ma.w.{{ $field | FieldSymbolLower }}.v
				{{- else}}
				ma.ca_{{ $field | FieldSymbolLower }}.w = &ma.w.{{ $field | FieldSymbolLower }}
				{{- end}}
				ma.ca_{{ $field | FieldSymbolLower }}.fcb = ma.fcb_{{ $field | FieldSymbolLower }}
				return &ma.ca_{{ $field | FieldSymbolLower }}
			{{- end}}
			default:
				panic("unreachable")
			}
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) Finish() error {
			if ma.state != maState_initial {
				panic("misuse")
			}
			//FIXME check if all required fields are set
			ma.state = maState_finished
			if ma.fcb != nil {
				return ma.fcb()
			} else {
				return nil
			}
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
		g.TypeName + ".KeyAssembler",
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
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__KeyAssembler) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				return ka.AssignString(v2)
			}
		}
		func (_{{ .Type | TypeSymbol }}__KeyAssembler) Style() ipld.NodeStyle {
			return _String__Style{}
		}
	`, w, g.AdjCfg, g)

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

	// Generate the finisher callbacks that we'll insert into child assemblers.
	//  In all cases, it's responsible for advancing the parent's state machine.
	//  If it's not a maybe:
	//   If the child assembler 'z' is true, raise error.
	//  If it's a maybe:
	//   If it's not nullable:
	//    If the child assembler 'z' is true, raise error.
	//    If the maybe uses pointer types, we copy out the pointer.
	//   If it is nullable:
	//    If the child assembler 'z' is true, set 'm' to 'Null'.
	//    If the child assembler 'z' is false:
	//     If the maybe uses pointer types, we copy out the pointer.
	//      (The pointer will not be nil.  Absent would've been on a divergent track much earlier.)
	//     (If the maybe doesn't use pointer types, the pointer is to the embed location; no action needed.)
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) fcb_{{ .Field | FieldSymbolLower }}() error {
			{{- if .Field.IsNullable }}
			if na.ca_{{ .Field | FieldSymbolLower }}.z == true {
				na.w.{{ .Field | FieldSymbolLower }}.m = schema.Maybe_Null
			} else {
				na.w.{{ .Field | FieldSymbolLower }}.m = schema.Maybe_Value
			}
			{{- else}}
			if na.ca_{{ .Field | FieldSymbolLower }}.z == true {
				return mixins.MapAssembler{"{{ .PkgName }}.{{ .Field.Type.Name }}"}.AssignNull()
			}
			{{- end}}
			{{- if and .Field.IsOptional (not .Field.IsNullable) }}
			na.w.{{ .Field | FieldSymbolLower }}.m = schema.Maybe_Value
			{{- end}}
			{{- if and .Field.IsMaybe (.Field.Type | MaybeUsesPtr) }}
			na.w.{{ .Field | FieldSymbolLower }}.v = na.ca_{{ .Field | FieldSymbolLower }}.w
			{{- end}}
			na.ca_{{ .Field | FieldSymbolLower }}.w = nil
			na.state = maState_initial
			return nil
		}
	`, w, g.AdjCfg, struct {
		PkgName  string
		Type     schema.TypeStruct
		TypeName string
		Field    schema.StructField
	}{
		g.PkgName,
		g.Type,
		g.TypeName,
		f,
	})
}
