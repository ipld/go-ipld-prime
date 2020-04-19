package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type stringGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.StringTraits
	PkgName string
	Type    schema.TypeString
}

// --- native content and specializations --->

func (g stringGenerator) EmitNativeType(w io.Writer) {
	// Using a struct with a single member is the same size in memory as a typedef,
	//  while also having the advantage of meaning we can block direct casting,
	//   which is desirable because the compiler then ensures our validate methods can't be evaded.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct{ x string }
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsString` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) String() string {
			return n.x
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNativeBuilder(w io.Writer) {
	// Generate a single-step construction function -- this is easy to do for a scalar,
	//  and all representations of scalar kind can be expected to have a method like this.
	// The function is attached to the nodestyle for convenient namespacing;
	//  it needs no new memory, so it would be inappropriate to attach to the builder or assembler.
	// The function is directly used internally by anything else that might involve recursive destructuring on the same scalar kind
	//  (for example, structs using stringjoin strategies that have one of this type as a field, etc).
	// FUTURE: should engage validation flow.
	doTemplate(`
		func (_{{ .Type | TypeSymbol }}__Style) fromString(w *_{{ .Type | TypeSymbol }}, v string) error {
			*w = _{{ .Type | TypeSymbol }}{v}
			return nil
		}
	`, w, g.AdjCfg, g)
	// And generate a publicly exported version of that single-step constructor, too.
	//  (Just don't expose the details about allocation, because you can't meaningfully use that from outside the package.)
	doTemplate(`
		func (_{{ .Type | TypeSymbol }}__Style) FromString(v string) ({{ .Type | TypeSymbol }}, error) {
			n := _{{ .Type | TypeSymbol }}{v}
			return &n, nil
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNativeMaybe(w io.Writer) {
	// REVIEW: can this be extracted to the mixins package?  it doesn't appear to vary for kind.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Maybe struct {
			m schema.Maybe
			v {{if not (MaybeUsesPtr .Type) }}_{{end}}{{ .Type | TypeSymbol }}
		}
		type Maybe{{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}__Maybe

		func (m Maybe{{ .Type | TypeSymbol }}) IsNull() bool {
			return m.m == schema.Maybe_Null
		}
		func (m Maybe{{ .Type | TypeSymbol }}) IsUndefined() bool {
			return m.m == schema.Maybe_Absent
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Exists() bool {
			return m.m == schema.Maybe_Value
		}
		func (m Maybe{{ .Type | TypeSymbol }}) AsNode() ipld.Node {
			switch m.m {
				case schema.Maybe_Absent:
					return ipld.Undef
				case schema.Maybe_Null:
					return ipld.Null
				case schema.Maybe_Value:
					return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
				default:
					panic("unreachable")
			}
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Must() {{ .Type | TypeSymbol }} {
			if !m.Exists() {
				panic("unbox of a maybe rejected")
			}
			return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
		}
	`, w, g.AdjCfg, g)
}

// --- type info --->

func (g stringGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g stringGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	// Perhaps surprisingly, the way to get the representation node pointer
	//  does not actually depend on what the representation strategy is.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Representation() ipld.Node {
			return (*_{{ .Type | TypeSymbol }}__Repr)(n)
		}
	`, w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g stringGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g stringGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
		var _ schema.TypedNode = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) AsString() (string, error) {
			return n.x, nil
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNodeMethodStyle(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNodeStyleType(w io.Writer) {
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

func (g stringGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return stringBuilderGenerator{
		g.AdjCfg,
		mixins.StringAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type stringBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.StringAssemblerTraits
	PkgName string
	Type    schema.TypeString
}

func (g stringBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}
	`, w, g.AdjCfg, g)
}
func (g stringBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			if *nb.m != schema.Maybe_Value {
				panic("invalid state: cannot call Build on an assembler that's not finished")
			}
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			var m schema.Maybe
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{&w, &m}}
		}
	`, w, g.AdjCfg, g)
}
func (g stringBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			m *schema.Maybe
		}
	`, w, g.AdjCfg, g)
}
func (g stringBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNull() error {
			switch *na.m {
			case allowNull:
				*na.m = schema.Maybe_Null
				return nil
			case schema.Maybe_Absent:
				return mixins.StringAssembler{"{{ .PkgName }}.{{ .TypeName }}"}.AssignNull()
			case schema.Maybe_Value, schema.Maybe_Null:
				panic("invalid state: cannot assign into assembler that's already finished")
			}
			panic("unreachable")
		}
	`, w, g.AdjCfg, g)
}
func (g stringBuilderGenerator) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignString(v string) error {
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
func (g stringBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
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
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				return na.AssignString(v2)
			}
		}
	`, w, g.AdjCfg, g)
}
func (g stringBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// Nothing needed here for string kinds.
}
