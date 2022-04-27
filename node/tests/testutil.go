package tests

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

// This file is full of helper functions.  Most are moderately embarassing.
//
// We should probably turn half of this into Wish Checkers;
//  they'd probably be much less fragile and give better error messages that way.
//  On the other hand, the functions for condensing two-arg returns wouldn't go away anyway.

// various benchmarks assign their final result here,
// in order to defuse the possibility of their work being elided.
var sink interface{} //lint:ignore U1000 used by benchmarks

// purely to syntactically flip large inline closures so we can see the argument at the top rather than the bottom of the block.
func withNode(n datamodel.Node, cb func(n datamodel.Node)) {
	cb(n)
}
