package ipldfree

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/tests"
)

func TestNodeBuilder(t *testing.T) {
	tests.TestBuildingScalars(t, NodeBuilder())
	tests.TestBuildingRecursives(t, NodeBuilder())
}

func TestMutableNode(t *testing.T) {
	// this should likely become legacy stuff and go away
	mutNodeFac := func() ipld.MutableNode { return &Node{} }
	tests.TestScalars(t, mutNodeFac)
	tests.TestRecursives(t, mutNodeFac)
}

func TestTokening(t *testing.T) {
	mutNodeFac := func() ipld.MutableNode { return &Node{} }
	tests.TestScalarMarshal(t, mutNodeFac)
	tests.TestRecursiveMarshal(t, mutNodeFac)
	tests.TestScalarUnmarshal(t, Unmarshal)
	tests.TestRecursiveUnmarshal(t, Unmarshal)
}
