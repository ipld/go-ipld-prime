package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func BenchmarkSpec_Walk_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_Map3StrInt(b, Style__Any{})
}

func BenchmarkSpec_Walk_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Walk_MapNStrMap3StrInt(b, Style__Any{})
}
