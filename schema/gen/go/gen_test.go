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

	emitType := func(tg typeGenerator, w io.Writer) {
		tg.EmitNodeType(w)
		tg.EmitNodeMethodReprKind(w)
		tg.EmitNodeMethodTraverseField(w)
		tg.EmitNodeMethodTraverseIndex(w)
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
	}

	f := openOrPanic("_test/minima.go")
	emitMinima(f)

	tStrang := schema.SpawnString("Strang")
	tStract := schema.SpawnStruct("Stract",
		[]schema.StructField{schema.SpawnStructField(
			"aField", tStrang, false, false,
		)},
		schema.StructRepresentation_Map{},
	)
	tStract2 := schema.SpawnStruct("Stract2",
		[]schema.StructField{schema.SpawnStructField(
			"nulble", tStrang, false, true,
		)},
		schema.StructRepresentation_Map{},
	)
	tStract3 := schema.SpawnStruct("Stract3",
		[]schema.StructField{schema.SpawnStructField(
			"noptble", tStrang, true, true,
		)},
		schema.StructRepresentation_Map{},
	)

	f = openOrPanic("_test/tStrang.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindString(tStrang), f)

	f = openOrPanic("_test/Stract.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract), f)

	f = openOrPanic("_test/Stract2.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract2), f)

	f = openOrPanic("_test/Stract3.go")
	emitFileHeader(f)
	emitType(NewGeneratorForKindStruct(tStract3), f)
}
