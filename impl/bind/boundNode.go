package ipldbind

import (
	"fmt"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/polydawn/refmt/obj/atlas"
)

var (
	_ ipld.Node = Node{}
)

/*
	Node binds to some Go object in memory, using the definitions provided
	by refmt's object atlasing tools.

	This binding does not provide a serialization valid for hashing; to
	compute a CID, you'll have to convert to another kind of node.
	If you're not sure which kind serializable node to use, try `ipldcbor.Node`.
*/
type Node struct {
	kind  ipld.ReprKind // compute during bind
	bound reflect.Value // should already be ptr-unwrapped
	atlas atlas.Atlas
}

/*
	Bind binds any go value into being interfacable as a Node, using the provided
	atlas to understand how to traverse it.
*/
func Bind(bindme interface{}, atl atlas.Atlas) ipld.Node {
	rv := reflect.ValueOf(bindme)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return Node{
		kind:  determineReprKind(rv),
		bound: rv,
		atlas: atl,
	}
}

func determineReprKind(rv reflect.Value) ipld.ReprKind {
	switch rv.Type().Kind() {
	case reflect.Bool:
		return ipld.ReprKind_Bool
	case reflect.String:
		return ipld.ReprKind_String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return ipld.ReprKind_Int
	case reflect.Float32, reflect.Float64:
		return ipld.ReprKind_Float
	case reflect.Complex64, reflect.Complex128:
		panic(fmt.Errorf("invalid ipldbind.Node: ipld has no concept for complex numbers"))
	case reflect.Array, reflect.Slice:
		return ipld.ReprKind_List
	case reflect.Map, reflect.Struct:
		return ipld.ReprKind_Map
	case reflect.Interface:
		determineReprKind(rv.Elem())
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Errorf("invalid ipldbind.Node: cannot atlas over type %q; it's a %v", rv.Type().Name(), rv.Kind()))
	case reflect.Ptr:
		// might've already been traversed during bind, but interface path can find more.
		determineReprKind(rv.Elem())
	}
	panic("unreachable")
}

func (n Node) Kind() ipld.ReprKind {
	return n.kind
}

func (n Node) AsBool() (v bool, _ error) {
	reflect.ValueOf(&v).Elem().Set(n.bound)
	return
}
func (n Node) AsString() (v string, _ error) {
	reflect.ValueOf(&v).Elem().Set(n.bound)
	return
}
func (n Node) AsInt() (v int, _ error) {
	reflect.ValueOf(&v).Elem().Set(n.bound)
	return
}
func (n Node) AsLink() (v cid.Cid, _ error) {
	reflect.ValueOf(&v).Elem().Set(n.bound)
	return
}

func (n Node) Keys() ([]string, int) {
	return nil, 0 // FIXME
	// TODO: REVIEW: structs have clear key order; maps do not.  what do?
}

func (n Node) TraverseField(pth string) (ipld.Node, error) {
	v := n.bound

	// Traverse.
	//  Honor any atlent overrides if present;
	//  Use kind-based fallbacks if necessary.
	atlent, exists := n.atlas.Get(reflect.ValueOf(v.Type()).Pointer())
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
			return Node{}, fmt.Errorf("traverse failed: type %q has no field named %q", v.Type().Name(), pth)
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
			return Node{determineReprKind(v), v, n.atlas}, nil
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
	case // primitives: wrap in node
		reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return Node{determineReprKind(v), v, n.atlas}, nil
	case // recursives: wrap in node
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Struct,
		reflect.Interface:
		return Node{determineReprKind(v), v, n.atlas}, nil
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

func (n Node) TraverseIndex(idx int) (ipld.Node, error) {
	panic("NYI")
}
