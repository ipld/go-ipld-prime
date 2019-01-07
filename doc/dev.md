Developer Docs
==============

Core Concepts
-------------

### IPLD Big Picture

IPLD is a system for describing, storing, sharing, and hashing data.

IPLD comes in several layers:

1. The IPLD Data Model
2. The IPLD Schema system
3. "IPFN"

The IPLD Data Model is a specification of the core serializable types of data.
You can think of it as roughly isomorphic to the familiar JSON and CBOR
serialization systems (and yes, those parsers usually work "out of the box").
You can read more about
[the Data Model specification in the Specs repo](https://github.com/ipld/specs/blob/master/IPLD-Data-Model-v1.md).

The IPLD Schema system builds atop the IPLD Data Model to add a system of
user-defined types.  These type schema are meant to be easy to write, easy to
read, and easy for the machine to parse and check efficiently.
IPLD Schema bring common concepts like "structs" with named "fields", as well
as "enums" and even "unions" -- and make clear definitions of how to map all of
these concepts onto the basic serializable types of the Data Model.

You can use IPLD at any of these levels you choose.  We recommend having a
schema for your data, because it makes validity checking, sanitization,
migration, and documentation all easier in the long run -- but the Data Model
is also absolutely usable stand-alone.

IPFN is mentioned here as an intent to design something between "dependent types"
and a full network and host agnostic computation platform.  It is not part of
this repo.

### IPLD Nodes

Everything is a Node.  Maps are a node.  Arrays are node.  Ints are a node.
Everything in the IPLD Data Model is a Node.

Nodes are traversable.  Given a node which is one of the recursive kinds
(e.g. map or array), you can list child indexes, and you can traverse that
index to get another Node.

Nodes are trees.  It is not permissible to make a cycle of Node references.

Everything in the IPLD *Schema* layer is *also* a node.
This sometimes means you'll jump over *several* nodes in the raw Data Model
representation of the data when traversing one link in the Schema layer.

(Go programmers already familiar with the standard library `reflect` package
may find it useful to think of `ipld.Node` as similar to `reflect.Value`.)

Code Layout
-----------

### Core / Data Model packages

- `go-ipld-prime` -- the core interfaces are defined here.
	- `ipld.Node` -- see the section on Nodes above; this is *the* interface.
	- `ipld.ReprKind` -- this enumeration describes all the major kinds at the Data Model layer.
- `go-ipld-prime/fluent` -- a user-facing helper package.
	- `fluent.Node` -- does *not* implement `ipld.Node` for once; its methods don't return errors, but carry them until checked instead.
- `go-ipld-prime/impl/free` -- one of the implementations of `ipld.Node`.
	- `ipldfree.Node` -- internally uses go wildcard types to store all data.  Not itself hashable.
- `go-ipld-prime/impl/cbor` -- one of the implementations of `ipld.Node`.
	- `ipldcbor.Node` -- stores data as a skiplist over cbor-formatted byte slices.  Hashable!
- `go-ipld-prime/impl/bind` -- one of the implementations of `ipld.Node`.
	- `ipldbind.Node` -- uses `refmt` to map data onto user-provided Go types.  Not itself hashable.

Node trees can be heterogeneous of their implementation types.
For a complex example:
one should be able to deserialize some data into an `ipldcbor.Node`;
replace *some* substructure of the tree with `ipldfree.Node` datum;
and convert the whole tree back to `ipldcbor.Node` to serialize out again.
In doing so, we might expect that all of the data which started in a CBOR-backed
node and was never replaced with new data would actually be carried by reference
into the final tree, and as a result, re-serialization is the simple act of
writing the bytes back out again as they came.

### Schema system packages

// Note: this is somewhat more aspirational as yet.

- `go-ipld-prime/typed` -- functions for operating over the typed nodes.
- `go-ipld-prime/typed/ast` -- a (self-describing!) tree of schema declarations.
	- n.b. `ast.Type` nodes have to be reified into `typed.Type` to be used; this involves running a completeness check on the whole set of `ast.Node` at once (e.g. making sure no references dangle, etc).
- `go-ipld-prime/typed/ast/parser` -- (future work) a parser for the schema DSL.
- `go-ipld-prime/typed/ast/parser/fmt` -- (future work) a canonicalizing printer for the schema DSL.
- `go-ipld-prime/advlay/chamt` -- (future work) "advanced layouts" for CHAMT.

Note the relationship between the `advlay` and the `typed` families of packages
is not fully fleshed out.  It seems likely that the Schema system will provide
a good clear place to indicate and configure the use of Advanced Layouts, but
perhaps we'll additionally find a good API that allows Advanced Layouts to be
exposed as an `ipld.Node` that works with purely Data Model layer concepts too.
(All Advanced Layouts we've considered so far present as either a map, array,
or byte range; so making them work as pure Data Model should be feasible.)
