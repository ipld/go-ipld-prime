Developer Docs
==============

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
