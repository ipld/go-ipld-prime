package datamodel_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
)

func TestErrWrongKind_String(t *testing.T) {
	qt.Check(t, datamodel.KindSet{}.String(), qt.Equals, `<empty KindSet>`)
	qt.Check(t, datamodel.ErrWrongKind{}.Error(), qt.Equals, `func called on wrong kind: "" called on a INVALID node, but only makes sense on <empty KindSet>`)
}
