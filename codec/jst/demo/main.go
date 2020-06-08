package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/codec/jst"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func main() {
	fixture := `[
		{"path": "./foo",  "moduleName": "whiz.org/teamBar/foo", "status": "changed"},
		{"path": "./baz",  "moduleName": "whiz.org/teamBar/baz", "status": "green"},
		{"path": "./quxx", "moduleName": "example.net/quxx",     "status": "lit",
		  "subtable": [
		    {"widget": "shining",       "property": "neat", "familiarity": 14},
		    {"widget": "shimmering",    "property": "neat", "familiarity": 140},
		    {"widget": "scintillating",                     "familiarity": 0},
		    {"widget": "irridescent",   "property": "yes"},
		  ]}
	]`
	nb := basicnode.Style.Any.NewBuilder()
	if err := dagjson.Decoder(nb, bytes.NewBufferString(fixture)); err != nil {
		panic(err)
	}
	n := nb.Build()

	if err := jst.MarshalConfigured(jst.Config{
		Indent: []byte{' ', ' '},
		Color:  jst.Color{Enabled: true},
	}, n, os.Stdout); err != nil {
		fmt.Printf("\nerror: %s\n", err)
		os.Exit(5)
	}
	fmt.Println()
}
