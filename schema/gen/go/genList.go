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

func (listGenerator) IsRepr() bool { return false } // hint used in some generalized templates.

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
	// FUTURE: come back to this -- surely something nice can be done here.
}

func (g listGenerator) EmitNativeBuilder(w io.Writer) {
	// FUTURE: come back to this -- not yet clear what exactly might be most worth emitting here.
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

func (g listGenerator) EmitNodeMethodLookupByIndex(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupByIndex(idx int) (ipld.Node, error) {
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

func (g listGenerator) EmitNodeMethodLookupByNode(w io.Writer) {
	// LookupByNode will procede by coercing to int if it can; or fail; those are really the only options.
	// REVIEW: how much coercion is done by other types varies quite wildly.  so we should figure out if that inconsistency is acceptable, and at least document it if so.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) LookupByNode(k ipld.Node) (ipld.Node, error) {
			idx, err := k.AsInt()
			if err != nil {
				return nil, err
			}
			return n.LookupByIndex(idx)
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

func (g listGenerator) EmitNodeMethodPrototype(w io.Writer) {
	emitNodeMethodPrototype_typical(w, g.AdjCfg, g)
}

func (g listGenerator) EmitNodePrototypeType(w io.Writer) {
	emitNodePrototypeType_typical(w, g.AdjCfg, g)
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
	emitNodeAssemblerMethodBeginList_listoid(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	emitNodeAssemblerMethodAssignNull_recursive(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	emitNodeAssemblerMethodAssignNode_listoid(w, g.AdjCfg, g)
}
func (g listBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	emitNodeAssemblerHelper_listoid_tidyHelper(w, g.AdjCfg, g)
	emitNodeAssemblerHelper_listoid_listAssemblerMethods(w, g.AdjCfg, g)
}
