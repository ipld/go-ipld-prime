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

	tString := schema.SpawnString("String")

	f = openOrPanic("_test/tString.go")
	EmitFileHeader(pkgName, f)
	EmitEntireType(NewStringReprStringGenerator(pkgName, tString, adjCfg), f)
}
