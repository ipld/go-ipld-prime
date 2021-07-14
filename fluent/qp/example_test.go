package qp_test

import (
	"os"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

// TODO: can we make ListEntry/MapEntry less verbose?

func Example() {
	n, err := qp.BuildMap(basicnode.Prototype.Any, 4, func(ma ipld.MapAssembler) {
		qp.MapEntry(ma, "some key", qp.String("some value"))
		qp.MapEntry(ma, "another key", qp.String("another value"))
		qp.MapEntry(ma, "nested map", qp.Map(2, func(ma ipld.MapAssembler) {
			qp.MapEntry(ma, "deeper entries", qp.String("deeper values"))
			qp.MapEntry(ma, "more deeper entries", qp.String("more deeper values"))
		}))
		qp.MapEntry(ma, "nested list", qp.List(2, func(la ipld.ListAssembler) {
			qp.ListEntry(la, qp.Int(1))
			qp.ListEntry(la, qp.Int(2))
		}))
	})
	if err != nil {
		panic(err)
	}
	dagjson.Encode(n, os.Stdout)

	// Output:
	// {"some key":"some value","another key":"another value","nested map":{"deeper entries":"deeper values","more deeper entries":"more deeper values"},"nested list":[1,2]}
}
