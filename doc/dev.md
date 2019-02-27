Developer Docs
==============

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
