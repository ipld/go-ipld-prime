package ipldbind

import (
	"fmt"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/polydawn/refmt/obj/atlas"
)

var (
	_ ipld.Node = &Node{}
)

/*
	Node binds to some Go object in memory, using the definitions provided
	by refmt's object atlasing tools.

	This binding does not provide a serialization valid for hashing; to
	compute a CID, you'll have to convert to another kind of node.
	If you're not sure which kind serializable node to use, try `ipldcbor.Node`.
*/
type Node struct {
	bound reflect.Value
	atlas atlas.Atlas
}

/*
	Bind binds any go value into being interfacable as a Node, using the provided
	atlas to understand how to traverse it.
*/
func Bind(bindme interface{}, atl atlas.Atlas) ipld.Node {
	return &Node{
		bound: reflect.ValueOf(bindme),
		atlas: atl,
	}
}

func (n *Node) GetField(pth string) (v interface{}, _ error) {
	return v, traverseField(n.bound, pth, n.atlas, reflect.ValueOf(v))
}
func (n *Node) GetFieldBool(pth string) (v bool, _ error) {
	return v, traverseField(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetFieldString(pth string) (v string, _ error) {
	return v, traverseField(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetFieldInt(pth string) (v int, _ error) {
	return v, traverseField(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetFieldLink(pth string) (v cid.Cid, _ error) {
	return v, traverseField(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}

func (n *Node) GetIndex(pth int) (v interface{}, _ error) {
	return v, traverseIndex(n.bound, pth, n.atlas, reflect.ValueOf(v))
}
func (n *Node) GetIndexBool(pth int) (v bool, _ error) {
	return v, traverseIndex(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetIndexString(pth int) (v string, _ error) {
	return v, traverseIndex(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetIndexInt(pth int) (v int, _ error) {
	return v, traverseIndex(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}
func (n *Node) GetIndexLink(pth int) (v cid.Cid, _ error) {
	return v, traverseIndex(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}

func traverseField(v reflect.Value, pth string, atl atlas.Atlas, assignTo reflect.Value) error {
	// Unwrap any pointers.
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Traverse.
	//  Honor any atlent overrides if present;
	//  Use kind-based fallbacks if necessary.
	atlent, exists := atl.Get(reflect.ValueOf(v.Type()).Pointer())
	if exists {
		switch {
		case atlent.MarshalTransformFunc != nil:
			panic(fmt.Errorf("invalid ipldbind.Node: type %q atlas specifies transform, but ipld doesn't support this power level", v.Type().Name()))
		case atlent.StructMap != nil:
			for _, fe := range atlent.StructMap.Fields {
				if fe.SerialName == pth {
					v = fe.ReflectRoute.TraverseToValue(v)
					break
				}
			}
			return fmt.Errorf("traverse failed: type %q has no field named %q", v.Type().Name(), pth)
		case atlent.UnionKeyedMorphism != nil:
			panic(fmt.Errorf("invalid ipldbind.Node: type %q atlas specifies union, but ipld doesn't know how to make sense of this", v.Type().Name()))
		case atlent.MapMorphism != nil:
			v = v.MapIndex(reflect.ValueOf(pth))
		default:
			panic("unreachable")
		}
	} else {
		switch v.Type().Kind() {
		case // primitives: set them
			reflect.Bool,
			reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128:
			panic(fmt.Errorf("invalid ipldbind.Node: atlas for type %q is union; ipld doesn't know how to make sense of this", v.Type().Name()))
		case // recursives: wrap in node
			reflect.Array,
			reflect.Slice,
			reflect.Map,
			reflect.Struct,
			reflect.Interface:
			assignTo.Set(reflect.ValueOf(Node{v, atl}))
		case // esotera: reject with panic
			reflect.Chan,
			reflect.Func,
			reflect.UnsafePointer:
			panic(fmt.Errorf("invalid ipldbind.Node: cannot atlas over type %q; it's a %v", v.Type().Name(), v.Kind()))
		case // pointers: should've already been unwrapped
			reflect.Ptr:
			panic("unreachable")
		}
	}

	// Unwrap any pointers.
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Assign into the result.
	//  Either assign the result directly (for primitives)
	//  Or wrap with a Node and assign that (for recursives).
	// TODO decide what to do with typedef'd primitives.
	switch v.Type().Kind() {
	case // primitives: set them
		reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		assignTo.Set(v)
		return nil
	case // recursives: wrap in node
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Struct,
		reflect.Interface:
		assignTo.Set(reflect.ValueOf(Node{v, atl}))
		return nil
	case // esotera: reject with panic
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		panic(fmt.Errorf("invalid ipldbind.Node: cannot atlas over type %q; it's a %v", v.Type().Name(), v.Kind()))
	case // pointers: should've already been unwrapped
		reflect.Ptr:
		panic("unreachable")
	default:
		panic("unreachable")
	}
}

func traverseIndex(v reflect.Value, pth int, atl atlas.Atlas, assignTo reflect.Value) error {
	panic("NYI")
}
