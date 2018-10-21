package ipld

import "github.com/ipfs/go-cid"

type Node interface {
	// GetField resolves a path in the the object and returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	//
	// If a Node is returned, it will be a unrooted node -- that is,
	// it can be used to view the fields below it, but since it was not
	// originally stored as a full node, you cannot immediately take
	// a link to it for embedding in other objects (you'd have to make
	// a new RootNode with the same content first, then store that).
	//
	// A variety of GetField{Foo} methods exists to yield specific types,
	// and will panic upon encountering unexpected types.
	GetField(path string) (interface{}, error)
	GetFieldBool(path string) (bool, error)     // See GetField docs.
	GetFieldString(path string) (string, error) // See GetField docs.
	GetFieldInt(path string) (int, error)       // See GetField docs.
	GetFieldLink(path string) (cid.Cid, error)  // See GetField docs.

	// GetIndex is the equivalent of GetField but for indexing into an array
	// (or a numerically-keyed map).  Like GetField, it returns
	// either a primitive (e.g. string, int, etc), a link (type CID),
	// or another Node.
	//
	// A variety of GetIndex{Foo} methods exists to yield specific types,
	// and will panic upon encountering unexpected types.
	GetIndex(idx int) (interface{}, error)
	GetIndexBool(idx int) (bool, error)     // See GetIndex docs.
	GetIndexString(idx int) (string, error) // See GetIndex docs.
	GetIndexInt(idx int) (int, error)       // See GetIndex docs.
	GetIndexLink(idx int) (cid.Cid, error)  // See GetIndex docs.

	// REVIEW this whole interface is still *very* unfriendly to chaining.
	// Friendly would be returning `(Node)` at all times, and having final
	// dereferencing options on leaf nodes, and treating it as a Maybe at all
	// other times.  Multireturn methods are antithetical to graceful fluent
	// chaining in golang, syntactically.
	// Panics might also be a viable option.  I suspect we do quite usually
	// want to do large amounts of these traversal operations, and bailing at
	// any point is desired to be terse.  We could provide a thunkwrapper.
}

type SerializableNode interface {
	CID() cid.Cid
}

type MutableNode interface {
	// FUTURE: setter methods go here
}
