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
		{"Bool", declaration.TypeBool{}},
		{"String", declaration.TypeString{}},
		{"DemoMapOfStringToString", declaration.TypeMap{
			KeyType:   "String",
			ValueType: declaration.TypeName("String"),
		}},
	}
	y, _ := os.OpenFile("test/gen_types.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	gm := generationMonad{
		typesFile: y,
	}
	fmt.Fprintf(gm.typesFile, "package whee\n\n")
	for _, x := range fixture {
		gm.writeType(x.name, x.typ)
	}
}
