package basicnode_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, basicnode.Prototype__String{})
}
