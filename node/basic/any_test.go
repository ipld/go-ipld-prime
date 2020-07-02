package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestAnyBeingString(t *testing.T) {
	tests.SpecTestString(t, Prototype__Any{})
}

func TestAnyBeingMapStrInt(t *testing.T) {
	tests.SpecTestMapStrInt(t, Prototype__Any{})
}

func TestAnyBeingMapStrMapStrInt(t *testing.T) {
	tests.SpecTestMapStrMapStrInt(t, Prototype__Any{})
}
