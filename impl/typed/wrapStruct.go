package typed

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	"github.com/ipld/go-ipld-prime/schema"
)

var _ Node = wrapnodeStruct{}

type wrapnodeStruct struct {
	ipld.Node
	typ schema.TypeStruct
}

// Most of the 'nope' methods from the inner node are fine;
// we add the extra things required for typed.Node;
// we decorate the getters and iterators to handle the distinct path around optionals
// and return a different error for missing fields;
// length becomes fixed to a constant;
// and we replace the builder with a complete wrapper that maintains type rules.

// (We could override more of the Node methods to return errors with accurate type name, though.)

func (tn wrapnodeStruct) Type() schema.Type {
	return tn.typ
}

func (tn wrapnodeStruct) LookupString(key string) (ipld.Node, error) {
	for _, field := range tn.typ.Fields() {
		if field.Name() != key {
			continue
		}
		v, e1 := tn.Node.LookupString(key)
		if e1 == nil {
			return v, nil // null or set both flow through here
		}
		if _, ok := e1.(ipld.ErrNotExists); ok {
			return ipld.Undef, nil // we assume the type allows this, or this node shouldn't have been possible to construct in the first place
		}
		return nil, e1
	}
	return nil, ErrNoSuchField{Type: tn.typ, FieldName: key}
}

func (tn wrapnodeStruct) MapIterator() ipld.MapIterator {
	return &wrapnodeStruct_Iterator{&tn, 0}
}

type wrapnodeStruct_Iterator struct {
	node *wrapnodeStruct
	idx  int
}

func (itr *wrapnodeStruct_Iterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.idx >= itr.node.Length() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	field := itr.node.typ.Fields()[itr.idx]
	k = ipldfree.String(field.Name())
	v, e1 := itr.node.LookupString(field.Name())
	if e1 != nil {
		if _, ok := e1.(ipld.ErrNotExists); ok {
			v = ipld.Undef // we assume the type allows this, or this node shouldn't have been possible to construct in the first place
		} else {
			return k, nil, e1
		}
	}
	itr.idx++
	return
}
func (itr *wrapnodeStruct_Iterator) Done() bool {
	return itr.idx >= itr.node.Length()
}

func (tn wrapnodeStruct) Length() int {
	return len(tn.typ.Fields())
}

func (tn wrapnodeStruct) NodeBuilder() ipld.NodeBuilder {
	panic("todo")
}

func (tn wrapnodeStruct) Representation() ipld.Node {
	switch rs := tn.typ.RepresentationStrategy().(type) {
	case schema.StructRepresentation_Map:
		panic("todo") // TODO: add new source file for each of these.
	case schema.StructRepresentation_Tuple:
		panic("todo") // TODO: add new source file for each of these.
	case schema.StructRepresentation_StringPairs:
		panic("todo") // TODO: add new source file for each of these.
	case schema.StructRepresentation_StringJoin:
		panic("todo") // TODO: add new source file for each of these.
	default:
		_ = rs
		panic("unreachable (schema.StructRepresentation sum type)")
	}
}

// The builder is a more complete straightjacket; it wouldn't be correct to
// assume that the builder we're delegating internal storage to would reject
// other kinds (e.g. CreateString) entirely, and our type requires that.

type wrapnodeStruct_Builder struct {
	utnb ipld.NodeBuilder
	typ  schema.TypeStruct
}

func (nb wrapnodeStruct_Builder) CreateMap() (ipld.MapBuilder, error) {
	mb, err := nb.utnb.CreateMap()
	if err != nil {
		return nil, err
	}
	needs := make(map[string]struct{}, len(nb.typ.Fields()))
	for _, field := range nb.typ.Fields() {
		if !field.IsOptional() {
			needs[field.Name()] = struct{}{}
		}
	}
	return &wrapnodeStruct_MapBuilder{mb, nb.typ, needs}, nil
}
func (nb wrapnodeStruct_Builder) AmendMap() (ipld.MapBuilder, error) {
	panic("TODO") // TODO
}
func (nb wrapnodeStruct_Builder) CreateList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) AmendList() (ipld.ListBuilder, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.AmendList", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateNull() (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateNull", AppropriateKind: ipld.ReprKindSet_JustNull, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateBool(v bool) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateInt(v int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateInt", AppropriateKind: ipld.ReprKindSet_JustInt, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateFloat(v float64) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateString(v string) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateBytes(v []byte) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (nb wrapnodeStruct_Builder) CreateLink(v ipld.Link) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: string(nb.typ.Name()), MethodName: "NodeBuilder.CreateLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}

type wrapnodeStruct_MapBuilder struct {
	utmb  ipld.MapBuilder
	typ   schema.TypeStruct
	needs map[string]struct{}
	// We have to remember if anything was intentionally unset so we can check at the end.
	// Or, initialize with the set of things that need to be set, and decrement it; easier.
}

func (mb *wrapnodeStruct_MapBuilder) Insert(k, v ipld.Node) error {
	ks, err := k.AsString()
	if err != nil {
		return err
	}
	// Check that the field exists at all.
	field := mb.typ.Field(ks)
	if field == nil {
		return ErrNoSuchField{Type: mb.typ, FieldName: ks}
	}
	// Check that the value is assignable to this field, or return error.
	vt, ok := v.(Node)
	switch {
	case v.IsNull():
		if !field.IsNullable() {
			return fmt.Errorf("type mismatch on struct field assignment: cannot assign null to non-nullable field")
		}
		// if null and nullable: carry on.
	case ok:
		if mb.typ.Field(ks).Type() != vt.Type() {
			return fmt.Errorf("type mismatch on struct field assignment")
		}
		// if typed node, and it matches: carry on.
	default:
		return fmt.Errorf("need typed.Node for insertion into struct") // FUTURE: maybe if it's a basic enough thing we sholud attempt coerce?
	}
	// Insert the value, and note it's now been set.
	if err := mb.utmb.Insert(k, v); err != nil {
		return err
	}
	delete(mb.needs, ks)
	return nil
}
func (mb *wrapnodeStruct_MapBuilder) Delete(k ipld.Node) error {
	panic("delete not supported on this type") // I have serious questions about whether the delete method deserves to exist.
}
func (mb *wrapnodeStruct_MapBuilder) Build() (ipld.Node, error) {
	if len(mb.needs) > 0 {
		return nil, fmt.Errorf("missing required fields") // TODO say which
	}
	n, err := mb.Build()
	if err != nil {
		return nil, err
	}
	return wrapnodeStruct{n, mb.typ}, nil
}

// TODO and soon the nb methods for getting child builders.
// also those will have fun methods for handling the ability to have undefined..?
//     no they shouldn't actually that's important -- those features really only occur in Insert methods.
