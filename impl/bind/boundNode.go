package ipldbind

import (
	"fmt"
	"reflect"

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

func (n *Node) GetField(pth []string) (v interface{}, _ error) {
	return v, traverse(n.bound, pth, n.atlas, reflect.ValueOf(v))
}
func (n *Node) GetFieldString(pth []string) (v string, _ error) {
	return v, traverse(n.bound, pth, n.atlas, reflect.ValueOf(&v).Elem())
}

// traverse is the internal impl behind GetField.
// It's preferable to have this function for recursing because the defn of
// the GetField function is more for caller-friendliness than performance.
func traverse(v reflect.Value, pth []string, atl atlas.Atlas, assignTo reflect.Value) error {
	// Handle the terminating case of expected leaf nodes first.
	if len(pth) == 0 {
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
		case reflect.Array: // array: ... REVIEW, like map, is it acceptable to leak concrete types?
		case reflect.Slice: // slice: same as array
		case reflect.Map: // map: ... REVIEW: can we just return this?  it might have more concrete types in it, and that's kind of unstandard.
		case reflect.Struct: // struct: wrap in Node
			// TODO
		case reflect.Interface: // interface: ... REVIEW: i don't know what to do with this
		case reflect.Chan: // chan: not acceptable in ipld objects
		case reflect.Func: // func: not acceptable in ipld objects
		case reflect.Ptr: // ptr: TODO needs an unwrap round
		case reflect.UnsafePointer: // unsafe: not acceptable in ipld objects
		}
	}
	// Handle traversal to deeper nodes.
	//  If we get a primitive here, it's an error, because we haven't handled all path segments yet.

	atlent, exists := atl.Get(reflect.ValueOf(v.Type()).Pointer())
	if !exists {
		panic(fmt.Errorf("invalid ipldbind.Node: atlas missing entry for type %q", v.Type().Name()))
	}
	errIfPathNonEmpty := func() error {
		if len(pth) > 1 {
			return fmt.Errorf("getField reached leaf before following all path segements")
		}
		return nil
	}
	// TODO all these cases
	switch atlent.Type.Kind() {
	case // primitives found when expecting to path deeper cause an error return.
		reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return errIfPathNonEmpty()
	case reflect.Array:
	case reflect.Slice:
	case reflect.Map:
	case reflect.Struct:
	case reflect.Interface:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Ptr:
	case reflect.UnsafePointer:
	}
	return nil
}
