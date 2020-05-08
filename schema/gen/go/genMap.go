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

func (mapGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

// --- native content and specializations --->

func (g mapGenerator) EmitNativeType(w io.Writer) {
	// Maps do double bookkeeping.
	// - 'm' is used for quick lookup.
	// - 't' is used for both for order maintainence, and for allocation amortization for both keys and values.
	// Note that the key in 'm' is *not* a pointer.
	// The value in 'm' is a pointer into 't' (except when it's a maybe; maybes are already pointers).
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			m map[_{{ .Type.KeyType | TypeSymbol }}]{{if .Type.ValueIsNullable }}Maybe{{else}}*_{{end}}{{ .Type.ValueType | TypeSymbol }}
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
			v _{{ .Type.ValueType | TypeSymbol }}{{if .Type.ValueIsNullable }}__Maybe{{end}}
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNativeAccessors(w io.Writer) {
	// Generate a speciated Lookup as well as LookupMaybe method.
	// The LookupMaybe method is needed if the map value is nullable and you're going to distinguish nulls
	//  (and may also be convenient if you would rather handle Maybe_Absent than an error for not-found).
	// The Lookup method works fine if the map value isn't nullable
	//  (and should be preferred in that case, because boxing something into a maybe when it wasn't already stored that way costs an alloc(!),
	//   and may additionally incur a memcpy if the maybe for the value type doesn't use pointers internally).
	// REVIEW: is there a way we can make this less twisty?  it is VERY unfortunate if the user has to know what sort of map it is to know which method to prefer.
	//  Maybe the Lookup method on maps that have nullable values should just always have a MaybeT return type?
	//   But then this means the Lookup method doesn't "need" an error as part of its return signiture, which just shuffles differences around.
	doTemplate(`
		func (n *_{{ .Type | TypeSymbol }}) LookupMaybe(k {{ .Type.KeyType | TypeSymbol }}) Maybe{{ .Type.ValueType | TypeSymbol }} {
			v, ok := n.m[*k]
			if !ok {
				return &_{{ .Type | TypeSymbol }}__valueAbsent
			}
			{{- if .Type.ValueIsNullable }}
			return v
			{{- else}}
			return &_{{ .Type.ValueType | TypeSymbol }}__Maybe{
				m: schema.Maybe_Value,
				v: {{ if not (MaybeUsesPtr .Type.ValueType) }}*{{end}}v,
			}
			{{- end}}
		}

		var _{{ .Type | TypeSymbol }}__valueAbsent = _{{ .Type.ValueType | TypeSymbol }}__Maybe{m:schema.Maybe_Absent}

		// TODO generate also a plain Lookup method that doesn't box and alloc if this type contains non-nullable values!
	`, w, g.AdjCfg, g)
	// FUTURE: also a speciated iterator?
}

func (g mapGenerator) EmitNativeBuilder(w io.Writer) {
	// Not yet clear what exactly might be most worth emitting here.
}

func (g mapGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
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
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g mapGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g mapGenerator) EmitNodeTypeAssertions(w io.Writer) {
	emitNodeTypeAssertions_typical(w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodLookupString(w io.Writer) {
	// What should be coercible in which directions (and how surprising that is) is an interesting question.
	//  Most of the answer comes from considering what needs to be possible when working with PathSegment:
	//   we *must* be able to accept a string in a PathSegment and be able to use it to navigate a map -- even if the map has complex keys.
	//   For that to work out, it means if the key type doesn't have a string type kind, we must be willing to reach into its representation and use the fromString there.
	//  If the key type *does* have a string kind at the type level, we'll use that; no need to consider going through the representation.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupString(k string) (ipld.Node, error) {
			var k2 _{{ .Type.KeyType | TypeSymbol }}
			{{- if eq .Type.KeyType.Kind.String "String" }}
			if err := (_{{ .Type.KeyType | TypeSymbol }}__Style{}).fromString(&k2, k); err != nil {
				return nil, err // TODO wrap in some kind of ErrInvalidKey
			}
			{{- else}}
			if err := (_{{ .Type.KeyType | TypeSymbol }}__ReprStyle{}).fromString(&k2, k); err != nil {
				return nil, err // TODO wrap in some kind of ErrInvalidKey
			}
			{{- end}}
			v, exists := n.m[k2]
			if !exists {
				return ipld.Undef, ipld.ErrNotExists{ipld.PathSegmentOfString(k)}
			}
			{{- if .Type.ValueIsNullable }}
			if v.m == schema.Maybe_Null {
				return ipld.Null, nil
			}
			return {{ if not (MaybeUsesPtr .Type.ValueType) }}&{{end}}v.v, nil
			{{- else}}
			return v, nil
			{{- end}}
		}
	`, w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeMethodLookup(w io.Writer) {
	// LookupNode will procede by cast if it can; or simply error if that doesn't work.
	//  There's no attempt to turn the node (or its repr) into a string and then reify that into a key;
	//   if you used a Node here, you should've meant it.
	// REVIEW: by comparison structs will coerce anything stringish silently...!  so we should figure out if that inconsistency is acceptable, and at least document it if so.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Lookup(k ipld.Node) (ipld.Node, error) {
			k2, ok := k.({{ .Type.KeyType | TypeSymbol }})
			if !ok {
				panic("todo invalid key type error")
				// 'ipld.ErrInvalidKey{TypeName:"{{ .PkgName }}.{{ .Type.Name }}", Key:&_String{k}}' doesn't quite cut it: need room to explain the type, and it's not guaranteed k can be turned into a string at all
			}
			v, exists := n.m[*k2]
			if !exists {
				return ipld.Undef, ipld.ErrNotExists{ipld.PathSegmentOfString(k2.String())}
			}
			{{- if .Type.ValueIsNullable }}
			if v.m == schema.Maybe_Null {
				return ipld.Null, nil
			}
			return {{ if not (MaybeUsesPtr .Type.ValueType) }}&{{end}}v.v, nil
			{{- else}}
			return v, nil
			{{- end}}
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
			x := &itr.n.t[itr.idx]
			k = &x.k
			{{- if .Type.ValueIsNullable }}
			switch x.v.m {
			case schema.Maybe_Null:
				v = ipld.Null
			case schema.Maybe_Value:
				v = {{ if not (MaybeUsesPtr .Type.ValueType) }}&{{end}}x.v.v
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
	emitNodeMethodStyle_typical(w, g.AdjCfg, g)
}

func (g mapGenerator) EmitNodeStyleType(w io.Writer) {
	emitNodeStyleType_typical(w, g.AdjCfg, g)
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

func (mapBuilderGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

func (g mapBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'm' is the **m**aybe which communicates our completeness to the parent if we're a child assembler.
	// - 'state' is what it says on the tin.  this is used for the map state (the broad transitions between null, start-map, and finish are handled by 'm' for consistency.)
	// - there's no equivalent of the 'f' (**f**ocused next) field in struct assemblers -- that's implicitly the last row of the 'w.t'.
	//
	// - 'cm' is **c**hild **m**aybe and is used for the completion message from children.
	//    It's used for values if values aren't allowed to be nullable and thus don't have their own per-value maybe slot we can use.
	//    It's always used for key assembly, since keys are never allowed to be nullable and thus etc.
	// - 'ka' and 'va' are the key assembler and value assembler respectively.
	//    Perhaps surprisingly, we can get away with using the assemblers for each type just straight up, no wrappers necessary;
	//     All of the required magic is handled through maybe pointers and some tidy methods used during state transitions.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
			state maState

			cm schema.Maybe
			ka _{{ .Type.KeyType | TypeSymbol }}__Assembler
			va _{{ .Type.ValueType | TypeSymbol }}__Assembler
		}

		func (na *_{{ .Type | TypeSymbol }}__Assembler) reset() {
			na.state = maState_initial
			na.ka.reset()
			na.va.reset()
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
			na.w.m = make(map[_{{ .Type.KeyType | TypeSymbol }}]{{if .Type.ValueIsNullable }}Maybe{{else}}*_{{end}}{{ .Type.ValueType | TypeSymbol }}, sizeHint)
			na.w.t = make([]_{{ .Type | TypeSymbol }}__entry, 0, sizeHint)
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	emitNodeAssemblerMethodAssignNull_recursive(w, g.AdjCfg, g)
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
	g.emitMapAssemblerKeyTidyHelper(w)
	g.emitMapAssemblerValueTidyHelper(w)
	g.emitMapAssemblerMethods(w)
}
func (g mapBuilderGenerator) emitMapAssemblerKeyTidyHelper(w io.Writer) {
	// This function attempts to clean up the state machine to acknolwedge key assembly finish.
	//  If the child was finished and we just collected it, return true and update state to maState_expectValue.
	//   Collecting the child includes updating the 'ma.w.m' to point into the relevant row of 'ma.w.t', since that couldn't be done earlier,
	//    AND initializing the 'ma.va' (since we're already holding relevant offsets into 'ma.w.t').
	//  Otherwise, if it wasn't done, return false;
	//   and the caller is almost certain to emit an error momentarily.
	// The function will only be called when the current state is maState_midKey.
	//  (In general, the idea is that if the user is doing things correctly,
	//   this function will only be called when the child is in fact finished.)
	// Completion info always comes via 'cm', and we reset it to its initial condition of Maybe_Absent here.
	//  At the same time, we nil the 'w' pointer for the child assembler; otherwise its own state machine would probably let it modify 'w' again!
	doTemplate(`
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) keyFinishTidy() bool {
			switch ma.cm {
			case schema.Maybe_Value:
				ma.ka.w = nil
				tz := &ma.w.t[len(ma.w.t)-1]
				ma.cm = schema.Maybe_Absent
				ma.state = maState_expectValue
				ma.w.m[tz.k] = &tz.v
				{{- if .Type.ValueIsNullable }}
				{{- if not (MaybeUsesPtr .Type.ValueType) }}
				ma.va.w = &tz.v.v
				{{- end}}
				ma.va.m = &tz.v.m
				tz.v.m = allowNull
				{{- else}}
				ma.va.w = &tz.v
				ma.va.m = &ma.cm
				{{- end}}
				ma.ka.reset()
				return true
			default:
				return false
			}
		}
	`, w, g.AdjCfg, g)
}
func (g mapBuilderGenerator) emitMapAssemblerValueTidyHelper(w io.Writer) {
	// This function attempts to clean up the state machine to acknolwedge child value assembly finish.
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
			{{- if .Type.ValueIsNullable }}
			tz := &ma.w.t[len(ma.w.t)-1]
			switch tz.v.m {
			case schema.Maybe_Null:
				ma.state = maState_initial
				ma.va.reset()
				return true
			case schema.Maybe_Value:
				{{- if (MaybeUsesPtr .Type.ValueType) }}
				tz.v.v = ma.va.w
				{{- end}}
				ma.state = maState_initial
				ma.va.reset()
				return true
			{{- else}}
			switch ma.cm {
			case schema.Maybe_Value:
				ma.va.w = nil
				ma.cm = schema.Maybe_Absent
				ma.state = maState_initial
				ma.va.reset()
				return true
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
	// REVIEW: there's a copy-by-value of k2 that's avoidable.  But it simplifies the error path.  Worth working on?
	// REVIEW: processing the key via the reprStyle of the key if it's type kind isn't string is currently supported, but should it be?  or is that more confusing than valuable?
	//  Very possible that it shouldn't be supported: the full-on keyAssembler route won't accept this, so consistency with that might be best.
	//  On the other hand, lookups by string *do* support this kind of processing (and it must, or PathSegment utility becomes unacceptably damaged), so either way, something feels surprising.
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

			var k2 _{{ .Type.KeyType | TypeSymbol }}
			{{- if eq .Type.KeyType.Kind.String "String" }}
			if err := (_{{ .Type.KeyType | TypeSymbol }}__Style{}).fromString(&k2, k); err != nil {
				return nil, err // TODO wrap in some kind of ErrInvalidKey
			}
			{{- else}}
			if err := (_{{ .Type.KeyType | TypeSymbol }}__ReprStyle{}).fromString(&k2, k); err != nil {
				return nil, err // TODO wrap in some kind of ErrInvalidKey
			}
			{{- end}}
			if _, exists := ma.w.m[k2]; exists {
				return nil, ipld.ErrRepeatedMapKey{&k2}
			}
			ma.w.t = append(ma.w.t, _{{ .Type | TypeSymbol }}__entry{k: k2})
			tz := &ma.w.t[len(ma.w.t)-1]
			ma.state = maState_midValue

			ma.w.m[k2] = &tz.v
			{{- if .Type.ValueIsNullable }}
			{{- if not (MaybeUsesPtr .Type.ValueType) }}
			ma.va.w = &tz.v.v
			{{- end}}
			ma.va.m = &tz.v.m
			tz.v.m = allowNull
			{{- else}}
			ma.va.w = &tz.v
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
			ma.w.t = append(ma.w.t, _{{ .Type | TypeSymbol }}__entry{})
			ma.state = maState_midKey
			ma.ka.m = &ma.cm
			ma.ka.w = &ma.w.t[len(ma.w.t)-1].k
			return &ma.ka
		}
		func (ma *_{{ .Type | TypeSymbol }}__Assembler) AssembleValue() ipld.NodeAssembler {
			switch ma.state {
			case maState_initial:
				panic("invalid state: AssembleValue cannot be called when no key is primed")
			case maState_midKey:
				if !ma.keyFinishTidy() {
					panic("invalid state: AssembleValue cannot be called when in the middle of assembling a key")
				} // if tidy success: carry on
			case maState_expectValue:
				// carry on
			case maState_midValue:
				panic("invalid state: AssembleValue cannot be called when in the middle of assembling another value")
			case maState_finished:
				panic("invalid state: AssembleValue cannot be called on an assembler that's already finished")
			}
			ma.state = maState_midValue
			return &ma.va
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
