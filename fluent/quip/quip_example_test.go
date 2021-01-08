package quip_test

import (
	"os"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/quip"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func Example() {
	var err error
	n := quip.BuildMap(&err, basicnode.Prototype.Any, 4, func(ma ipld.MapAssembler) {
		quip.AssignMapEntryString(&err, ma, "some key", "some value")
		quip.AssignMapEntryString(&err, ma, "another key", "another value")
		quip.AssembleMapEntry(&err, ma, "nested map", func(na ipld.NodeAssembler) {
			quip.AssembleMap(&err, na, 2, func(ma ipld.MapAssembler) {
				quip.AssignMapEntryString(&err, ma, "deeper entries", "deeper values")
				quip.AssignMapEntryString(&err, ma, "more deeper entries", "more deeper values")
			})
		})
		quip.AssembleMapEntry(&err, ma, "nested list", func(na ipld.NodeAssembler) {
			quip.AssembleList(&err, na, 2, func(la ipld.ListAssembler) {
				quip.AssignListEntryInt(&err, la, 1)
				quip.AssignListEntryInt(&err, la, 2)
			})
		})
	})
	if err != nil {
		panic(err)
	}
	dagjson.Encoder(n, os.Stdout)

	// Output:
	// {
	// 	"some key": "some value",
	// 	"another key": "another value",
	// 	"nested map": {
	// 		"deeper entries": "deeper values",
	// 		"more deeper entries": "more deeper values"
	// 	},
	// 	"nested list": [
	// 		1,
	// 		2
	// 	]
	// }
}
