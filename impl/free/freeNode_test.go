package ipldfree

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/tests"
)

func Test(t *testing.T) {
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
