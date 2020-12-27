Nodes
=====

Everything is a Node.  Maps are a node.  Arrays are node.  Ints are a node.
Everything in the IPLD Data Model is a Node.

Nodes are traversable.  Given a node which is one of the recursive kinds
(e.g. map or array), you can list child indexes, and you can traverse that
index to get another Node.

Nodes are trees.  It is not permissible to make a cycle of Node references.

(Go programmers already familiar with the standard library `reflect` package
may find it useful to think of `ipld.Node` as similar to `reflect.Value`.)


Overview of Important Types
---------------------------

- `ipld.Node` -- see the section on Nodes above; this is *the* interface.
- `ipld.NodeBuilder` -- an interface for building new Nodes.
- `ipld.Kind` -- this enumeration describes all the major kinds at the Data Model layer.
- `fluent.Node` -- similar to `ipld.Node`, but methods don't return errors, instead
  carrying them until actually checked and/or using panics, making them easy to compose
  in long expressions.  See the full page for [Fluent APIs](./fluent.md) for more.


the Node interface
------------------

Node is an interface for inspecting single values in IPLD.
It has methods for
extracting go primitive values for all the scalar node kinds,
'traverse' methods for maps and lists which return individual child nodes,
and iterator methods for maps.

Node exposes only reader methods -- A Node is immutable.

See the [godocs for Node](https://godoc.org/github.com/ipld/go-ipld-prime#Node).


### kinds

The `Node.Kind()` method returns an `ipld.Kind` enum value describing what
kind of data this node contains in terms of the IPLD Data Model.

The validity of many other methods can be anticipated by switching on the kind:
for example, `AsString` is definitely going to error if `Kind() == ipld.Kind_Map`,
and `LookupByString` is definitely going to error if `Kind() == ipld.Kind_String`.

See the [godocs for Kind](https://godoc.org/github.com/ipld/go-ipld-prime#Kind).


Node implementations
--------------------

Since `Node` is an interface, it can have many implementations.
We use different implementations to satisfy different performance objectives.
There are several implementations in the core library, and users can bring their own.

In go-ipld-prime, we have several implementations:

- `impl/free` (imports as `ipldfree`) -- a generic unspecialized implementation
  of `Node` which can contain any kind of content, internally using Go wildcard types.
- `impl/bind` (imports as `ipldbind`) -- an implementation of `Node` which uses
  reflection to bind to existing Go native objects and present them as nodes.
- `impl/cbor` (imports as `ipldcbor`) -- an implementation of `Node` which stores
  all content as cbor-encoded bytes -- interesting primarily because when doing
  a bulk parse of cbor data, this can effectively defer parsing of values until
  they're actually inspected, sometimes enabling considerable processing saving.
- `typed` -- Typed nodes add a few more features onto the base Node interface
  and add additional logical constraints to their contents -- more on this later
  in the [Typed Nodes](#typed-nodes) section.

Different Node implementations can be mixed freely!

For example: we can use `ipldcbor.Decode` to get an `ipldcbor.Node`
(which internally is a sort of skip-list over the original CBOR byte slices),
use `traverse.Transform` to replace just *some* nodes internally with new data
we build out of an `ipldfree.NodeBuilder`, and then use `ipldcbor.Encode` again
to emit the updated CBOR.  In this case, we are both able to use whatever kind
of Node we want for our new data, and since we're keeping the rest of the tree
in an `ipldcbor.Node` form we were able to skip all parsing of parts of the tree
we aren't interested in, and when emitting bytes again we can recycle all the
byte slices from the very first read -- in other words, zero copy operation
(except for specifically the parts we've updated).  This is pretty neat!

> Note: 'cbor' and 'bind' nodes are not yet fully supported.


using NodeBuilder
-----------------

NodeBuilder is an interface that describes creating new Node instances.
NodeBuilders are a typical builder pattern: they are mutable and stateful.

See the [godocs for NodeBuilder](https://godoc.org/github.com/ipld/go-ipld-prime#NodeBuilder).

Each package exporting a `Node` implementation will generally also export a
corresponding `NodeBuilder()` factory func which can be used to start producing
nodes of that implementation type.

Also note the `Node.NodeBuilder()` method, which returns a NodeBuilder from any Node.
NodeBuilder instances obtained from these methods produce new nodes of the same
implementation type (e.g., if it's an `ipldcbor.Node`, it'll give you an
`ipldcbor.NodeBuilder`).

NodeBuilder instances derived from `Node.NodeBuilder()` have several uses:

1. In the `traversal.Transform` system, this mechanism for getting a NodeBuilder
  is how we pick an implementation for creating new nodes as we propagate
  updates through the parents of mutation point.
2. For typed nodes, NodeBuilder instances obtained from `Node.NodeBuilder()`
  also continue to carry type constraints matching the node they came from.
3. Finally, builders from `Node.NodeBuilder()` have some interesting powers
  for the recursive kinds (maps and lists): they can "append"!  Such appends
  are still copy-on-write (e.g., won't mutate the original Node), but can often
  save on memory use.


Typed Nodes
-----------

TODO

Everything in the IPLD *Schema* layer is *also* a node.
This sometimes means you'll jump over *several* nodes in the raw Data Model
representation of the data when traversing one link in the Schema layer.

NodeBuilder for typed nodes functions the same as other NodeBuilder instances,
but again with more constraints on correct usage: NodeBuilder for typed nodes
will enforce the constraint of their associated type when used.
`typed.Type` has a `Type.NodeBuilder` method to yield appropriately constrained
builder instances; and `Node.NodeBuilder` on existing nodes will continue to
carry on the same constraints as the node its derived from.
