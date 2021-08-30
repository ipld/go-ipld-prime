package printer

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestSimpleData(t *testing.T) {
	n, _ := qp.BuildMap(basicnode.Prototype.Any, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "some key", qp.String("some value"))
		qp.MapEntry(ma, "another key", qp.String("another value"))
		qp.MapEntry(ma, "nested map", qp.Map(2, func(ma datamodel.MapAssembler) {
			qp.MapEntry(ma, "deeper entries", qp.String("deeper values"))
			qp.MapEntry(ma, "more deeper entries", qp.String("more deeper values"))
		}))
		qp.MapEntry(ma, "nested list", qp.List(2, func(la datamodel.ListAssembler) {
			qp.ListEntry(la, qp.Int(1))
			qp.ListEntry(la, qp.Int(2))
		}))
	})
	qt.Check(t, Sprint(n), qt.Equals, `map{
	string{"some key"}: string{"some value"}
	string{"another key"}: string{"another value"}
	string{"nested map"}: map{
		string{"deeper entries"}: string{"deeper values"}
		string{"more deeper entries"}: string{"more deeper values"}
	}
	string{"nested list"}: list{
		0: int{1}
		1: int{2}
	}
}`)
}
