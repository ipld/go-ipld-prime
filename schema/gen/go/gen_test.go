package gengo

import (
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/schema"
)

func TestNuevo(t *testing.T) {
	os.Mkdir("_test", 0755)
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}

	f := openOrPanic("_test/minima.go")
	EmitMinima("whee", f)

	tString := schema.SpawnString("String")
	tInt := schema.SpawnInt("Int")
	tBytes := schema.SpawnBytes("Bytes")
	tLink := schema.SpawnLink("Link")
	tIntLink := schema.SpawnLinkReference("IntLink", tInt)
	tIntList := schema.SpawnList("IntList", tInt, false)
	tNullableIntList := schema.SpawnList("NullableIntList", tInt, true)

	tStract := schema.SpawnStruct("Stract",
		[]schema.StructField{schema.SpawnStructField(
			"aField", tString, false, false,
		)},
		schema.StructRepresentation_Map{},
	)
	tStract2 := schema.SpawnStruct("Stract2",
		[]schema.StructField{schema.SpawnStructField(
			"nulble", tString, false, true,
		)},
		schema.StructRepresentation_Map{},
	)
	tStract3 := schema.SpawnStruct("Stract3",
		[]schema.StructField{schema.SpawnStructField(
			"noptble", tString, true, true,
		)},
		schema.StructRepresentation_Map{},
	)
	tStroct := schema.SpawnStruct("Stroct",
		[]schema.StructField{
			schema.SpawnStructField("f1", tString, false, false),
			schema.SpawnStructField("f2", tString, true, false),
			schema.SpawnStructField("f3", tString, true, true),
			schema.SpawnStructField("f4", tString, false, true),
		},
		schema.StructRepresentation_Map{},
	)

	tKindsStroct := schema.SpawnStruct("KindsStroct",
		[]schema.StructField{
			schema.SpawnStructField("inty", tInt, false, false),
			schema.SpawnStructField("bytey", tBytes, false, false),
			schema.SpawnStructField("linky", tLink, false, false),
			schema.SpawnStructField("intListy", tIntList, false, false),
			schema.SpawnStructField("nullableIntListy", tNullableIntList, false, false),
		},
		schema.StructRepresentation_Map{},
	)
	f = openOrPanic("_test/tString.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindString(tString), f)

	f = openOrPanic("_test/tInt.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindInt(tInt), f)

	f = openOrPanic("_test/tBytes.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindBytes(tBytes), f)

	f = openOrPanic("_test/tLink.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindLink(tLink), f)

	f = openOrPanic("_test/tIntLink.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindLink(tIntLink), f)

	f = openOrPanic("_test/tIntList.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindList(tIntList), f)

	f = openOrPanic("_test/tNullableIntList.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindList(tNullableIntList), f)

	f = openOrPanic("_test/Stract.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindStruct(tStract), f)

	f = openOrPanic("_test/Stract2.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindStruct(tStract2), f)

	f = openOrPanic("_test/Stract3.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindStruct(tStract3), f)

	f = openOrPanic("_test/Stroct.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindStruct(tStroct), f)

	f = openOrPanic("_test/KindsStroct.go")
	EmitFileHeader("whee", f)
	EmitEntireType(NewGeneratorForKindStruct(tKindsStroct), f)
}
