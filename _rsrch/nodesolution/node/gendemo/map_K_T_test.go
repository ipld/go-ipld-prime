package gendemo

import (
	"testing"

	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/node/mixins/tests"
)

func TestGennedMapStrInt(t *testing.T) {
	tests.SpecTestMapStrInt(t, Type__Map_K_T{})
}

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, Type__Map_K_T{})
}
func BenchmarkMapStrInt_3n_AssembleDirectly(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleDirectly(b, Type__Map_K_T{})
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, Type__Map_K_T{})
}

func BenchmarkMapStrInt_25n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleStandard(b, Type__Map_K_T{})
}
func BenchmarkMapStrInt_25n_AssembleDirectly(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_AssembleDirectly(b, Type__Map_K_T{})
}
func BenchmarkMapStrInt_25n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_25n_Iteration(b, Type__Map_K_T{})
}
