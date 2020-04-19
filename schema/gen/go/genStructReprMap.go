package gengo

import (
	"io"
	"strconv"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

var _ TypeGenerator = &structReprMapGenerator{}

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

	// We do also want some constants for our fields;
	//  they'll make iterators able to work faster.
	//  These might be the same strings as the type-level field names
	//   (in face, they are, unless renames are used)... but that's fine.
	//    We get simpler code by just doing this unconditionally.
	doTemplate(`
		var (
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $field := .Type.Fields }}
			fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}_serial = _String{"{{ $field | $type.RepresentationStrategy.GetFieldKey }}"}
			{{- end }}
		)
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
				{{- if $field.IsMaybe }}
				return n.{{ $field | FieldSymbolLower }}.v.Representation(), nil
				{{- else}}
				return n.{{ $field | FieldSymbolLower }}.Representation(), nil
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

	// First: Determine if there are any optionals at all.
	//  If there are none, some control flow symbols need to not be emitted.
	fields := g.Type.Fields()
	haveOptionals := false
	for _, field := range fields {
		if field.IsOptional() {
			haveOptionals = true
			break
		}
	}

	// Second: Count how many trailing fields are optional.
	//  The 'Done' predicate gets more complex when in the trailing optionals.
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
		{{ if .HaveOptionals }}advance:{{end -}}
			if itr.idx >= {{ len .Type.Fields }} {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			switch itr.idx {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				k = &fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}_serial
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
				{{- if $field.IsMaybe }}
				v = itr.n.{{ $field | FieldSymbolLower}}.v.Representation()
				{{- else}}
				v = itr.n.{{ $field | FieldSymbolLower}}.Representation()
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
		HaveOptionals              bool
		HaveTrailingOptionals      bool
		BeginTrailingOptionalField int
	}{
		g.Type,
		haveOptionals,
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
			var m schema.Maybe
			*nb = _{{ .Type | TypeSymbol }}__ReprBuilder{_{{ .Type | TypeSymbol }}__ReprAssembler{w: &w, m: &m, state: maState_initial}}
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'm' is the **m**aybe which communicates our completeness to the parent if we're a child assembler.
	// - 'state' is what it says on the tin.  this is used for the map state (the broad transitions between null, start-map, and finish are handled by 'm' for consistency.)
	// - 's' is a bitfield for what's been **s**et.
	// - 'f' is the **f**ocused field that will be assembled next.
	//
	// - 'cm' is **c**hild **m**aybe and is used for the completion message from children that aren't allowed to be nullable (for those that are, their own maybe.m is used).
	// - the 'ca_*' fields embed **c**hild **a**ssemblers -- these are embedded so we can yield pointers to them without causing new allocations.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprAssembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
			state maState
			s int
			f int

			cm schema.Maybe
			{{range $field := .Type.Fields -}}
			ca_{{ $field | FieldSymbolLower }} _{{ $field.Type | TypeSymbol }}__ReprAssembler
			{{end -}}
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	// We currently disregard sizeHint.  It's not relevant to us.
	//  We could check it strictly and emit errors; presently, we don't.
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) BeginMap(int) (ipld.MapAssembler, error) {
			switch *na.m {
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: it makes no sense to 'begin' twice on the same assembler!")
			}
			*na.m = midvalue
			{{- if .Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ .Type | TypeSymbol }}{}
			}
			{{- end}}
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNull() error {
			switch *na.m {
			case allowNull:
				*na.m = schema.Maybe_Null
				return nil
			case schema.Maybe_Absent:
				return mixins.MapAssembler{"{{ .PkgName }}.{{ .TypeName }}"}.AssignNull()
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: cannot assign null into an assembler that's already begun working on recursive structures!")
			}
			panic("unreachable")
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// AssignNode goes through three phases:
	// 1. is it null?  Jump over to AssignNull (which may or may not reject it).
	// 2. is it our own type?  Handle specially -- we might be able to do efficient things.
	// 3. is it the right kind to morph into us?  Do so.
	//
	// We do not set m=midvalue in phase 3 -- it shouldn't matter unless you're trying to pull off concurrent access, which is wrong and unsafe regardless.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__ReprAssembler) AssignNode(v ipld.Node) error {
			if v.IsNull() {
				return na.AssignNull()
			}
			if v2, ok := v.(*_{{ .Type | TypeSymbol }}); ok {
				switch *na.m {
				case schema.Maybe_Value, schema.Maybe_Null:
					panic("invalid state: cannot assign into assembler that's already finished")
				case midvalue:
					panic("invalid state: cannot assign null into an assembler that's already begun working on recursive structures!")
				}
				{{- if .Type | MaybeUsesPtr }}
				if na.w == nil {
					na.w = v2
					*na.m = schema.Maybe_Value
					return nil
				}
				{{- end}}
				*na.w = *v2
				*na.m = schema.Maybe_Value
				return nil
			}
			if v.ReprKind() != ipld.ReprKind_Map {
				return ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .Type.Name }}.Repr", MethodName: "AssignNode", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: v.ReprKind()}
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
func (g structReprMapReprBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	g.emitMapAssemblerChildTidyHelper(w)
	g.emitMapAssemblerMethods(w)
	g.emitKeyAssembler(w)
}
func (g structReprMapReprBuilderGenerator) emitMapAssemblerChildTidyHelper(w io.Writer) {
	// This is exactly the same as the matching method on the type-level assembler;
	//  everything that differs happens to be hidden behind the 'f' indirection, which is numeric.
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) valueFinishTidy() bool {
			switch ma.f {
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				{{- if $field.IsNullable }}
				switch ma.w.{{ $field | FieldSymbolLower }}.m {
				case schema.Maybe_Null:
					ma.state = maState_initial
					return true
				case schema.Maybe_Value:
					{{- if (MaybeUsesPtr $field.Type) }}
					ma.w.{{ $field | FieldSymbolLower }}.v = ma.ca_{{ $field | FieldSymbolLower }}.w
					{{- end}}
					ma.state = maState_initial
					return true
				default:
					return false
				}
				{{- else if $field.IsOptional }}
				switch ma.w.{{ $field | FieldSymbolLower }}.m {
				case schema.Maybe_Value:
					{{- if (MaybeUsesPtr $field.Type) }}
					ma.w.{{ $field | FieldSymbolLower }}.v = ma.ca_{{ $field | FieldSymbolLower }}.w
					{{- end}}
					ma.state = maState_initial
					return true
				default:
					return false
				}
				{{- else}}
				switch ma.cm {
				case schema.Maybe_Value:
					{{- /* while defense in depth here might avoid some 'wat' outcomes, it's not strictly necessary for safety */ -}}
					{{- /* ma.ca_{{ $field | FieldSymbolLower }}.w = nil */ -}}
					{{- /* ma.ca_{{ $field | FieldSymbolLower }}.m = nil */ -}}
					ma.cm = schema.Maybe_Absent
					ma.state = maState_initial
					return true
				default:
					return false
				}
				{{- end}}
			{{- end}}
			default:
				panic("unreachable")
			}
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) emitMapAssemblerMethods(w io.Writer) {
	// FUTURE: some of the setup of the child assemblers could probably be DRY'd up.
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
			switch ma.state {
			case maState_initial:
				// carry on
			case maState_midKey:
				panic("invalid state: AssembleEntry cannot be called when in the middle of assembling another key")
			case maState_expectValue:
				panic("invalid state: AssembleEntry cannot be called when expecting start of value assembly")
			case maState_midValue:
				if !ma.valueFinishTidy() {
					panic("invalid state: AssembleEntry cannot be called when in the middle of assembling a value")
				} // if tidy success: carry on
			case maState_finished:
				panic("invalid state: AssembleEntry cannot be called on an assembler that's already finished")
			}
			switch k {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case "{{ $field | $type.RepresentationStrategy.GetFieldKey }}":
				if ma.s & fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} != 0 {
					return nil, ipld.ErrRepeatedMapKey{&fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}_serial}
				}
				ma.s += fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}
				ma.state = maState_midValue
				{{- if $field.IsMaybe }}
				ma.ca_{{ $field | FieldSymbolLower }}.w = {{if not (MaybeUsesPtr $field.Type) }}&{{end}}ma.w.{{ $field | FieldSymbolLower }}.v
				ma.ca_{{ $field | FieldSymbolLower }}.m = &ma.w.{{ $field | FieldSymbolLower }}.m
				{{if $field.IsNullable }}ma.w.{{ $field | FieldSymbolLower }}.m = allowNull{{end}}
				{{- else}}
				ma.ca_{{ $field | FieldSymbolLower }}.w = &ma.w.{{ $field | FieldSymbolLower }}
				ma.ca_{{ $field | FieldSymbolLower }}.m = &ma.cm
				{{- end}}
				return &ma.ca_{{ $field | FieldSymbolLower }}, nil
			{{- end}}
			default:
				return nil, ipld.ErrInvalidKey{TypeName:"{{ .PkgName }}.{{ .Type.Name }}.Repr", Key:&_String{k}}
			}
		}
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) AssembleKey() ipld.NodeAssembler {
			switch ma.state {
			case maState_initial:
				// carry on
			case maState_midKey:
				panic("invalid state: AssembleKey cannot be called when in the middle of assembling another key")
			case maState_expectValue:
				panic("invalid state: AssembleKey cannot be called when expecting start of value assembly")
			case maState_midValue:
				if !ma.valueFinishTidy() {
					panic("invalid state: AssembleKey cannot be called when in the middle of assembling a value")
				} // if tidy success: carry on
			case maState_finished:
				panic("invalid state: AssembleKey cannot be called on an assembler that's already finished")
			}
			ma.state = maState_midKey
			return (*_{{ .Type | TypeSymbol }}__ReprKeyAssembler)(ma)
		}
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) AssembleValue() ipld.NodeAssembler {
			switch ma.state {
			case maState_initial:
				panic("invalid state: AssembleValue cannot be called when no key is primed")
			case maState_midKey:
				panic("invalid state: AssembleValue cannot be called when in the middle of assembling a key")
			case maState_expectValue:
				// carry on
			case maState_midValue:
				panic("invalid state: AssembleValue cannot be called when in the middle of assembling another value")
			case maState_finished:
				panic("invalid state: AssembleValue cannot be called on an assembler that's already finished")
			}
			ma.state = maState_midValue
			switch ma.f {
			{{- range $i, $field := .Type.Fields }}
			case {{ $i }}:
				{{- if $field.IsMaybe }}
				ma.ca_{{ $field | FieldSymbolLower }}.w = {{if not (MaybeUsesPtr $field.Type) }}&{{end}}ma.w.{{ $field | FieldSymbolLower }}.v
				ma.ca_{{ $field | FieldSymbolLower }}.m = &ma.w.{{ $field | FieldSymbolLower }}.m
				{{if $field.IsNullable }}ma.w.{{ $field | FieldSymbolLower }}.m = allowNull{{end}}
				{{- else}}
				ma.ca_{{ $field | FieldSymbolLower }}.w = &ma.w.{{ $field | FieldSymbolLower }}
				ma.ca_{{ $field | FieldSymbolLower }}.m = &ma.cm
				{{- end}}
				return &ma.ca_{{ $field | FieldSymbolLower }}
			{{- end}}
			default:
				panic("unreachable")
			}
		}
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) Finish() error {
			switch ma.state {
			case maState_initial:
				// carry on
			case maState_midKey:
				panic("invalid state: Finish cannot be called when in the middle of assembling a key")
			case maState_expectValue:
				panic("invalid state: Finish cannot be called when expecting start of value assembly")
			case maState_midValue:
				if !ma.valueFinishTidy() {
					panic("invalid state: Finish cannot be called when in the middle of assembling a value")
				} // if tidy success: carry on
			case maState_finished:
				panic("invalid state: Finish cannot be called on an assembler that's already finished")
			}
			//FIXME check if all required fields are set
			ma.state = maState_finished
			*ma.m = schema.Maybe_Value
			return nil
		}
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) KeyStyle() ipld.NodeStyle {
			return _String__Style{}
		}
		func (ma *_{{ .Type | TypeSymbol }}__ReprAssembler) ValueStyle(k string) ipld.NodeStyle {
			panic("todo structbuilder mapassembler repr valuestyle")
		}
	`, w, g.AdjCfg, g)
}
func (g structReprMapReprBuilderGenerator) emitKeyAssembler(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__ReprKeyAssembler _{{ .Type | TypeSymbol }}__ReprAssembler
	`, w, g.AdjCfg, g)
	stubs := mixins.StringAssemblerTraits{
		g.PkgName,
		g.TypeName + ".ReprKeyAssembler",
		"_" + g.AdjCfg.TypeSymbol(g.Type) + "__ReprKey",
	}
	// This key assembler can disregard any idea of complex keys because it's at the representation level!
	//  Map keys must always be plain strings at the representation level.
	stubs.EmitNodeAssemblerMethodBeginMap(w)
	stubs.EmitNodeAssemblerMethodBeginList(w)
	stubs.EmitNodeAssemblerMethodAssignNull(w)
	stubs.EmitNodeAssemblerMethodAssignBool(w)
	stubs.EmitNodeAssemblerMethodAssignInt(w)
	stubs.EmitNodeAssemblerMethodAssignFloat(w)
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__ReprKeyAssembler) AssignString(k string) error {
			if ka.state != maState_midKey {
				panic("misuse: KeyAssembler held beyond its valid lifetime")
			}
			switch k {
			{{- $type := .Type -}} {{- /* ranging modifies dot, unhelpfully */ -}}
			{{- range $i, $field := .Type.Fields }}
			case "{{ $field | $type.RepresentationStrategy.GetFieldKey }}":
				if ka.s & fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }} != 0 {
					return ipld.ErrRepeatedMapKey{&fieldName__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}_serial}
				}
				ka.s += fieldBit__{{ $type | TypeSymbol }}_{{ $field | FieldSymbolUpper }}
				ka.state = maState_expectValue
				ka.f = {{ $i }}
			{{- end}}
			default:
				return ipld.ErrInvalidKey{TypeName:"{{ .PkgName }}.{{ .Type.Name }}.Repr", Key:&_String{k}}
			}
			return nil
		}
	`, w, g.AdjCfg, g)
	stubs.EmitNodeAssemblerMethodAssignBytes(w)
	stubs.EmitNodeAssemblerMethodAssignLink(w)
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__ReprKeyAssembler) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				return ka.AssignString(v2)
			}
		}
		func (_{{ .Type | TypeSymbol }}__ReprKeyAssembler) Style() ipld.NodeStyle {
			return _String__Style{}
		}
	`, w, g.AdjCfg, g)
}
