package basicnode_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestList(t *testing.T) {
	tests.SpecTestListString(t, basicnode.Prototype.List)
}
