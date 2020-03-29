package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type stringGenerator struct {
	mixins.StringTraits
	PkgName string
	Type    schema.TypeString
	Symbol  string // defaults to Type.Name but can be overriden.
}

// --- native content and specializations --->

func (g stringGenerator) EmitNativeType(w io.Writer) {
	// Using a struct with a single member is the same size in memory as a typedef,
	//  while also having the advantage of meaning we can block direct casting,
	//   which is desirable because the compiler then ensures our validate methods can't be evaded.
	doTemplate(`
		type _{{ .Symbol }} struct{ x string }
		type {{ .Symbol }} = *_{{ .Symbol }}
	`, w, g)
}

func (g stringGenerator) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsString` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (n {{ .Symbol }}) String() string {
			return n.x
		}
	`, w, g)
}

func (g stringGenerator) EmitNativeBuilder(w io.Writer) {
	// Scalar types are easy to generate a constructor function for.
	// REVIEW: if this is useful and should be on by default; it also adds a decent amount of noise to a package.
	// FUTURE: should engage validation flow.
	doTemplate(`
		func New{{ .Symbol }}(v string) {{ .Symbol }} {
			n := _{{ .Symbol }}{v}
			return &n
		}
	`, w, g)
}

func (g stringGenerator) EmitNativeMaybe(w io.Writer) {
	// REVIEW: can this be extracted to the mixins package?  it doesn't even vary for kind.
	// REVIEW: what conventions and interfaces are required around Maybe types is very non-finalized.
	doTemplate(`
		type Maybe{{ .Symbol }} struct {
			m schema.Maybe
			n {{ .Symbol }}
		}

		func (m Maybe{{ .Symbol }}) IsNull() bool {
			return m.m == schema.Maybe_Null
		}
		func (m Maybe{{ .Symbol }}) IsUndefined() bool {
			return m.m == schema.Maybe_Absent
		}
		func (m Maybe{{ .Symbol }}) Exists() bool {
			return m.m == schema.Maybe_Value
		}
		func (m Maybe{{ .Symbol }}) Must() *{{ .Symbol }} {
			if !m.Exists() {
				panic("unbox of a maybe rejected")
			}
			return &m.n
		}
	`, w, g)
}

// --- type info --->

func (g stringGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g)
}

// --- TypedNode interface satisfaction --->

func (g stringGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Symbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g)
}

func (g stringGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	// Perhaps surprisingly, the way to get the representation node pointer
	//  does not actually depend on what the representation strategy is.
	doTemplate(`
		func (n {{ .Symbol }}) Representation() ipld.Node {
			return (*_{{ .Symbol }}__Repr)(n)
		}
	`, w, g)
}

// --- Node interface satisfaction --->

func (g stringGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g stringGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Symbol }})(&_{{ .Symbol }}{})
		var _ schema.TypedNode = ({{ .Symbol }})(&_{{ .Symbol }}{})
	`, w, g)
}

func (g stringGenerator) EmitNodeMethodAsString(w io.Writer) {
	doTemplate(`
		func (n {{ .Symbol }}) AsString() (string, error) {
			return n.x, nil
		}
	`, w, g)
}

func (g stringGenerator) EmitNodeMethodStyle(w io.Writer) {
	doTemplate(`
		func ({{ .Symbol }}) Style() ipld.NodeStyle {
			return nil // TODO
		}
	`, w, g)
}

func (g stringGenerator) EmitNodeStyleType(w io.Writer) {
	doTemplate(`
		type _{{ .Symbol }}__Style struct{}

		func (_{{ .Symbol }}__Style) NewBuilder() ipld.NodeBuilder {
			var nb _{{ .Symbol }}__Builder
			nb.Reset()
			return &nb
		}
	`, w, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g stringGenerator) EmitNodeBuilder(w io.Writer) {
	doTemplate(`
		type _{{ .Symbol }}__Builder struct {
			_{{ .Symbol }}__Assembler
		}

		func (nb *_{{ .Symbol }}__Builder) Build() ipld.Node {
			return nb.w
		}
		func (nb *_{{ .Symbol }}__Builder) Reset() {
			var w _{{ .Symbol }}
			*nb = _{{ .Symbol }}__Builder{_{{ .Symbol }}__Assembler{w: &w}}
		}
	`, w, g)
}

func (g stringGenerator) EmitNodeAssembler(w io.Writer) {
	doTemplate(`
		type _{{ .Symbol }}__Assembler struct {
			w *_{{ .Symbol }}
		}

		func (_{{ .Symbol }}__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.BeginMap(0)
		}
		func (_{{ .Symbol }}__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.BeginList(0)
		}
		func (_{{ .Symbol }}__Assembler) AssignNull() error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignNull()
		}
		func (_{{ .Symbol }}__Assembler) AssignBool(bool) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignBool(false)
		}
		func (_{{ .Symbol }}__Assembler) AssignInt(int) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignInt(0)
		}
		func (_{{ .Symbol }}__Assembler) AssignFloat(float64) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignFloat(0)
		}
		func (na *_{{ .Symbol }}__Assembler) AssignString(v string) error {
			*na.w = _{{ .Symbol }}{v}
			return nil
		}
		func (_{{ .Symbol }}__Assembler) AssignBytes([]byte) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignBytes(nil)
		}
		func (_{{ .Symbol }}__Assembler) AssignLink(ipld.Link) error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .Type.Name }}"}.AssignLink(nil)
		}
		func (na *_{{ .Symbol }}__Assembler) AssignNode(v ipld.Node) error {
			if v2, err := v.AsString(); err != nil {
				return err
			} else {
				*na.w = _{{ .Symbol }}{v2}
				return nil
			}
		}
		func (_{{ .Symbol }}__Assembler) Style() ipld.NodeStyle {
			return _{{ .Symbol }}__Style{}
		}
	`, w, g)
}
