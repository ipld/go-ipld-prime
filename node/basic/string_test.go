package basicnode_test

import (
	"testing"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, basicnode.Prototype__String{})
}
