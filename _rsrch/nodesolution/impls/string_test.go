package impls

import (
	"testing"

	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/impls/tests"
)

func TestString(t *testing.T) {
	tests.SpecTestString(t, Style__String{})
}
