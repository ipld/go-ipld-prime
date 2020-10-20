package rot13adl

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basic"
)

func TestLogicalNodeRoundtrip(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	// Inspect the high level node.
	s, err := n.AsString()
	Wish(t, err, ShouldEqual, nil)
	Wish(t, s, ShouldEqual, "abcd")
}

func TestNodeInternal(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	// Poke its insides directly to see that the transformation occured.
	Wish(t, n.(*_R13String).synthesized, ShouldEqual, "abcd")
	Wish(t, n.(*_R13String).raw, ShouldEqual, "nopq")
}

func TestReify(t *testing.T) {
	t.Run("using unspecialized raw node", func(t *testing.T) {
		// Build substrate-shaped data using basicnode.
		sn := basicnode.NewString("nopq")
		// Reify it.
		synth, err := Reify(sn)
		// Inspect the resulting high level node.
		Require(t, err, ShouldEqual, nil)
		Wish(t, synth.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		s, err := synth.AsString()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, "abcd")
	})
	t.Run("using substrate node", func(t *testing.T) {
		// Build substrate-shaped data, in the substrate type right from the start.
		snb := Prototype.SubstrateRoot.NewBuilder()
		snb.AssignString("nopq")
		sn := snb.Build()
		// Reify it.
		synth, err := Reify(sn)
		// Inspect the resulting high level node.
		Require(t, err, ShouldEqual, nil)
		Wish(t, synth.ReprKind(), ShouldEqual, ipld.ReprKind_String)
		s, err := synth.AsString()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, s, ShouldEqual, "abcd")
	})
}

func TestInspectingSubstrate(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	Require(t, err, ShouldEqual, nil)
	n := nb.Build()
	// Ask it about its substrate, and inspect that.
	sn := n.(*_R13String).Substrate()
	// TODO: the cast above isn't available outside this package: we should probably make an interface with `Substrate()` and make it available.
	//  Is it reasonable to make this part of a standard feature detection pattern,
	//   and make that interface reside in the main IPLD package?  Or in an `adl` package that contains such standard interfaces?
	ss, err := sn.AsString()
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ss, ShouldEqual, "nopq")
}
