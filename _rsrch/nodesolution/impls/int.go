package impls

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

var (
	_ ipld.Node          = plainInt(0)
	_ ipld.NodeStyle     = Style__Int{}
	_ ipld.NodeBuilder   = &plainInt__Builder{}
	_ ipld.NodeAssembler = &plainInt__Assembler{}
)

func Int(value int) ipld.Node {
	return plainInt(value)
}

// plainInt is a simple boxed int that complies with ipld.Node.
type plainInt int

// -- Node interface methods -->

func (plainInt) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (plainInt) LookupString(string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupString", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) Lookup(key ipld.Node) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "Lookup", AppropriateKind: ipld.ReprKindSet_JustMap, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{MethodName: "LookupSegment", AppropriateKind: ipld.ReprKindSet_Recursive, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) MapIterator() ipld.MapIterator {
	return nil
}
func (plainInt) ListIterator() ipld.ListIterator {
	return nil
}
func (plainInt) Length() int {
	return -1
}
func (plainInt) IsUndefined() bool {
	return false
}
func (plainInt) IsNull() bool {
	return false
}
func (plainInt) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Int}
}
func (n plainInt) AsInt() (int, error) {
	return int(n), nil
}
func (plainInt) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsString() (string, error) {
	return "", ipld.ErrWrongKind{MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Int}
}
func (plainInt) Style() ipld.NodeStyle {
	panic("todo")
}

// -- NodeStyle -->

type Style__Int struct{}

func (Style__Int) NewBuilder() ipld.NodeBuilder {
	var w plainInt
	return &plainInt__Builder{plainInt__Assembler{w: &w}}
}

// -- NodeBuilder -->

type plainInt__Builder struct {
	plainInt__Assembler
}

func (nb *plainInt__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *plainInt__Builder) Reset() {
	var w plainInt
	*nb = plainInt__Builder{plainInt__Assembler{w: &w}}
}

// -- NodeAssembler -->

type plainInt__Assembler struct {
	w *plainInt
}

func (plainInt__Assembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (plainInt__Assembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainInt__Assembler) AssignNull() error                                      { panic("no") }
func (plainInt__Assembler) AssignBool(bool) error                                  { panic("no") }
func (na *plainInt__Assembler) AssignInt(v int) error {
	*na.w = plainInt(v)
	return nil
}
func (plainInt__Assembler) AssignFloat(float64) error { panic("no") }
func (plainInt__Assembler) AssignString(string) error { panic("no") }
func (plainInt__Assembler) AssignBytes([]byte) error  { panic("no") }
func (na *plainInt__Assembler) Assign(v ipld.Node) error {
	if s, err := v.AsInt(); err != nil {
		return err
	} else {
		*na.w = plainInt(s)
		return nil
	}
}
func (plainInt__Assembler) Style() ipld.NodeStyle {
	return Style__Int{}
}
