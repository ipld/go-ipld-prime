package impls

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/impls/mixins"
)

var (
	_ ipld.Node          = plainString("")
	_ ipld.NodeStyle     = Style__String{}
	_ ipld.NodeBuilder   = &plainString__Builder{}
	_ ipld.NodeAssembler = &plainString__Assembler{}
)

func String(value string) ipld.Node {
	return plainString(value)
}

// plainString is a simple boxed string that complies with ipld.Node.
// It's useful for many things, such as boxing map keys.
//
// The implementation is a simple typedef of a string;
// handling it as a Node incurs 'runtime.convTstring',
// which is about the best we can do.
type plainString string

// -- Node interface methods -->

func (plainString) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (plainString) LookupString(string) (ipld.Node, error) {
	return mixins.String{"string"}.LookupString("")
}
func (plainString) Lookup(key ipld.Node) (ipld.Node, error) {
	return mixins.String{"string"}.Lookup(nil)
}
func (plainString) LookupIndex(idx int) (ipld.Node, error) {
	return mixins.String{"string"}.LookupIndex(0)
}
func (plainString) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.String{"string"}.LookupSegment(seg)
}
func (plainString) MapIterator() ipld.MapIterator {
	return nil
}
func (plainString) ListIterator() ipld.ListIterator {
	return nil
}
func (plainString) Length() int {
	return -1
}
func (plainString) IsUndefined() bool {
	return false
}
func (plainString) IsNull() bool {
	return false
}
func (plainString) AsBool() (bool, error) {
	return mixins.String{"string"}.AsBool()
}
func (plainString) AsInt() (int, error) {
	return mixins.String{"string"}.AsInt()
}
func (plainString) AsFloat() (float64, error) {
	return mixins.String{"string"}.AsFloat()
}
func (x plainString) AsString() (string, error) {
	return string(x), nil
}
func (plainString) AsBytes() ([]byte, error) {
	return mixins.String{"string"}.AsBytes()
}
func (plainString) AsLink() (ipld.Link, error) {
	return mixins.String{"string"}.AsLink()
}
func (plainString) Style() ipld.NodeStyle {
	return Style__String{}
}

// -- NodeStyle -->

type Style__String struct{}

func (Style__String) NewBuilder() ipld.NodeBuilder {
	var w plainString
	return &plainString__Builder{plainString__Assembler{w: &w}}
}

// -- NodeBuilder -->

type plainString__Builder struct {
	plainString__Assembler
}

func (nb *plainString__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *plainString__Builder) Reset() {
	var w plainString
	*nb = plainString__Builder{plainString__Assembler{w: &w}}
}

// -- NodeAssembler -->

type plainString__Assembler struct {
	w *plainString
}

func (plainString__Assembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (plainString__Assembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainString__Assembler) AssignNull() error                                      { panic("no") }
func (plainString__Assembler) AssignBool(bool) error                                  { panic("no") }
func (plainString__Assembler) AssignInt(int) error                                    { panic("no") }
func (plainString__Assembler) AssignFloat(float64) error                              { panic("no") }
func (na *plainString__Assembler) AssignString(v string) error {
	*na.w = plainString(v)
	return nil
}
func (plainString__Assembler) AssignBytes([]byte) error { panic("no") }
func (na *plainString__Assembler) AssignNode(v ipld.Node) error {
	if s, err := v.AsString(); err != nil {
		return err
	} else {
		*na.w = plainString(s)
		return nil
	}
}
func (plainString__Assembler) Style() ipld.NodeStyle {
	return Style__String{}
}
