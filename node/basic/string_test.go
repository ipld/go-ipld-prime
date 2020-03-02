package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, Style__String{})
}
