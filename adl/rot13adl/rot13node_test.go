package rot13adl

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestLogicalNodeRoundtrip(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	qt.Assert(t, err, qt.IsNil)
	n := nb.Build()
	// Inspect the high level node.
	s, err := n.AsString()
	qt.Check(t, err, qt.IsNil)
	qt.Check(t, s, qt.Equals, "abcd")
}

func TestNodeInternal(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	qt.Assert(t, err, qt.IsNil)
	n := nb.Build()
	// Poke its insides directly to see that the transformation occured.
	qt.Check(t, n.(*_R13String).synthesized, qt.Equals, "abcd")
	qt.Check(t, n.(*_R13String).raw, qt.Equals, "nopq")
}

func TestReify(t *testing.T) {
	t.Run("using unspecialized raw node", func(t *testing.T) {
		// Build substrate-shaped data using basicnode.
		sn := basicnode.NewString("nopq")
		// Reify it.
		synth, err := Reify(sn)
		// Inspect the resulting high level node.
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, synth.Kind(), qt.Equals, datamodel.Kind_String)
		s, err := synth.AsString()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, qt.Equals, "abcd")
	})
	t.Run("using substrate node", func(t *testing.T) {
		// Build substrate-shaped data, in the substrate type right from the start.
		snb := Prototype.SubstrateRoot.NewBuilder()
		snb.AssignString("nopq")
		sn := snb.Build()
		// Reify it.
		synth, err := Reify(sn)
		// Inspect the resulting high level node.
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, synth.Kind(), qt.Equals, datamodel.Kind_String)
		s, err := synth.AsString()
		qt.Check(t, err, qt.IsNil)
		qt.Check(t, s, qt.Equals, "abcd")
	})
}

func TestInspectingSubstrate(t *testing.T) {
	// Build high level node.
	nb := Prototype.Node.NewBuilder()
	err := nb.AssignString("abcd")
	qt.Assert(t, err, qt.IsNil)
	n := nb.Build()
	// Ask it about its substrate, and inspect that.
	sn := n.(R13String).Substrate()
	// TODO: It's unfortunate this is only available as a concrete type cast: we should probably make a standard feature detection interface with `Substrate()`.
	//  Is it reasonable to make this part of a standard feature detection pattern,
	//   and make that interface reside in the main IPLD package?  Or in an `adl` package that contains such standard interfaces?
	ss, err := sn.AsString()
	qt.Check(t, err, qt.IsNil)
	qt.Check(t, ss, qt.Equals, "nopq")
}
