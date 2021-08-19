// Package 'must' provides another alternative to the 'fluent' package,
// providing many helpful functions for wrapping methods with multiple returns
// into a single return (converting errors into panics).
//
// It's useful especially for testing code and other situations where panics
// are not problematic.
//
// Unlike the 'fluent' package, panics are not of any particular type.
// There is no equivalent to the `fluent.Recover` feature in the 'must' package.
//
// Because golang supports implied destructuring of multiple-return functions
// into arguments for another funtion of matching arity, most of the 'must'
// functions can used smoothly in a pointfree/chainable form, like this:
//
//		must.Node(SomeNodeBuilder{}.CreateString("a"))
//
package must

import (
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/schema"
)

// must.NotError simply panics if given an error.
// It helps turn multi-line code into one-liner code in situations where
// you simply don't care.
func NotError(e error) {
	if e != nil {
		panic(e)
	}
}

// must.Node helps write pointfree/chainable-style code
// by taking a Node and an error and transforming any error into a panic.
//
// Because golang supports implied destructuring of multiple-return functions
// into arguments for another funtion of matching arity, it can be used like this:
//
//		must.Node(SomeNodeBuilder{}.CreateString("a"))
//
func Node(n datamodel.Node, e error) datamodel.Node {
	if e != nil {
		panic(e)
	}
	return n
}

// must.TypedNode helps write pointfree/chainable-style code
// by taking a Node and an error and transforming any error into a panic.
// It will also cast the `datamodel.Node` to a `schema.TypedNode`, panicking if impossible.
//
// Because golang supports implied destructuring of multiple-return functions
// into arguments for another funtion of matching arity, it can be used like this:
//
//		must.TypedNode(SomeNodeBuilder{}.CreateString("a"))
//
func TypedNode(n datamodel.Node, e error) schema.TypedNode {
	if e != nil {
		panic(e)
	}
	return n.(schema.TypedNode)
}

// must.True panics if the given bool is false.
func True(v bool) {
	if !v {
		panic("must.True")
	}
}

// must.String unboxes the given Node via AsString,
// panicking in the case that the Node isn't of string kind,
// and otherwise returning the bare native string.
func String(n datamodel.Node) string {
	if v, err := n.AsString(); err != nil {
		panic(err)
	} else {
		return v
	}
}

// must.Int unboxes the given Node via AsInt,
// panicking in the case that the Node isn't of int kind,
// and otherwise returning the bare native int.
func Int(n datamodel.Node) int64 {
	if v, err := n.AsInt(); err != nil {
		panic(err)
	} else {
		return v
	}
}
