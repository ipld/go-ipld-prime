package realgen

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
	"github.com/ipld/go-ipld-prime/schema"
)

// Code generated go-ipld-prime DO NOT EDIT.

type _String struct{ x string }
type String = *_String

func (n String) String() string {
	return n.x
}
func NewString(v string) String {
	n := _String{v}
	return &n
}

type _String__Maybe struct {
	m schema.Maybe
	v String
}
type MaybeString = *_String__Maybe

func (m MaybeString) IsNull() bool {
	return m.m == schema.Maybe_Null
}
func (m MaybeString) IsUndefined() bool {
	return m.m == schema.Maybe_Absent
}
func (m MaybeString) Exists() bool {
	return m.m == schema.Maybe_Value
}
func (m MaybeString) AsNode() ipld.Node {
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
func (m MaybeString) Must() String {
	if !m.Exists() {
		panic("unbox of a maybe rejected")
	}
	return m.v
}

var _ ipld.Node = (String)(&_String{})
var _ schema.TypedNode = (String)(&_String{})

func (String) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_String
}
func (String) LookupString(string) (ipld.Node, error) {
	return mixins.String{"realgen.String"}.LookupString("")
}
func (String) Lookup(ipld.Node) (ipld.Node, error) {
	return mixins.String{"realgen.String"}.Lookup(nil)
}
func (String) LookupIndex(idx int) (ipld.Node, error) {
	return mixins.String{"realgen.String"}.LookupIndex(0)
}
func (String) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.String{"realgen.String"}.LookupSegment(seg)
}
func (String) MapIterator() ipld.MapIterator {
	return nil
}
func (String) ListIterator() ipld.ListIterator {
	return nil
}
func (String) Length() int {
	return -1
}
func (String) IsUndefined() bool {
	return false
}
func (String) IsNull() bool {
	return false
}
func (String) AsBool() (bool, error) {
	return mixins.String{"realgen.String"}.AsBool()
}
func (String) AsInt() (int, error) {
	return mixins.String{"realgen.String"}.AsInt()
}
func (String) AsFloat() (float64, error) {
	return mixins.String{"realgen.String"}.AsFloat()
}
func (n String) AsString() (string, error) {
	return n.x, nil
}
func (String) AsBytes() ([]byte, error) {
	return mixins.String{"realgen.String"}.AsBytes()
}
func (String) AsLink() (ipld.Link, error) {
	return mixins.String{"realgen.String"}.AsLink()
}
func (String) Style() ipld.NodeStyle {
	return _String__Style{}
}

type _String__Style struct{}

func (_String__Style) NewBuilder() ipld.NodeBuilder {
	var nb _String__Builder
	nb.Reset()
	return &nb
}

type _String__Builder struct {
	_String__Assembler
}

func (nb *_String__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *_String__Builder) Reset() {
	var w _String
	*nb = _String__Builder{_String__Assembler{w: &w}}
}
func (nb *_String__Builder) AssignNull() error {
	return mixins.StringAssembler{"realgen.String"}.AssignNull()
}
func (nb *_String__Builder) AssignString(v string) error {
	*nb.w = _String{v}
	return nil
}
func (nb *_String__Builder) AssignNode(v ipld.Node) error {
	if v2, err := v.AsString(); err != nil {
		return err
	} else {
		return nb.AssignString(v2)
	}
}

type _String__Assembler struct {
	w   *_String
	z   bool
	fcb func() error
}

func (_String__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	return mixins.StringAssembler{"realgen.String"}.BeginMap(0)
}
func (_String__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	return mixins.StringAssembler{"realgen.String"}.BeginList(0)
}
func (na *_String__Assembler) AssignNull() error {
	na.z = true
	return na.fcb()
}
func (_String__Assembler) AssignBool(bool) error {
	return mixins.StringAssembler{"realgen.String"}.AssignBool(false)
}
func (_String__Assembler) AssignInt(int) error {
	return mixins.StringAssembler{"realgen.String"}.AssignInt(0)
}
func (_String__Assembler) AssignFloat(float64) error {
	return mixins.StringAssembler{"realgen.String"}.AssignFloat(0)
}
func (na *_String__Assembler) AssignString(v string) error {
	if na.w == nil {
		na.w = &_String{v}
		return na.fcb()
	}
	*na.w = _String{v}
	return na.fcb()
}
func (_String__Assembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{"realgen.String"}.AssignBytes(nil)
}
func (_String__Assembler) AssignLink(ipld.Link) error {
	return mixins.StringAssembler{"realgen.String"}.AssignLink(nil)
}
func (na *_String__Assembler) AssignNode(v ipld.Node) error {
	if v.IsNull() {
		return na.AssignNull()
	}
	if v2, err := v.AsString(); err != nil {
		return err
	} else {
		return na.AssignString(v2)
	}
}
func (_String__Assembler) Style() ipld.NodeStyle {
	return _String__Style{}
}
func (String) Type() schema.Type {
	return nil /*TODO:typelit*/
}
func (n String) Representation() ipld.Node {
	return (*_String__Repr)(n)
}

type _String__Repr = _String

var _ ipld.Node = &_String__Repr{}

type _String__ReprStyle = _String__Style

func (_String__ReprStyle) construct(w *_String, v string) error {
	*w = _String{v}
	return nil
}

type _String__ReprAssembler = _String__Assembler
