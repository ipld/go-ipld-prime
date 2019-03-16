package ipldfree

import (
	"testing"

	"github.com/ipld/go-ipld-prime/tests"
)

func TestNodeBuilder(t *testing.T) {
	tests.TestBuildingScalars(t, NodeBuilder())
	tests.TestBuildingRecursives(t, NodeBuilder())
}

func TestTokening(t *testing.T) {
	tests.TestScalarMarshal(t, NodeBuilder())
	tests.TestRecursiveMarshal(t, NodeBuilder())
	tests.TestScalarUnmarshal(t, NodeBuilder())
	tests.TestRecursiveUnmarshal(t, NodeBuilder())
}
