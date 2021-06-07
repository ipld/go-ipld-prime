package basicnode_test

import (
	"testing"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_Map3StrInt(b, basicnode.Prototype__Any{})
}

func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_MapNStrMap3StrInt(b, basicnode.Prototype__Any{})
}
