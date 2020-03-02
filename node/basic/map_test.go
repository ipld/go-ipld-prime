package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func TestMap(t *testing.T) {
	tests.SpecTestMapStrInt(t, Style__Map{})
	tests.SpecTestMapStrMapStrInt(t, Style__Map{})
}

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, Style__Map{})
}
func BenchmarkMapStrInt_3n_AssembleDirectly(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleDirectly(b, Style__Map{})
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, Style__Map{})
}

func BenchmarkMapStrInt_25n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleStandard(b, Style__Map{})
}
func BenchmarkMapStrInt_25n_AssembleDirectly(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleDirectly(b, Style__Map{})
}
func BenchmarkMapStrInt_25n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_Iteration(b, Style__Map{})
}

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, Style__Map{})
}
func BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b, Style__Map{})
}
func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_MapNStrMap3StrInt(b, Style__Map{})
}

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, Style__Map{})
}
func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b, Style__Map{})
}
