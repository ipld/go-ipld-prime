package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, Prototype__String{})
}
