package basicnode

import (
	"testing"

	"github.com/ipld/go-ipld-prime/_rsrch/nodesolution/node/mixins/tests"
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
