package gengo

import (
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/schema"
)

func TestSmoke(t *testing.T) {
	os.Mkdir("_test", 0755)
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}
	var f *os.File

	pkgName := "whee"
	adjCfg := &AdjunctCfg{}

	f = openOrPanic("_test/minima.go")
	EmitInternalEnums(pkgName, f)

	tString := schema.SpawnString("String")

	tStroct := schema.SpawnStruct("Stroct",
		[]schema.StructField{
			schema.SpawnStructField("f1", tString, false, false),
			schema.SpawnStructField("f2", tString, true, false),
			schema.SpawnStructField("f3", tString, true, true),
			schema.SpawnStructField("f4", tString, false, true),
		},
		schema.StructRepresentation_Map{},
	)

	f = openOrPanic("_test/tString.go")
	EmitFileHeader(pkgName, f)
	EmitEntireType(NewStringReprStringGenerator(pkgName, tString, adjCfg), f)

	f = openOrPanic("_test/tStroct.go")
	EmitFileHeader(pkgName, f)
	EmitEntireType(NewStructReprMapGenerator(pkgName, tStroct, adjCfg), f)
}
