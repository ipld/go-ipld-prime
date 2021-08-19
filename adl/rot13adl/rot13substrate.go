package rot13adl

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/mixins"
	"github.com/ipld/go-ipld-prime/schema"
)

// Substrate returns the root node of the raw internal data form of the ADL's content.
func (n *_R13String) Substrate() datamodel.Node {
	// This is a very minor twist in the case of the rot13 ADL.
	//  However, for larger ADLs (especially those relating to multi-block collections),
	//   this could be quite a bit more involved, and would almost certainly be the root node of a larger tree.
	return (*_Substrate)(n)
}

// -- Node -->

var _ datamodel.Node = (*_Substrate)(nil)

// Somewhat unusually for an ADL, there's only one substrate node type,
// and we actually made it have the same in-memory structure as the synthesized view node.
//
// When implementing other more complex ADLs, it will probably be more common to have
// the synthesized high-level node type either embed or have a pointer to the substrate root node.
type _Substrate _R13String

// REVIEW: what on earth we think the "TypeName" strings in error messages and other references to this node should be.
//  At the moment, it shares a prefix with the synthesized node, which is potentially confusing (?),
//  and I'm not sure what, if any, suffix actually makes meaningful sense to a user either.
//  I added the segment ".internal." to the middle of the name mangle; does this seem helpful?

func (*_Substrate) Kind() datamodel.Kind {
	return datamodel.Kind_String
}
func (*_Substrate) LookupByString(string) (datamodel.Node, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.LookupByString("")
}
func (*_Substrate) LookupByNode(datamodel.Node) (datamodel.Node, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.LookupByNode(nil)
}
func (*_Substrate) LookupByIndex(idx int64) (datamodel.Node, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.LookupByIndex(0)
}
func (*_Substrate) LookupBySegment(seg datamodel.PathSegment) (datamodel.Node, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.LookupBySegment(seg)
}
func (*_Substrate) MapIterator() datamodel.MapIterator {
	return nil
}
func (*_Substrate) ListIterator() datamodel.ListIterator {
	return nil
}
func (*_Substrate) Length() int64 {
	return -1
}
func (*_Substrate) IsAbsent() bool {
	return false
}
func (*_Substrate) IsNull() bool {
	return false
}
func (*_Substrate) AsBool() (bool, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.AsBool()
}
func (*_Substrate) AsInt() (int64, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.AsInt()
}
func (*_Substrate) AsFloat() (float64, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.AsFloat()
}
func (n *_Substrate) AsString() (string, error) {
	return n.raw, nil
}
func (*_Substrate) AsBytes() ([]byte, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.AsBytes()
}
func (*_Substrate) AsLink() (datamodel.Link, error) {
	return mixins.String{TypeName: "rot13adl.internal.Substrate"}.AsLink()
}
func (*_Substrate) Prototype() datamodel.NodePrototype {
	return _Substrate__Prototype{}
}

// -- NodePrototype -->

var _ datamodel.NodePrototype = _Substrate__Prototype{}

type _Substrate__Prototype struct {
	// There's no configuration to this ADL.
}

func (np _Substrate__Prototype) NewBuilder() datamodel.NodeBuilder {
	return &_Substrate__Builder{}
}

// -- NodeBuilder -->

var _ datamodel.NodeBuilder = (*_Substrate__Builder)(nil)

type _Substrate__Builder struct {
	_Substrate__Assembler
}

func (nb *_Substrate__Builder) Build() datamodel.Node {
	if nb.m != schema.Maybe_Value {
		panic("invalid state: cannot call Build on an assembler that's not finished")
	}
	return nb.w
}
func (nb *_Substrate__Builder) Reset() {
	*nb = _Substrate__Builder{}
}

// -- NodeAssembler -->

var _ datamodel.NodeAssembler = (*_Substrate__Assembler)(nil)

type _Substrate__Assembler struct {
	w *_Substrate
	m schema.Maybe // REVIEW: if the package where this Maybe enum lives is maybe not the right home for it after all.  Or should this line use something different?  We're only using some of its values after all.
}

func (_Substrate__Assembler) BeginMap(sizeHint int64) (datamodel.MapAssembler, error) {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.BeginMap(0)
}
func (_Substrate__Assembler) BeginList(sizeHint int64) (datamodel.ListAssembler, error) {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.BeginList(0)
}
func (na *_Substrate__Assembler) AssignNull() error {
	// REVIEW: unclear how this might compose with some other context (like a schema) which does allow nulls.  Probably a wrapper type?
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignNull()
}
func (_Substrate__Assembler) AssignBool(bool) error {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignBool(false)
}
func (_Substrate__Assembler) AssignInt(int64) error {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignInt(0)
}
func (_Substrate__Assembler) AssignFloat(float64) error {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignFloat(0)
}
func (na *_Substrate__Assembler) AssignString(v string) error {
	switch na.m {
	case schema.Maybe_Value:
		panic("invalid state: cannot assign into assembler that's already finished")
	}
	na.w = &_Substrate{
		raw:         v,
		synthesized: unrotate(v),
	}
	na.m = schema.Maybe_Value
	return nil
}
func (_Substrate__Assembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignBytes(nil)
}
func (_Substrate__Assembler) AssignLink(datamodel.Link) error {
	return mixins.StringAssembler{TypeName: "rot13adl.internal.Substrate"}.AssignLink(nil)
}
func (na *_Substrate__Assembler) AssignNode(v datamodel.Node) error {
	if v.IsNull() {
		return na.AssignNull()
	}
	if v2, ok := v.(*_Substrate); ok {
		switch na.m {
		case schema.Maybe_Value:
			panic("invalid state: cannot assign into assembler that's already finished")
		}
		na.w = v2
		na.m = schema.Maybe_Value
		return nil
	}
	if v2, err := v.AsString(); err != nil {
		return err
	} else {
		return na.AssignString(v2)
	}
}
func (_Substrate__Assembler) Prototype() datamodel.NodePrototype {
	return _Substrate__Prototype{}
}
