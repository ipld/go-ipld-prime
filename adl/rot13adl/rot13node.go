/*
	rot13adl is a demo ADL -- its purpose is to show what an ADL and its public interface can look like.
	It implements a "rot13" string: when creating data through the ADL, the user gives it a regular string;
	the ADL will create aninternal representation of it which has the characters altered in a reversable way.

	It provides reference and example materal, but it's very unlikely you want to use it in real situations ;)

	There are several ways to move data in and out of the ADL:

		- treat it like a regular IPLD map:
			- using the exported NodePrototype can be used to get a NodeBuilder which can accept keys and values;
			- using the resulting Node and doing lookup operations on it like a regular map;
		- load up raw substrate data and `Reify()` it into the synthesized form, and *then* treat it like a regular map:
			- this is handy if the raw data already parsed into Nodes.
			- optionally, use `SubstrateRootPrototype` as the prototype for loading the raw substrate data;
			  any kind of Node is a valid input to Reify, but this one will generally have optimal performance.
		- take the synthesized form and inspect its substrate data:
			- the `Substrate()` method will return another ipld.Node which is the root of the raw substrate data,
			  and can be walked normally like any other ipld.Node.
*/
package rot13adl

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
	"github.com/ipld/go-ipld-prime/schema"
)

// -- Node -->

var _ ipld.Node = (*_R13String)(nil)

type _R13String struct {
	raw         string // the raw content, before our ADL lens is applied to it.
	synthesized string // the content that the ADL presents.  calculated proactively from the original, in this implementation (though you could imagine implementing it lazily, in either direction, too).
}

func (*_R13String) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (*_R13String) LookupByString(string) (ipld.Node, error) {
	return mixins.String{"rot13adl.R13String"}.LookupByString("")
}
func (*_R13String) LookupByNode(ipld.Node) (ipld.Node, error) {
	return mixins.String{"rot13adl.R13String"}.LookupByNode(nil)
}
func (*_R13String) LookupByIndex(idx int) (ipld.Node, error) {
	return mixins.String{"rot13adl.R13String"}.LookupByIndex(0)
}
func (*_R13String) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.String{"rot13adl.R13String"}.LookupBySegment(seg)
}
func (*_R13String) MapIterator() ipld.MapIterator {
	return nil
}
func (*_R13String) ListIterator() ipld.ListIterator {
	return nil
}
func (*_R13String) Length() int {
	return -1
}
func (*_R13String) IsAbsent() bool {
	return false
}
func (*_R13String) IsNull() bool {
	return false
}
func (*_R13String) AsBool() (bool, error) {
	return mixins.String{"rot13adl.R13String"}.AsBool()
}
func (*_R13String) AsInt() (int, error) {
	return mixins.String{"rot13adl.R13String"}.AsInt()
}
func (*_R13String) AsFloat() (float64, error) {
	return mixins.String{"rot13adl.R13String"}.AsFloat()
}
func (n *_R13String) AsString() (string, error) {
	return n.synthesized, nil
}
func (*_R13String) AsBytes() ([]byte, error) {
	return mixins.String{"rot13adl.R13String"}.AsBytes()
}
func (*_R13String) AsLink() (ipld.Link, error) {
	return mixins.String{"rot13adl.R13String"}.AsLink()
}
func (*_R13String) Prototype() ipld.NodePrototype {
	return _R13String__Prototype{}
}

// -- NodePrototype -->

var _ ipld.NodePrototype = _R13String__Prototype{}

type _R13String__Prototype struct {
	// There's no configuration to this ADL.

	// A more complex ADL might have some kind of parameters here.
	//
	// The general contract of a NodePrototype is supposed to be that:
	// when you get one from an existing Node,
	//  it should have enough information to create a new Node that
	//   could "replace" the previous one in whatever context it's in.
	// For ADLs, that means it should carry most of the configuration.
	//
	// An ADL that does multi-block stuff might also need functions like a LinkLoader passed in through here.
}

func (np _R13String__Prototype) NewBuilder() ipld.NodeBuilder {
	return &_R13String__Builder{}
}

// -- NodeBuilder -->

var _ ipld.NodeBuilder = (*_R13String__Builder)(nil)

type _R13String__Builder struct {
	_R13String__Assembler
}

func (nb *_R13String__Builder) Build() ipld.Node {
	if nb.m != schema.Maybe_Value {
		panic("invalid state: cannot call Build on an assembler that's not finished")
	}
	return nb.w
}
func (nb *_R13String__Builder) Reset() {
	*nb = _R13String__Builder{}
}

// -- NodeAssembler -->

var _ ipld.NodeAssembler = (*_R13String__Assembler)(nil)

type _R13String__Assembler struct {
	w *_R13String
	m schema.Maybe // REVIEW: if the package where this Maybe enum lives is maybe not the right home for it after all.  Or should this line use something different?  We're only using some of its values after all.
}

func (_R13String__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	return mixins.StringAssembler{"rot13adl.R13String"}.BeginMap(0)
}
func (_R13String__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	return mixins.StringAssembler{"rot13adl.R13String"}.BeginList(0)
}
func (na *_R13String__Assembler) AssignNull() error {
	// REVIEW: unclear how this might compose with some other context (like a schema) which does allow nulls.  Probably a wrapper type?
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignNull()
}
func (_R13String__Assembler) AssignBool(bool) error {
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignBool(false)
}
func (_R13String__Assembler) AssignInt(int) error {
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignInt(0)
}
func (_R13String__Assembler) AssignFloat(float64) error {
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignFloat(0)
}
func (na *_R13String__Assembler) AssignString(v string) error {
	switch na.m {
	case schema.Maybe_Value:
		panic("invalid state: cannot assign into assembler that's already finished")
	}
	na.w = &_R13String{
		raw:         rotate(v),
		synthesized: v,
	}
	na.m = schema.Maybe_Value
	return nil
}
func (_R13String__Assembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignBytes(nil)
}
func (_R13String__Assembler) AssignLink(ipld.Link) error {
	return mixins.StringAssembler{"rot13adl.R13String"}.AssignLink(nil)
}
func (na *_R13String__Assembler) AssignNode(v ipld.Node) error {
	if v.IsNull() {
		return na.AssignNull()
	}
	if v2, ok := v.(*_R13String); ok {
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
func (_R13String__Assembler) Prototype() ipld.NodePrototype {
	return _R13String__Prototype{}
}
