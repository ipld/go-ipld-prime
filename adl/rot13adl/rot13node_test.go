package rot13adl

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basic"
)

func TestLogicalNodeRoundtrip(t *testing.T) {
	nb := Prototype{}.NewBuilder()
	err := nb.AssignString("abcd")
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	s, err := n.AsString()
	Wish(t, err, ShouldEqual, nil)
	Wish(t, s, ShouldEqual, "abcd")
}

func TestNodeInternal(t *testing.T) {
	nb := Prototype{}.NewBuilder()
	err := nb.AssignString("abcd")
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	Wish(t, n.(*_R13String).raw, ShouldEqual, "nopq")
}

func TestReify(t *testing.T) {
	sn := basicnode.NewString("nopq")
	synth, err := Reify(sn)
	Require(t, err, ShouldEqual, nil)
	Wish(t, synth.ReprKind(), ShouldEqual, ipld.ReprKind_String)
	s, err := synth.AsString()
	Wish(t, err, ShouldEqual, nil)
	Wish(t, s, ShouldEqual, "abcd")
}
