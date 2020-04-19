package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type mapGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapTraits
	PkgName string
	Type    schema.TypeMap
}

// --- native content and specializations --->

func (g mapGenerator) EmitNativeType(w io.Writer) {
	// Maps do double bookkeeping.
	// - 'm' is used for quick lookup.
	// - 't' is used for both for order maintainence, and for allocation amortization for both keys and values.
	// Note that the key in 'm' is *not* a pointer.
	// The value in 'm' is a pointer into 't'.
	// REVIEW: does m's value need to be a pointer when it's a maybe?  perhaps not.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			m map[_{{ .Type.KeyType | TypeSymbol }}]*{{if .Type.ValueIsNullable }}Maybe{{else}}_{{end}}{{ .Type.ValueType | TypeSymbol }}
			t []_{{ .Type | TypeSymbol }}__entry
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
	// - address of 'k' is used when we return keys as nodes, such as in iterators.
	//    Having these in the 't' slice above amortizes moving all of them to heap at once,
	//     which makes iterators that have to return them as an interface much (much) lower cost -- no 'runtime.conv*' pain.
	// - address of 'v' is used in map values, to return, and of course also in iterators.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__entry struct {
			k _{{ .Type.KeyType | TypeSymbol }}
			v {{if .Type.ValueIsNullable }}Maybe{{else}}_{{end}}{{ .Type.ValueType | TypeSymbol }}
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNativeAccessors(w io.Writer) {
	// The speciated Lookup method *always* returns a MaybeT,
	//  because it can always be Absent;
	//   this saves us from needing multiple returns for an error.
	// If the MaybeT in question is UsePtr=false and the object is large and the map value isn't nullable,
	//  then this may cost some unfortunately large memcopies.
	// REVIEW: should we have a Lookup and a LookupMaybe function?  The former would be have no benefit over the latter if the value is nullable, but so what.
	doTemplate(`
		func (n *_{{ $type | TypeSymbol }}) LookupMaybe(k {{ .Type.KeyType | TypeSymbol }}) Maybe{{ .Type.ValueType | TypeSymbol }} {
			v, ok := n.m[*k]
			if !ok {
				return Maybe{{ .Type.ValueType | TypeSymbol }}{m:schema.Maybe_Absent}
			}
			{{- if .Type.ValueIsNullable }}
			return *v
			{{- else}}
			return Maybe{{ .Type.ValueType | TypeSymbol }}{m:schema.Maybe_Value, n:v}
			{{- end}}
		}
	`, w, g.AdjCfg, g)
	// FUTURE: also a speciated iterator?
}

func (g mapGenerator) EmitNativeBuilder(w io.Writer) {
	// Not yet clear what exactly might be most worth emitting here.
}

func (g mapGenerator) EmitNativeMaybe(w io.Writer) {
	// TODO these really seem to be maximally cookiecutter.  Can we extract them somewhere?
}

// --- type info --->

func (g mapGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g mapGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
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

func (g mapGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g mapGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
		var _ schema.TypedNode = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodLookupString(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupString(key string) (ipld.Node, error) {
			// TODO if complex key: use its fromString constructor
			// TODO if simple key: could use fromString for 'String', i guess?  it should inline?  or we can just do it bare obviously.
			//  ... yeah, let's standardize this.  these shouldn't even be visibly distinct cases, actually.
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodLookup(w io.Writer) {
	// FIXME maps can have complex keys
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Lookup(key ipld.Node) (ipld.Node, error) {
			// TODO cast-check it into our key type; flat barf if not match
			// REVIEW structs will coerce anything stringish silently...!  so we should figure out how to document the inconsistency here.
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodMapIterator(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) MapIterator() ipld.MapIterator {
			return &_{{ .Type | TypeSymbol }}__MapItr{n, 0}
		}

		type _{{ .Type | TypeSymbol }}__MapItr struct {
			n {{ .Type | TypeSymbol }}
			idx  int
		}

		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Next() (k ipld.Node, v ipld.Node, _ error) {
			if itr.idx >= len(itr.n.t) {
				return nil, nil, ipld.ErrIteratorOverread{}
			}
			x := itr.n.t[itr.idx]
			k = x.k
			{{- if .Type.ValueIsNullable }}
			switch x.v.m {
			case schema.Maybe_Null:
				v = ipld.Null
			case schema.Maybe_Value:
				v = {{ if not (MaybeUsesPtr .Type.ValueType) }}&{{end}}x.v.n
			}
			{{- else}}
			v = &x.v
			{{- end}}
			itr.idx++
			return
		}
		func (itr *_{{ .Type | TypeSymbol }}__MapItr) Done() bool {
			return itr.idx >= len(itr.n.t)
		}

	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Length() int {
			return len(n.t)
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodStyle(w io.Writer) {
	// REVIEW: this appears to be standard even across kinds; can we extract it?
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeStyleType(w io.Writer) {
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

func (g mapGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return mapBuilderGenerator{
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

type mapBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.MapAssemblerTraits
	PkgName string
	Type    schema.TypeMap
}

func (g mapBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			if nb.state != maState_finished {
				panic("invalid state: assembler for {{ .PkgName }}.{{ .Type.Name }} must be 'finished' before Build can be called!")
			}
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			var m schema.Maybe
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{w: &w, m: &m, state: maState_initial}}
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'm' is the **m**aybe which communicates our completeness to the parent if we're a child assembler.
	// - 'state' is what it says on the tin.  this is used for the map state (the broad transitions between null, start-map, and finish are handled by 'm' for consistency.)
	// - there's no equivalent of the 'f' (**f**ocused next) field in struct assemblers -- that's implicitly the last row of the 'w.t'.
	//
	// - 'cm' is **c**hild **m**aybe and is used for the completion message from children if values aren't allowed to be nullable and thus don't have their own per-value maybe slot we can use.
	// - 'ka' and 'va' are obviously the key assembler and value assembler respectively.
	//     Perhaps surprisingly, we can get away with using the value assembler for the value type just straight up, no wrappers necessary.
	//     TODO keys are probably not that simple, are they.
	//       Oh lordie.  We can't even do the 'w.m' assignment until key assembly is _finished_ can we.  if it's a complex key, that's... fun.  we'll need another unit of temp space: 'wk'.
	//       Keys can even be unions.  So they can *definitely* recurse.
	//       ... critical question: is it always the kind of direct recursion like struct stringjoin reprs, where it's nonstateful?
	//         Well, at the representation level, yes.  At the type level: no.  Cool!
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
			state maState

			{{- if not .Type.ValueIsNullable }}
			cm schema.Maybe
			{{- end}}
			ka _{{ .Type | TypeSymbol }}__KeyAssembler
			va _{{ .Type | TypeSymbol }}__ValueAssembler
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			switch *na.m {
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			case midvalue:
				panic("invalid state: it makes no sense to 'begin' twice on the same assembler!")
			}
			*na.m = midvalue
			if sizeHint < 0 {
				sizeHint = 0
			}
			{{- if .Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ .Type | TypeSymbol }}{}
			}
			{{- end}}
			na.w.m = make(map[_{{ .Type.KeyType | TypeSymbol }}]*{{if .Type.ValueIsNullable }}Maybe{{else}}_{{end}}{{ .Type.ValueType | TypeSymbol }}, sizeHint)
			na.w.t = make([]_{{ .Type | TypeSymbol }}__entry, 0, sizeHint)
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// DRY: this seems awfully similar -- almost exact, even -- amongst anything mapoid.  Can we extract?
	//  Might even be something we can turn into a util function, not just a template dry.  Only parameters are '*m', kind mixin, and type name.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNull() error {
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
func (g mapBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// AssignNode goes through three phases:
	// 1. is it null?  Jump over to AssignNull (which may or may not reject it).
	// 2. is it our own type?  Handle specially -- we might be able to do efficient things.
	// 3. is it the right kind to morph into us?  Do so.
	//
	// We do not set m=midvalue in phase 3 -- it shouldn't matter unless you're trying to pull off concurrent access, which is wrong and unsafe regardless.
	//
	// DRY: this seems awfully similar -- almost exact, even -- amongst anything mapoid.  Can we extract?
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
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
func (g mapBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	g.emitMapAssemblerChildTidyHelper(w)
	g.emitMapAssemblerMethods(w)
	g.emitKeyAssembler(w)
}
func (g mapBuilderGenerator) emitMapAssemblerChildTidyHelper(w io.Writer) {
	// This function attempts to clean up the state machine to acknolwedge child assembly finish.
	//  If the child was finished and we just collected it, return true and update state to maState_initial.
	//  Otherwise, if it wasn't done, return false;
	//   and the caller is almost certain to emit an error momentarily.
	// The function will only be called when the current state is maState_midValue.
	//  (In general, the idea is that if the user is doing things correctly,
	//   this function will only be called when the child is in fact finished.)
	// If 'cm' is used, we reset it to its initial condition of Maybe_Absent here.
	//  At the same time, we nil the 'w' pointer for the child assembler; otherwise its own state machine would probably let it modify 'w' again!
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) valueFinishTidy() bool {
			tz := &ma.w.t[len(ma.w.t)-1]
			switch tz.m {
			{{- if .Type.ValueIsNullable }}
			case schema.Maybe_Null:
				ma.state = maState_initial
				return true
			case schema.Maybe_Value:
				{{- if (MaybeUsesPtr $field.Type) }}
				tz.v = ma.va.w
				{{- end}}
				ma.state = maState_initial
				return true
			{{- else}}
			case schema.Maybe_Value:
				ma.va.w = nil
				ma.cm = schema.Maybe_Absent
				ma.state = maState_initial
			{{- end}}
			default:
				return false
			}
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) emitMapAssemblerMethods(w io.Writer) {
	// FUTURE: some of the setup of the child assemblers could probably be DRY'd up.
	// DRY: a lot of the state transition fences again are common for all mapoids, and could probably even be a function over '*state'... crap, except for the valueFinishTidy function, which is definitely not extractable.
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleEntry(k string) (ipld.NodeAssembler, error) {
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

			// TODO need key reification again
			//  or invalid key: return nil, ipld.ErrInvalidKey{TypeName:"{{ .PkgName }}.{{ .Type.Name }}", Key:&_String{k}}
			if _, ok := ma.w.m[k2] {
				return ipld.ErrRepeatedMapKey{k}
			}
			ma.w.t = append(ma.w.t, _{{ .Type | TypeSymbol }}__entry{k: k2})
			ma.w.m[k2] = &ma.w.t[len(ma.w.t)-1].v
			ma.state = maState_midValue

			{{- if .Type.ValueIsNullable }}
			ma.va.w = ma.w.t[len(ma.w.t)-1].v.n
			ma.va.m = &ma.w.t[len(ma.w.t)-1].v.m
			ma.w.t[len(ma.w.t)-1].v.m = allowNull
			{{- else}}
			ma.va.w = ma.w.t[len(ma.w.t)-1].v.n
			ma.va.m = &ma.cm
			{{- end}}
			return &ma.va, nil
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleKey() ipld.NodeAssembler {
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
			return (*_{{ .Type | TypeSymbol }}__KeyAssembler)(ma)
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleValue() ipld.NodeAssembler {
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
			return &ma.va, nil
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) Finish() error {
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
			ma.state = maState_finished
			*ma.m = schema.Maybe_Value
			return nil
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) KeyStyle() ipld.NodeStyle {
			return _{{ .Type.KeyType | TypeSymbol }}__Style{}
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) ValueStyle(_ string) ipld.NodeStyle {
			return _{{ .Type.ValueType | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) emitKeyAssembler(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__KeyAssembler _{{ .Type | TypeSymbol }}__Assembler
	`, w, g.AdjCfg, g)
	stubs := mixins.StringAssemblerTraits{
		g.PkgName,
		g.TypeName + ".KeyAssembler",
		"_" + g.AdjCfg.TypeSymbol(g.Type) + "__Key",
	}
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__KeyAssembler) BeginMap(_ int) error {
			if ka.state != maState_midKey {
				panic("misuse: KeyAssembler held beyond its valid lifetime")
			}
			// TODO this could turn into either a struct (as long as it has a string repr) or a union (same rule (implicitly, can't be kinded))...
			//   but string reprs don't matter here: both of these act like maps at the type level.
		}
	`, w, g.AdjCfg, g)
	stubs.EmitNodeAssemblerMethodBeginList(w)
	stubs.EmitNodeAssemblerMethodAssignNull(w)
	stubs.EmitNodeAssemblerMethodAssignBool(w)
	stubs.EmitNodeAssemblerMethodAssignInt(w)
	stubs.EmitNodeAssemblerMethodAssignFloat(w)
	// DRY: this is almost entirely the same as the body of the AssembleEntry method except for the final return and where we put 'ka.state'.  extract.
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__KeyAssembler) AssignString(k string) error {
			if ka.state != maState_midKey {
				panic("misuse: KeyAssembler held beyond its valid lifetime")
			}

			// TODO need key reification again
			//  or invalid key: return nil, ipld.ErrInvalidKey{TypeName:"{{ .PkgName }}.{{ .Type.Name }}", Key:&_String{k}}
			if _, ok := ka.w.m[k2] {
				return ipld.ErrRepeatedMapKey{k}
			}
			ka.w.t = append(ka.w.t, _{{ .Type | TypeSymbol }}__entry{k: k2})
			ka.w.m[k2] = &ka.w.t[len(ka.w.t)-1].v
			ka.state = maState_expectValue

			{{- if .Type.ValueIsNullable }}
			ma.va.w = ma.w.t[len(ma.w.t)-1].v.n
			ma.va.m = &ma.w.t[len(ma.w.t)-1].v.m
			ma.w.t[len(ma.w.t)-1].v.m = allowNull
			{{- else}}
			ma.va.w = ma.w.t[len(ma.w.t)-1].v.n
			ma.va.m = &ma.cm
			{{- end}}
			return nil
		}
	`, w, g.AdjCfg, g)
	stubs.EmitNodeAssemblerMethodAssignBytes(w)
	stubs.EmitNodeAssemblerMethodAssignLink(w)
	doTemplate(`
		func (ka *_{{ .Type | TypeSymbol }}__KeyAssembler) AssignNode(v ipld.Node) error {
			// TODO more key reification
		}
		func (_{{ .Type | TypeSymbol }}__KeyAssembler) Style() ipld.NodeStyle {
			return _{{ .Type.KeyType | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}
