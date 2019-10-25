package gengo

import (
	"io"
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

	emitType := func(tg typedNodeGenerator, w io.Writer) {
		tg.EmitNativeType(w)
		tg.EmitNativeAccessors(w)
		tg.EmitNativeBuilder(w)
		tg.EmitNativeMaybe(w)
		tg.EmitNodeType(w)
		tg.EmitTypedNodeMethodType(w)
		tg.EmitNodeMethodReprKind(w)
		tg.EmitNodeMethodLookupString(w)
		tg.EmitNodeMethodLookup(w)
		tg.EmitNodeMethodLookupIndex(w)
		tg.EmitNodeMethodLookupSegment(w)
		tg.EmitNodeMethodMapIterator(w)
		tg.EmitNodeMethodListIterator(w)
		tg.EmitNodeMethodLength(w)
		tg.EmitNodeMethodIsUndefined(w)
		tg.EmitNodeMethodIsNull(w)
		tg.EmitNodeMethodAsBool(w)
		tg.EmitNodeMethodAsInt(w)
		tg.EmitNodeMethodAsFloat(w)
		tg.EmitNodeMethodAsString(w)
		tg.EmitNodeMethodAsBytes(w)
		tg.EmitNodeMethodAsLink(w)

		tg.EmitNodeMethodNodeBuilder(w)
		tnbg := tg.GetNodeBuilderGen()
		tnbg.EmitNodebuilderType(w)
		tnbg.EmitNodebuilderConstructor(w)
		tnbg.EmitNodebuilderMethodCreateMap(w)
		tnbg.EmitNodebuilderMethodAmendMap(w)
		tnbg.EmitNodebuilderMethodCreateList(w)
		tnbg.EmitNodebuilderMethodAmendList(w)
		tnbg.EmitNodebuilderMethodCreateNull(w)
		tnbg.EmitNodebuilderMethodCreateBool(w)
		tnbg.EmitNodebuilderMethodCreateInt(w)
		tnbg.EmitNodebuilderMethodCreateFloat(w)
		tnbg.EmitNodebuilderMethodCreateString(w)
		tnbg.EmitNodebuilderMethodCreateBytes(w)
		tnbg.EmitNodebuilderMethodCreateLink(w)

		tg.EmitTypedNodeMethodRepresentation(w)
		rng := tg.GetRepresentationNodeGen()
		if rng == nil { // FIXME: hack to save me from stubbing tons right now, remove when done
			return
		}
		rng.EmitNodeType(w)
		rng.EmitNodeMethodReprKind(w)
		rng.EmitNodeMethodLookupString(w)
		rng.EmitNodeMethodLookup(w)
		rng.EmitNodeMethodLookupIndex(w)
		rng.EmitNodeMethodLookupSegment(w)
		rng.EmitNodeMethodMapIterator(w)
		rng.EmitNodeMethodListIterator(w)
		rng.EmitNodeMethodLength(w)
		rng.EmitNodeMethodIsUndefined(w)
		rng.EmitNodeMethodIsNull(w)
		rng.EmitNodeMethodAsBool(w)
		rng.EmitNodeMethodAsInt(w)
		rng.EmitNodeMethodAsFloat(w)
		rng.EmitNodeMethodAsString(w)
		rng.EmitNodeMethodAsBytes(w)
		rng.EmitNodeMethodAsLink(w)

		rng.EmitNodeMethodNodeBuilder(w)
		rnbg := rng.GetNodeBuilderGen()
		rnbg.EmitNodebuilderType(w)
		rnbg.EmitNodebuilderConstructor(w)
		rnbg.EmitNodebuilderMethodCreateMap(w)
		rnbg.EmitNodebuilderMethodAmendMap(w)
		rnbg.EmitNodebuilderMethodCreateList(w)
		rnbg.EmitNodebuilderMethodAmendList(w)
		rnbg.EmitNodebuilderMethodCreateNull(w)
		rnbg.EmitNodebuilderMethodCreateBool(w)
		rnbg.EmitNodebuilderMethodCreateInt(w)
		rnbg.EmitNodebuilderMethodCreateFloat(w)
		rnbg.EmitNodebuilderMethodCreateString(w)
		rnbg.EmitNodebuilderMethodCreateBytes(w)
		rnbg.EmitNodebuilderMethodCreateLink(w)
	}

	f := openOrPanic("_test/minima.go")
	emitMinima(f)

	tString := schema.SpawnString("String")
	tInt := schema.SpawnInt("Int")
	tBytes := schema.SpawnBytes("Bytes")
	tLink := schema.SpawnLink("Link")
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
	emitFileHeader(f)
	emitType(NewGeneratorForKindString(tString), f)

	f = openOrPanic("_test/tInt.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindInt(tInt), f)

	f = openOrPanic("_test/tBytes.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindBytes(tBytes), f)

	f = openOrPanic("_test/tLink.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindLink(tLink), f)

	f = openOrPanic("_test/tIntList.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindList(tIntList), f)

	f = openOrPanic("_test/tNullableIntList.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindList(tNullableIntList), f)

	f = openOrPanic("_test/Stract.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract), f)

	f = openOrPanic("_test/Stract2.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract2), f)

	f = openOrPanic("_test/Stract3.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract3), f)

	f = openOrPanic("_test/Stroct.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStroct), f)

	f = openOrPanic("_test/KindsStroct.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tKindsStroct), f)
}
