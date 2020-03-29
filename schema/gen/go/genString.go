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
	// Scalar types are easy to generate a constructor function for.
	// REVIEW: if this is useful and should be on by default; it also adds a decent amount of noise to a package.
	// FUTURE: should engage validation flow.
	doTemplate(`
		func New{{ .Type | TypeSymbol }}(v string) {{ .Type | TypeSymbol }} {
			n := _{{ .Type | TypeSymbol }}{v}
			return &n
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNativeMaybe(w io.Writer) {
	// REVIEW: can this be extracted to the mixins package?  it doesn't even vary for kind.
	// REVIEW: what conventions and interfaces are required around Maybe types is very non-finalized.
	doTemplate(`
		type Maybe{{ .Type | TypeSymbol }} struct {
			m schema.Maybe
			n {{ .Type | TypeSymbol }}
		}

		func (m Maybe{{ .Type | TypeSymbol }}) IsNull() bool {
			return m.m == schema.Maybe_Null
		}
		func (m Maybe{{ .Type | TypeSymbol }}) IsUndefined() bool {
			return m.m == schema.Maybe_Absent
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Exists() bool {
			return m.m == schema.Maybe_Value
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Must() {{ .Type | TypeSymbol }} {
			if !m.Exists() {
				panic("unbox of a maybe rejected")
			}
			return m.n
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
			return nil // TODO
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

func (g stringGenerator) EmitNodeBuilder(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}

		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{w: &w}}
		}
	`, w, g.AdjCfg, g)
}

func (g stringGenerator) EmitNodeAssembler(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
		}

		func (_{{ .Type | TypeSymbol }}__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.BeginMap(0)
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.BeginList(0)
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignNull() error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignNull()
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignBool(bool) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignBool(false)
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignInt(int) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignInt(0)
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignFloat(float64) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignFloat(0)
		}
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignString(v string) error {
			*na.w = _{{ .Type | TypeSymbol }}{v}
			return nil
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignBytes([]byte) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignBytes(nil)
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) AssignLink(ipld.Link) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignLink(nil)
		}
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				*na.w = _{{ .Type | TypeSymbol }}{v2}
				return nil
			}
		}
		func (_{{ .Type | TypeSymbol }}__Assembler) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}
