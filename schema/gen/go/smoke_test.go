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
			// Every field in this struct (including their order) is exercising an interesting case...
			schema.SpawnStructField("f1", tString, false, false), // plain field.
			schema.SpawnStructField("f2", tString, true, false),  // optional; later we have more than one optional field, nonsequentially.
			schema.SpawnStructField("f3", tString, false, true),  // nullable; but required.
			schema.SpawnStructField("f4", tString, true, true),   // optional and nullable; trailing optional.
			schema.SpawnStructField("f5", tString, true, false),  // optional; and the second one in a row, trailing.
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
