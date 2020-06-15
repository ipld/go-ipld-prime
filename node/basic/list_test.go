package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestList(t *testing.T) {
	tests.SpecTestListString(t, Style__List{})
}
