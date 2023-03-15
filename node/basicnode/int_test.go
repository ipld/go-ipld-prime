package basicnode_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/ipld/go-ipld-prime/must"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestBasicInt(t *testing.T) {
	m := basicnode.NewInt(3)
	b := m.Prototype().NewBuilder()
	b.AssignInt(4)
	n := b.Build()
	u := basicnode.NewUint(5)
	qt.Check(t, must.Int(m), qt.Equals, int64(3))
	qt.Check(t, must.Int(n), qt.Equals, int64(4))
	qt.Check(t, must.Int(u), qt.Equals, int64(5))
}

func TestIntErrors(t *testing.T) {
	x := basicnode.NewInt(3)

	_, err := x.LookupByIndex(0)
	errExpect := `func called on wrong kind: "LookupByIndex" called on a int node \(kind: int\), but only makes sense on list`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	_, err = x.LookupByString("n")
	errExpect = `func called on wrong kind: "LookupByString" called on a int node \(kind: int\), but only makes sense on map`
	qt.Check(t, err, qt.ErrorMatches, errExpect)

	_, err = x.LookupByNode(x)
	errExpect = `func called on wrong kind: "LookupByNode" called on a int node \(kind: int\), but only makes sense on map`
	qt.Check(t, err, qt.ErrorMatches, errExpect)
}
