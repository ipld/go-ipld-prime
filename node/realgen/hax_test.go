package realgen

import (
	"os/exec"
	"testing"

	"github.com/ipld/go-ipld-prime/node/tests"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go"
)

// i am the worst person and this is the worst code
// but it does do codegen when you test this package!
// (it's also legitimately trash tho, because if you get a compile error, you have to manually rm the relevant files, which is not fun.)
func init() {
	pkgName := "realgen"
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &gengo.AdjunctCfg{}
	ts.Accumulate(schema.SpawnInt("Int"))
	ts.Accumulate(schema.SpawnString("String"))
	ts.Accumulate(schema.SpawnStruct("Msg3",
		[]schema.StructField{
			schema.SpawnStructField("whee", ts.TypeByName("Int"), false, false),
			schema.SpawnStructField("woot", ts.TypeByName("Int"), false, false),
			schema.SpawnStructField("waga", ts.TypeByName("Int"), false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	gengo.Generate(".", pkgName, ts, adjCfg)
	exec.Command("go", "fmt").Run()
}

func BenchmarkMapStrInt_3n_AssembleStandard(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleStandard(b, _Msg3__Style{})
}
func BenchmarkMapStrInt_3n_AssembleEntry(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_AssembleEntry(b, _Msg3__Style{})
}
func BenchmarkMapStrInt_3n_Iteration(b *testing.B) {
	tests.SpecBenchmarkMapStrInt_3n_Iteration(b, _Msg3__Style{})
}

func BenchmarkSpec_Marshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt(b, _Msg3__Style{})
}
func BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b *testing.B) {
	tests.BenchmarkSpec_Marshal_Map3StrInt_CodecNull(b, _Msg3__Style{})
}

func BenchmarkSpec_Unmarshal_Map3StrInt(b *testing.B) {
	tests.BenchmarkSpec_Unmarshal_Map3StrInt(b, _Msg3__Style{})
}
