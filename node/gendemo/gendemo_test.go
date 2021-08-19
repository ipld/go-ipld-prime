package gendemo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
)

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, _Msg3__Prototype{})
}
func BenchmarkMapStrInt_3n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleEntry(b, _Msg3__Prototype{})
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, _Msg3__Prototype{})
}
func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, _Msg3__Prototype{})
}
func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, _Msg3__Prototype{})
}

func BenchmarkSpec_Marshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_MapNStrMap3StrInt(b, _Map__String__Msg3__Prototype{})
}
func BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_MapNStrMap3StrInt(b, _Map__String__Msg3__Prototype{})
}

// the standard 'walk' benchmarks don't work yet because those use selectors and use the prototype we give them for that, which...
//  does not fly: cramming selector keys into assemblers meant for struct types from our test corpus?  nope.
//   this is a known shortcut-become-bug with the design of the 'walk' benchmarks; we'll have to fix soon.
