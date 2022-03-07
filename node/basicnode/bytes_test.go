package basicnode_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestBytes(t *testing.T) {
	tests.SpecTestBytes(t, basicnode.Prototype__Bytes{})
}
