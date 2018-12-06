package typed

import (
	"github.com/ipld/go-ipld-prime"
)

type Node struct {
	// FUTURE: proxies most methods, plus adds just-in-time type checking on reads.
	// You can use a Validate call to force checking of the entire tree.
	// "Advanced Layouts" (e.g. HAMTs, etc) can be seamlessly presented as a plain map through this interface.
}

type MutableNode struct {
	// FUTURE: another impl of ipld.MutableNode we can return which checks all things at change time.
	// This can proxy to some other implementation type that does real storage.
}

func Validate(ts Universe, t Type, node ipld.Node) []error {
	// TODO need more methods to enable traversal, then come back here
	return nil
}
