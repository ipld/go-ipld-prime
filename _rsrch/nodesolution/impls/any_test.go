package impls

import (
	"testing"

	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/impls/tests"
)

func TestAnyBeingString(t *testing.T) {
	tests.SpecTestString(t, Style__Any{})
}

func TestAnyBeingMapStrInt(t *testing.T) {
	CheckMapStrInt(t, Style__Any{})
}

func TestAnyBeingMapStrMapStrInt(t *testing.T) {
	CheckMapStrMapStrInt(t, Style__Any{})
}
