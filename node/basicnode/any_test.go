package basicnode_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestAnyBeingString(t *testing.T) {
	tests.SpecTestString(t, basicnode.Prototype.Any)
}

func TestAnyBeingMapStrInt(t *testing.T) {
	tests.SpecTestMapStrInt(t, basicnode.Prototype.Any)
}

func TestAnyBeingMapStrMapStrInt(t *testing.T) {
	tests.SpecTestMapStrMapStrInt(t, basicnode.Prototype.Any)
}
