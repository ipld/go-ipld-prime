package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/node/mixins/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, Style__String{})
}
