package typegen

import (
	"fmt"
	"os"
	"testing"

	declaration "github.com/ipld/go-ipld-prime/typed/declaration"
)

func Test(t *testing.T) {
	fixture := []struct {
		name declaration.TypeName
		typ  declaration.Type
	}{
		// essential primitives:
		{"Bool", declaration.TypeBool{}},
		{"String", declaration.TypeString{}},
		// warmups:
		{"DemoMapOfStringToString", declaration.TypeMap{
			KeyType:   "String",
			ValueType: declaration.TypeName("String"),
		}},
		// okay, of real interest now:
		{"TypeName", declaration.TypeString{}},
		{"UnionRepresentation_Keyed", declaration.TypeMap{
			KeyType:   "String",
			ValueType: declaration.TypeName("TypeName"),
		}},
		{"AnonMap__TypeStruct__fields", declaration.TypeMap{ // hacky placeholder name; avoiding using power of TypeTerm for now.
			KeyType:   "String",
			ValueType: declaration.TypeName("StructField"),
		}},
		{"StructField", declaration.TypeStruct{
			// TODO TypeStruct is gonna be a barrel of laughs: `Fields: AnonMap__TypeStruct_fields{}` -- you're gonna need nodebuilders to even hardcode this.
			//   ahaha and the fields are unexported so we literally really do need a builder.
			Fields: declaration.AnonMap__TypeStruct__fields{}, //__Build ?
			// Type:     declaration.TypeName("TypeTerm"),
			// Optional: declaration.TypeName("Bool"),
			// Nullable: declaration.TypeName("Bool"),
		}},
	}
	os.Mkdir("test", 0755)
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}
	gm := generationMonad{
		typesFile:                  openOrPanic("test/types_gen.go"),
		hypergenericInterfacesFile: openOrPanic("test/nodeIface_gen.go"),
		methodsFile:                openOrPanic("test/methods_gen.go"),
	}
	fmt.Fprintf(gm.typesFile, "package whee\n\n")
	fmt.Fprintf(gm.hypergenericInterfacesFile, "package whee\n\n")
	fmt.Fprintf(gm.methodsFile, "package whee\n\n")
	fmt.Fprintf(gm.methodsFile, "import (\n")
	fmt.Fprintf(gm.methodsFile, "\t\"fmt\"\n")
	fmt.Fprintf(gm.methodsFile, ")\n\n")
	for _, x := range fixture {
		gm.writeType(x.name, x.typ)
		gm.writeMethods(x.name, x.typ)
		gm.writeNodeInterfaceMethods(x.name, x.typ)
	}
}
