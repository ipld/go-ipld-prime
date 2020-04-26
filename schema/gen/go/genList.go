package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type listGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.ListTraits
	PkgName string
	Type    schema.TypeList
}

// --- native content and specializations --->

func (g listGenerator) EmitNativeType(w io.Writer) {
	// Lists are a pretty straightforward struct enclosing a slice.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct {
			x []_{{ .Type.ValueType | TypeSymbol }}{{if .Type.ValueIsNullable }}__Maybe{{end}}
		}
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g listGenerator) EmitNativeAccessors(w io.Writer) {
	// TODO: come back to this
}

func (g listGenerator) EmitNativeBuilder(w io.Writer) {
	// Not yet clear what exactly might be most worth emitting here.
}

func (g listGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
}

// --- type info --->

func (g listGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g listGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g listGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g listGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g listGenerator) EmitNodeTypeAssertions(w io.Writer) {
	emitNodeTypeAssertions_typical(w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodeMethodLookupIndex(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupIndex(idx int) (ipld.Node, error) {
			if n.Length() <= idx {
				return nil, ipld.ErrNotExists{ipld.PathSegmentOfInt(idx)}
			}
			v := &n.x[idx]
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

func (g listGenerator) EmitNodeMethodLookup(w io.Writer) {
	// LookupNode will procede by coercing to int if it can; or fail; those are really the only options.
	// REVIEW: how much coercion is done by other types varies quite wildly.  so we should figure out if that inconsistency is acceptable, and at least document it if so.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Lookup(k ipld.Node) (ipld.Node, error) {
			idx, err := k.AsInt()
			if err != nil {
				return nil, err
			}
			return n.LookupIndex(idx)
		}
	`, w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodeMethodListIterator(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) ListIterator() ipld.ListIterator {
			return &_{{ .Type | TypeSymbol }}__ListItr{n, 0}
		}

		type _{{ .Type | TypeSymbol }}__ListItr struct {
			n {{ .Type | TypeSymbol }}
			idx  int
		}

		func (itr *_{{ .Type | TypeSymbol }}__ListItr) Next() (idx int, v ipld.Node, _ error) {
			if itr.idx >= len(itr.n.x) {
				return -1, nil, ipld.ErrIteratorOverread{}
			}
			idx = itr.idx
			x := &itr.n.x[itr.idx]
			{{- if .Type.ValueIsNullable }}
			switch x.m {
			case schema.Maybe_Null:
				v = ipld.Null
			case schema.Maybe_Value:
				v = {{ if not (MaybeUsesPtr .Type.ValueType) }}&{{end}}x.v
			}
			{{- else}}
			v = x
			{{- end}}
			itr.idx++
			return
		}
		func (itr *_{{ .Type | TypeSymbol }}__ListItr) Done() bool {
			return itr.idx >= len(itr.n.x)
		}

	`, w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodeMethodLength(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Length() int {
			return len(n.x)
		}
	`, w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodeMethodStyle(w io.Writer) {
	emitNodeMethodStyle_typical(w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodeStyleType(w io.Writer) {
	emitNodeStyleType_typical(w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g listGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return listBuilderGenerator{
		g.AdjCfg,
		mixins.ListAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type listBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.ListAssemblerTraits
	PkgName string
	Type    schema.TypeList
}

func (listBuilderGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

func (g listBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	// - 'w' is the "**w**ip" pointer.
	// - 'm' is the **m**aybe which communicates our completeness to the parent if we're a child assembler.
	// - 'state' is what it says on the tin.  this is used for the list state (the broad transitions between null, start-list, and finish are handled by 'm' for consistency with other types).
	//
	// - 'cm' is **c**hild **m**aybe and is used for the completion message from children.
	//    It's only present if list values *aren't* allowed to be nullable, since otherwise they have their own per-value maybe slot we can use.
	// - 'va' is the embedded child value assembler.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
			state laState

			{{ if not .Type.ValueIsNullable }}cm schema.Maybe{{end}}
			va _{{ .Type.ValueType | TypeSymbol }}__Assembler
		}

		func (na *_{{ .Type | TypeSymbol }}__Assembler) reset() {
			na.state = laState_initial
			na.va.reset()
		}
	`, w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerMethodBeginList(w io.Writer) {
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
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
			if sizeHint > 0 {
				na.w.x = make([]_{{ .Type.ValueType | TypeSymbol }}{{if .Type.ValueIsNullable }}__Maybe{{end}}, 0, sizeHint)
			}
			return na, nil
		}
	`, w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	emitNodeAssemblerMethodAssignNull_recursive(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// AssignNode goes through three phases:
	// 1. is it null?  Jump over to AssignNull (which may or may not reject it).
	// 2. is it our own type?  Handle specially -- we might be able to do efficient things.
	// 3. is it the right kind to morph into us?  Do so.
	//
	// We do not set m=midvalue in phase 3 -- it shouldn't matter unless you're trying to pull off concurrent access, which is wrong and unsafe regardless.
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
			if v.ReprKind() != ipld.ReprKind_List {
				return ipld.ErrWrongKind{TypeName: "{{ .PkgName }}.{{ .Type.Name }}", MethodName: "AssignNode", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: v.ReprKind()}
			}
			itr := v.ListIterator()
			for !itr.Done() {
				_, v, err := itr.Next()
				if err != nil {
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

func (g listBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	g.emitListAssemblerValueTidyHelper(w)
	g.emitListAssemblerMethods(w)
}
func (g listBuilderGenerator) emitListAssemblerValueTidyHelper(w io.Writer) {
	// This function attempts to clean up the state machine to acknolwedge child value assembly finish.
	//  If the child was finished and we just collected it, return true and update state to laState_initial.
	//  Otherwise, if it wasn't done, return false;
	//   and the caller is almost certain to emit an error momentarily.
	// The function will only be called when the current state is laState_midValue.
	//  (In general, the idea is that if the user is doing things correctly,
	//   this function will only be called when the child is in fact finished.)
	// If 'cm' is used, we reset it to its initial condition of Maybe_Absent here.
	//  At the same time, we nil the 'w' pointer for the child assembler; otherwise its own state machine would probably let it modify 'w' again!
	doTemplate(`
		func (la *_{{ .Type | TypeSymbol }}__Assembler) valueFinishTidy() bool {
			{{- if .Type.ValueIsNullable }}
			row := &la.w.x[len(la.w.x)-1]
			switch row.m {
			case schema.Maybe_Value:
				{{- if (MaybeUsesPtr .Type.ValueType) }}
				row.v = la.va.w
				{{- end}}
				fallthrough
			case schema.Maybe_Null:
				la.state = laState_initial
				la.va.reset()
				return true
			{{- else}}
			switch la.cm {
			case schema.Maybe_Value:
				la.va.w = nil
				la.cm = schema.Maybe_Absent
				la.state = laState_initial
				la.va.reset()
				return true
			{{- end}}
			default:
				return false
			}
		}
	`, w, g.AdjCfg, g)
}
func (g listBuilderGenerator) emitListAssemblerMethods(w io.Writer) {
	doTemplate(`
		func (la *_{{ .Type | TypeSymbol }}__Assembler) AssembleValue() ipld.NodeAssembler {
			switch la.state {
			case laState_initial:
				// carry on
			case laState_midValue:
				if !la.valueFinishTidy() {
					panic("invalid state: AssembleValue cannot be called when still in the middle of assembling the previous value")
				} // if tidy success: carry on
			case laState_finished:
				panic("invalid state: AssembleValue cannot be called on an assembler that's already finished")
			}
			la.w.x = append(la.w.x, _{{ .Type.ValueType | TypeSymbol }}{{if .Type.ValueIsNullable }}__Maybe{{end}}{})
			la.state = laState_midValue
			row := &la.w.x[len(la.w.x)-1]
			{{- if .Type.ValueIsNullable }}
			{{- if not (MaybeUsesPtr .Type.ValueType) }}
			la.va.w = &row.v
			{{- end}}
			la.va.m = &row.m
			row.m = allowNull
			{{- else}}
			la.va.w = row
			la.va.m = &la.cm
			{{- end}}
			return &la.va
		}
	`, w, g.AdjCfg, g)
	doTemplate(`
		func (la *_{{ .Type | TypeSymbol }}__Assembler) Finish() error {
			switch la.state {
			case laState_initial:
				// carry on
			case laState_midValue:
				if !la.valueFinishTidy() {
					panic("invalid state: Finish cannot be called when in the middle of assembling a value")
				} // if tidy success: carry on
			case laState_finished:
				panic("invalid state: Finish cannot be called on an assembler that's already finished")
			}
			la.state = laState_finished
			*la.m = schema.Maybe_Value
			return nil
		}
	`, w, g.AdjCfg, g)
	doTemplate(`
		func (la *_{{ .Type | TypeSymbol }}__Assembler) ValueStyle(_ int) ipld.NodeStyle {
			return _{{ .Type.ValueType | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}
