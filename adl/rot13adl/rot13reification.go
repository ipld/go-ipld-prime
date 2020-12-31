package rot13adl

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
)

// Reify examines data in a Node to see if it matches the shape for valid substrate data for this ADL,
// and if so, synthesizes and returns the high-level view of the ADL.
// If it succeeds in recognizing the raw data as this ADL,
// Reify returns a new Node which exhibits the logical behaviors of the ADL;
// otherwise, it returns an error.
//
// The input data can be any implementation of ipld.Node;
// it will be considered purely through that interface.
//
// If your application is expecting ADL data, this pipeline can be optimized
// by using the SubstratePrototype right from the start when unmarshalling;
// then, Reify can detect if the rawRoot parameter is of that implementation,
// and it can save some processing work internally that can be known to already be done.
//
// Reification will generally operate on the data in a single block
// (e.g. this function will not do any additional block loads and unmarshalling).
// This is important because some ADLs handle data so large that loading it all
// eagerly would be impractical (and in some cases outright impossible).
// However, it also necessarily implies that invalid data may lie beyond
// one of those lazy loads, and it won't be discovered at the time of Reify.
//
// In this demo ADL, we don't have multi-block content at all,
// so of course we don't have any additional block loads!
// However, ADL implementations may vary in their approaches to lazy vs eager loading.
// All ADLs should document their exact semantics regarding this --
// especially if it has any implications for boundaries of data validity checking.
//
// REVIEW: this function is currently not conforming to any particular interface;
// if we evolve the contract for ADLs to include an interface for reficiation functions,
// might we need to add context and link loader systems as parameters to it?
// Not all implementations might need it, as per previous paragraph; but some might.
// Reification for multiblock ADLs might also need link loader systems as a parameter here
// so they can capture them as config and hold them for use in future operations that do lazy loading.
//
func Reify(maybeSubstrateRoot ipld.Node) (ipld.Node, error) {
	// Reify is often very easy to implement,
	//  especially if you have an IPLD Schema that specifies the shape of the substrate data:
	// We can just check if the data in maybeSubstrateRoot happens to already be exactly the right type,
	//  and if so, take very direct shortcuts because we already know its been validated in shape;
	// otherwise, we create a new piece of memory for our native substrate memory layout,
	//  and assign into it from the raw node, validating in the process,
	//   which again just leans directly on the shape validation logic already given to us by the schema logic on that type.
	// (Checking the concrete type of maybeSubstrateRoot in search of a shortcut is seemingly a tad redundant,
	//  because the AssignNode path later also has such a check!
	//  However, doing it earlier allows us to avoid an allocation;
	//   the AssignNode path doesn't become available until after NewBuilder is invoked, and NewBuilder is where allocations happen.)

	// Check if we can recognize the maybeSubstrateRoot as being our own substrate types;
	//  if it is, we can shortcut pretty drastically.
	if x, ok := maybeSubstrateRoot.(*_Substrate); ok {
		// In this ADL implementation, the high level node has the exact same memory layout as the substrate root,
		//  and so our only remaining processing here is just to cast them, so that
		//   the node we return has the correct methodset exposed.
		return (*_R13String)(x), nil
	}

	// Shortcut didn't work.  Process via the data model.
	//  The AssignNode method on the substrate type already contains all the logic necessary for this, so we use that.
	nb := Prototype.SubstrateRoot.NewBuilder()
	if err := nb.AssignNode(maybeSubstrateRoot); err != nil {
		fmt.Errorf("rot13adl.Reify failed: data does not match expected shape for substrate: %w", err)
	}
	return (*_R13String)(nb.Build().(*_Substrate)), nil
}
