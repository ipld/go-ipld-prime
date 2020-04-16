package realgen

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
	"github.com/ipld/go-ipld-prime/schema"
)

// Code generated go-ipld-prime DO NOT EDIT.

type _Int struct{ x int }
type Int = *_Int

func (n Int) Int() int {
	return n.x
}
func NewInt(v int) Int {
	n := _Int{v}
	return &n
}

type _Int__Maybe struct {
	m schema.Maybe
	v Int
}
type MaybeInt = *_Int__Maybe

func (m MaybeInt) IsNull() bool {
	return m.m == schema.Maybe_Null
}
func (m MaybeInt) IsUndefined() bool {
	return m.m == schema.Maybe_Absent
}
func (m MaybeInt) Exists() bool {
	return m.m == schema.Maybe_Value
}
func (m MaybeInt) AsNode() ipld.Node {
	switch m.m {
	case schema.Maybe_Absent:
		return ipld.Undef
	case schema.Maybe_Null:
		return ipld.Null
	case schema.Maybe_Value:
		return m.v
	default:
		panic("unreachable")
	}
}
func (m MaybeInt) Must() Int {
	if !m.Exists() {
		panic("unbox of a maybe rejected")
	}
	return m.v
}

var _ ipld.Node = (Int)(&_Int{})
var _ schema.TypedNode = (Int)(&_Int{})

func (Int) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (Int) LookupString(string) (ipld.Node, error) {
	return mixins.Int{"realgen.Int"}.LookupString("")
}
func (Int) Lookup(ipld.Node) (ipld.Node, error) {
	return mixins.Int{"realgen.Int"}.Lookup(nil)
}
func (Int) LookupIndex(idx int) (ipld.Node, error) {
	return mixins.Int{"realgen.Int"}.LookupIndex(0)
}
func (Int) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Int{"realgen.Int"}.LookupSegment(seg)
}
func (Int) MapIterator() ipld.MapIterator {
	return nil
}
func (Int) ListIterator() ipld.ListIterator {
	return nil
}
func (Int) Length() int {
	return -1
}
func (Int) IsUndefined() bool {
	return false
}
func (Int) IsNull() bool {
	return false
}
func (Int) AsBool() (bool, error) {
	return mixins.Int{"realgen.Int"}.AsBool()
}
func (n Int) AsInt() (int, error) {
	return n.x, nil
}
func (Int) AsFloat() (float64, error) {
	return mixins.Int{"realgen.Int"}.AsFloat()
}
func (Int) AsString() (string, error) {
	return mixins.Int{"realgen.Int"}.AsString()
}
func (Int) AsBytes() ([]byte, error) {
	return mixins.Int{"realgen.Int"}.AsBytes()
}
func (Int) AsLink() (ipld.Link, error) {
	return mixins.Int{"realgen.Int"}.AsLink()
}
func (Int) Style() ipld.NodeStyle {
	return _Int__Style{}
}

type _Int__Style struct{}

func (_Int__Style) NewBuilder() ipld.NodeBuilder {
	var nb _Int__Builder
	nb.Reset()
	return &nb
}

type _Int__Builder struct {
	_Int__Assembler
}

func (nb *_Int__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *_Int__Builder) Reset() {
	var w _Int
	*nb = _Int__Builder{_Int__Assembler{w: &w}}
}
func (nb *_Int__Builder) AssignNull() error {
	return mixins.StringAssembler{"realgen.Int"}.AssignNull()
}
func (nb *_Int__Builder) AssignInt(v int) error {
	*nb.w = _Int{v}
	return nil
}
func (nb *_Int__Builder) AssignNode(v ipld.Node) error {
	if v2, err := v.AsInt(); err != nil {
		return err
	} else {
		return nb.AssignInt(v2)
	}
}

type _Int__Assembler struct {
	w   *_Int
	z   bool
	fcb func() error
}

func (_Int__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	return mixins.IntAssembler{"realgen.Int"}.BeginMap(0)
}
func (_Int__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	return mixins.IntAssembler{"realgen.Int"}.BeginList(0)
}
func (na *_Int__Assembler) AssignNull() error {
	na.z = true
	return na.fcb()
}
func (_Int__Assembler) AssignBool(bool) error {
	return mixins.IntAssembler{"realgen.Int"}.AssignBool(false)
}
func (na *_Int__Assembler) AssignInt(v int) error {
	if na.w == nil {
		na.w = &_Int{v}
		return na.fcb()
	}
	*na.w = _Int{v}
	return na.fcb()
}
func (_Int__Assembler) AssignFloat(float64) error {
	return mixins.IntAssembler{"realgen.Int"}.AssignFloat(0)
}
func (_Int__Assembler) AssignString(string) error {
	return mixins.IntAssembler{"realgen.Int"}.AssignString("")
}
func (_Int__Assembler) AssignBytes([]byte) error {
	return mixins.IntAssembler{"realgen.Int"}.AssignBytes(nil)
}
func (_Int__Assembler) AssignLink(ipld.Link) error {
	return mixins.IntAssembler{"realgen.Int"}.AssignLink(nil)
}
func (na *_Int__Assembler) AssignNode(v ipld.Node) error {
	if v.IsNull() {
		return na.AssignNull()
	}
	if v2, err := v.AsInt(); err != nil {
		return err
	} else {
		return na.AssignInt(v2)
	}
}
func (_Int__Assembler) Style() ipld.NodeStyle {
	return _Int__Style{}
}
func (Int) Type() schema.Type {
	return nil /*TODO:typelit*/
}
func (n Int) Representation() ipld.Node {
	return (*_Int__Repr)(n)
}

type _Int__Repr = _Int

var _ ipld.Node = &_Int__Repr{}

type _Int__ReprStyle = _Int__Style
type _Int__ReprAssembler = _Int__Assembler
