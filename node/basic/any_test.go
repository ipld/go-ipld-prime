package basicnode_test

import (
	"testing"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestAnyBeingString(t *testing.T) {
	tests.SpecTestString(t, basicnode.Prototype__Any{})
}

func TestAnyBeingMapStrInt(t *testing.T) {
	tests.SpecTestMapStrInt(t, basicnode.Prototype__Any{})
}

func TestAnyBeingMapStrMapStrInt(t *testing.T) {
	tests.SpecTestMapStrMapStrInt(t, basicnode.Prototype__Any{})
}
