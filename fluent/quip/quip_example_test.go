package quip_test

import (
	"os"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/quip"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

func Example() {
	nb := basicnode.Prototype.Any.NewBuilder()
	var err error
	quip.BuildMap(&err, nb, 4, func(ma ipld.MapAssembler) {
		quip.MapEntry(&err, ma, "some key", func(va ipld.NodeAssembler) {
			quip.AbsorbError(&err, va.AssignString("some value"))
		})
		quip.MapEntry(&err, ma, "another key", func(va ipld.NodeAssembler) {
			quip.AbsorbError(&err, va.AssignString("another value"))
		})
		quip.MapEntry(&err, ma, "nested map", func(va ipld.NodeAssembler) {
			quip.BuildMap(&err, va, 2, func(ma ipld.MapAssembler) {
				quip.MapEntry(&err, ma, "deeper entries", func(va ipld.NodeAssembler) {
					quip.AbsorbError(&err, va.AssignString("deeper values"))
				})
				quip.MapEntry(&err, ma, "more deeper entries", func(va ipld.NodeAssembler) {
					quip.AbsorbError(&err, va.AssignString("more deeper values"))
				})
			})
		})
		quip.MapEntry(&err, ma, "nested list", func(va ipld.NodeAssembler) {
			quip.BuildList(&err, va, 2, func(la ipld.ListAssembler) {
				quip.ListEntry(&err, la, func(va ipld.NodeAssembler) {
					quip.AbsorbError(&err, va.AssignInt(1))
				})
				quip.ListEntry(&err, la, func(va ipld.NodeAssembler) {
					quip.AbsorbError(&err, va.AssignInt(2))
				})
			})
		})
	})
	if err != nil {
		panic(err)
	}
	n := nb.Build()
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
