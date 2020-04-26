package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type intGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.IntTraits
	PkgName string
	Type    schema.TypeInt
}

// --- native content and specializations --->

func (g intGenerator) EmitNativeType(w io.Writer) {
	// Using a struct with a single member is the same size in memory as a typedef,
	//  while also having the advantage of meaning we can block direct casting,
	//   which is desirable because the compiler then ensures our validate methods can't be evaded.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct{ x int }
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsInt` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Int() int {
			return n.x
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeBuilder(w io.Writer) {
	// Generate a single-step construction function -- this is easy to do for a scalar,
	//  and all representations of scalar kind can be expected to have a method like this.
	// The function is attached to the nodestyle for convenient namespacing;
	//  it needs no new memory, so it would be inappropriate to attach to the builder or assembler.
	// FUTURE: should engage validation flow.
	doTemplate(`
		func (_{{ .Type | TypeSymbol }}__Style) FromInt(v int) ({{ .Type | TypeSymbol }}, error) {
			n := _{{ .Type | TypeSymbol }}{v}
			return &n, nil
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeMaybe(w io.Writer) {
	emitNativeMaybe(w, g.AdjCfg, g)
}

// --- type info --->

func (g intGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g intGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	emitTypicalTypedNodeMethodRepresentation(w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g intGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g intGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
		var _ schema.TypedNode = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) AsInt() (int, error) {
			return n.x, nil
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeMethodStyle(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeStyleType(w io.Writer) {
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

func (g intGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return intBuilderGenerator{
		g.AdjCfg,
		mixins.IntAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type intBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.IntAssemblerTraits
	PkgName string
	Type    schema.TypeInt
}

func (g intBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	emitEmitNodeBuilderType_typical(w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	emitNodeBuilderMethods_typical(w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	emitNodeAssemblerType_scalar(w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	emitNodeAssemblerMethodAssignNull_scalar(w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignInt(v int) error {
			switch *na.m {
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			}
			{{- if .Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ .Type | TypeSymbol }}{}
			}
			{{- end}}
			na.w.x = v
			*na.m = schema.Maybe_Value
			return nil
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	// AssignNode goes through three phases:
	// 1. is it null?  Jump over to AssignNull (which may or may not reject it).
	// 2. is it our own type?  Handle specially -- we might be able to do efficient things.
	// 3. is it the right kind to morph into us?  Do so.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			if v.IsNull() {
				return na.AssignNull()
			}
			if v2, ok := v.(*_{{ .Type | TypeSymbol }}); ok {
				switch *na.m {
				case schema.Maybe_Value, schema.Maybe_Null:
					panic("invalid state: cannot assign into assembler that's already finished")
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
			if v2, err := v.AsInt(); err != nil {
				return err
			} else {
				return na.AssignInt(v2)
			}
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// Nothing needed here for int kinds.
}
