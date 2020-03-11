package gendemo

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

var (
	_ ipld.Node          = plainString("")
	_ ipld.NodeStyle     = Style__String{}
	_ ipld.NodeBuilder   = &plainString__Builder{}
	_ ipld.NodeAssembler = &plainString__Assembler{}
)

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
	return mixins.String{"gendemo.String"}.LookupString("")
}
func (plainString) Lookup(key ipld.Node) (ipld.Node, error) {
	return mixins.String{"gendemo.String"}.Lookup(nil)
}
func (plainString) LookupIndex(idx int) (ipld.Node, error) {
	return mixins.String{"gendemo.String"}.LookupIndex(0)
}
func (plainString) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.String{"gendemo.String"}.LookupSegment(seg)
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
	return mixins.String{"gendemo.String"}.AsBool()
}
func (plainString) AsInt() (int, error) {
	return mixins.String{"gendemo.String"}.AsInt()
}
func (plainString) AsFloat() (float64, error) {
	return mixins.String{"gendemo.String"}.AsFloat()
}
func (x plainString) AsString() (string, error) {
	return string(x), nil
}
func (plainString) AsBytes() ([]byte, error) {
	return mixins.String{"gendemo.String"}.AsBytes()
}
func (plainString) AsLink() (ipld.Link, error) {
	return mixins.String{"gendemo.String"}.AsLink()
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

func (plainString__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	return mixins.StringAssembler{"gendemo.String"}.BeginMap(0)
}
func (plainString__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	return mixins.StringAssembler{"gendemo.String"}.BeginList(0)
}
func (plainString__Assembler) AssignNull() error {
	return mixins.StringAssembler{"gendemo.String"}.AssignNull()
}
func (plainString__Assembler) AssignBool(bool) error {
	return mixins.StringAssembler{"gendemo.String"}.AssignBool(false)
}
func (plainString__Assembler) AssignInt(int) error {
	return mixins.StringAssembler{"gendemo.String"}.AssignInt(0)
}
func (plainString__Assembler) AssignFloat(float64) error {
	return mixins.StringAssembler{"gendemo.String"}.AssignFloat(0)
}
func (na *plainString__Assembler) AssignString(v string) error {
	*na.w = plainString(v)
	return nil
}
func (plainString__Assembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{"gendemo.String"}.AssignBytes(nil)
}
func (plainString__Assembler) AssignLink(ipld.Link) error {
	return mixins.StringAssembler{"gendemo.String"}.AssignLink(nil)
}
func (na *plainString__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsString(); err != nil {
		return err
	} else {
		*na.w = plainString(v2)
		return nil
	}
}
func (plainString__Assembler) Style() ipld.NodeStyle {
	return Style__String{}
}
