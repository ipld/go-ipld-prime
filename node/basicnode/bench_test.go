package basicnode_test

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_Map3StrInt(b, basicnode.Prototype.Any)
}

func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_MapNStrMap3StrInt(b, basicnode.Prototype.Any)
}
