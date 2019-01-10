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
	tests.TestScalarMarshal(t, mutNodeFac)
	tests.TestRecursiveMarshal(t, mutNodeFac)
}
