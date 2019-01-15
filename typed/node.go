package typed

import (
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/typed/system"
)

// typed.Node is a superset of the ipld.Node interface, and has additional behaviors.
//
// A typed.Node can be inspected for its typesystem.Type and typesystem.Kind,
// which conveys much more and richer information than the Data Model layer
// ipld.ReprKind.
//
// There are many different implementations of typed.Node.
// One implementation can wrap any other existing ipld.Node (i.e., it's zero-copy)
// and promises that it has *already* been validated to match the typesystem.Type;
// another implementation similarly wraps any other existing ipld.Node, but
// defers to the typesystem validation checking to fields that are accessed;
// and when using code generation tools, all of the generated native Golang
// types produced by the codegen will each individually implement typed.Node.
//
// Note that typed.Node can wrap *other* typed.Node instances.
// Imagine you have two parts of a very large code base which have codegen'd
// components which are from different versions of a schema.  Smooth migrations
// and zero-copy type-safe data sharing between them: We can accommodate that!
type Node interface {
	ipld.Node

	Type() typesystem.Type
}
