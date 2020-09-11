package rot13adl

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
)

// Reify attempts to process raw Data Model data as substrate data to synthesize an ADL.
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
func Reify(rawRoot ipld.Node) (ipld.Node, error) {
	// Is it evidentally a valid substrate for this ADL?
	//  This is a pretty trivial check for rot13adl.
	//  (Other ADLs probably want to include some additional data and structure in them which allow more validation here,
	//   though in general, this validation should also usually stick to not crossing any block loading boundaries.)
	if rawRoot.ReprKind() != ipld.ReprKind_String {
		return nil, fmt.Errorf("cannot reify rot13adl: substrate root node is wrong kind (must be string)")
	}

	// Construct and return the reified node.
	//  If we can recognize the rawRoot as being our own substrate types,
	//   we can shortcut some things;
	//  Otherwise, just process it via the data model.
	if x, ok := rawRoot.(*_Substrate); ok {
		return (*_R13String)(x), nil
	}
	// Shortcut didn't work.  Process via the data model.
	s, _ := rawRoot.AsString()
	return &_R13String{
		raw:         s,
		synthesized: unrotate(s),
	}, nil
}
