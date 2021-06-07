package basicnode_test

import (
	"testing"

	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestMap(t *testing.T) {
	tests.SpecTestMapStrInt(t, basicnode.Prototype__Map{})
	tests.SpecTestMapStrMapStrInt(t, basicnode.Prototype__Map{})
	tests.SpecTestMapStrListStr(t, basicnode.Prototype__Map{})
}

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, basicnode.Prototype__Map{})
}
func BenchmarkMapStrInt_3n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleEntry(b, basicnode.Prototype__Map{})
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, basicnode.Prototype__Map{})
}

func BenchmarkMapStrInt_25n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleStandard(b, basicnode.Prototype__Map{})
}
func BenchmarkMapStrInt_25n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleEntry(b, basicnode.Prototype__Map{})
}
func BenchmarkMapStrInt_25n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_Iteration(b, basicnode.Prototype__Map{})
}

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, basicnode.Prototype__Map{})
}
func BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b, basicnode.Prototype__Map{})
}
func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_MapNStrMap3StrInt(b, basicnode.Prototype__Map{})
}

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, basicnode.Prototype__Map{})
}
func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b, basicnode.Prototype__Map{})
}
