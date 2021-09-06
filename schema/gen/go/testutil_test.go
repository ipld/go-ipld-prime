package gengo

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

// This file is full of helper functions.  Most are moderately embarassing.
//
// We should probably turn half of this into Wish Checkers;
//  they'd probably be much less fragile and give better error messages that way.
//  On the other hand, the functions for condensing two-arg returns wouldn't go away anyway.

func plz(n datamodel.Node, e error) datamodel.Node {
	if e != nil {
		panic(e)
	}
	return n
}
func plzStr(n datamodel.Node, e error) string {
	if e != nil {
		panic(e)
	}
	if s, ok := n.AsString(); ok == nil {
		return s
	} else {
		panic(ok)
	}
}
func str(n datamodel.Node) string {
	if s, ok := n.AsString(); ok == nil {
		return s
	} else {
		panic(ok)
	}
}
func erp(n datamodel.Node, e error) interface{} {
	if e != nil {
		return e
	}
	return n
}

// purely to syntactically flip large inline closures so we can see the argument at the top rather than the bottom of the block.
func withNode(n datamodel.Node, cb func(n datamodel.Node)) {
	cb(n)
}
