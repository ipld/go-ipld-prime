package basicnode_test

import (
	"testing"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestList(t *testing.T) {
	tests.SpecTestListString(t, basicnode.Prototype__List{})
}
