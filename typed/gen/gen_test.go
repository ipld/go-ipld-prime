package typegen

import (
	"io"
	"os"
	"testing"

	declaration "github.com/ipld/go-ipld-prime/typed/declaration"
)

func TestNuevo(t *testing.T) {
	os.Mkdir("test", 0755)
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}
	f := openOrPanic("test/neu.go")
	emitFileHeader(f)
	emitType := func(tg typeGenerator, w io.Writer) {
		tg.EmitNodeType(w)
		tg.EmitNodeMethodReprKind(w)
		tg.EmitNodeMethodTraverseField(w)
		tg.EmitNodeMethodTraverseIndex(w)
		tg.EmitNodeMethodMapIterator(w)
		tg.EmitNodeMethodListIterator(w)
		tg.EmitNodeMethodLength(w)
		tg.EmitNodeMethodIsNull(w)
		tg.EmitNodeMethodAsBool(w)
		tg.EmitNodeMethodAsInt(w)
		tg.EmitNodeMethodAsFloat(w)
		tg.EmitNodeMethodAsString(w)
		tg.EmitNodeMethodAsBytes(w)
		tg.EmitNodeMethodAsLink(w)
		tg.EmitNodeMethodNodeBuilder(w)
	}
	emitType(generateKindString{"Strang", declaration.TypeString{}}, f)
}
