Advanced Data Layouts
=====================

// Note: Advanced Data Layouts are an open research topic at present.
// Comments here are speculative and subject to change as implementations emerge.

IPLD Schemas may sometimes specify types which are annotated with an
Advanced Layout (for example, "HAMT").
The advanced layout is a special value which must be known to the parsing program.
The syntax of such a type still indicates what essential kind it is (namely,
map or array (-- Future: bytes?)).

For example:

```ipldsch
type FancyMap map {String:String}<HAMT>
```

Types with an advanced layout are still manipulated with the semantics usual for
their essential kind (e.g. `{String:String}<HAMT>` is still fundamentally a map
and handled as such), but may be internally composed of a more structurally
complex layout of many nodes, and these nodes may even span multiple blocks
(connected by Links (CIDs)).

Advanced Layouts allow data to be *sharded*, which allows extremely large
collections of data to be managed.  By sharding large collections of data
across several blocks, it becomes possible to transport those blocks separately,
which in turn enables incremental transfer of large collections,
transfering and loading only subsets of the data as necessary, and so on.
Advanced Layouts give us the best of sharding's performance capabilities
while still presenting the exact same interface to application logic.


Using Advanced Layouts without a Schema
---------------------------------------

Although we introduced Advanced Layouts by way of the schema system, the
implementation doesn't strictly depend on the schema system,
and Advanced Layouts can be used directly by interacting with the library code:
an Advanced Layout is just another implementation of `Node`!

TODO: more docs when we have implementations :)


Handling data without understanding the Advanced Layout
-------------------------------------------------------

All data can be interpreted with *or without* the Advanced Layout information.

Generally speaking, parsing some data with a schema in hand that provides
the Advanced Layout will allow traversal of the advanced structure transparently
(a single path segment may jump across as many internal nodes as the advanced
layout algorithm determines); and parsing some data without a schema is still
entirely possible at the Data Model layer.

Being able to handle all data as completely normal Data Model content even when
at the application layer it may be seen with an Advanced Layout means we can
make all kinds of content addressing, data transfer, pinning, etc features work
in general: all these low level systems work without any need of schemas or
any specific understanding of any advanced layout algorithm.

This duality can also be used intentionally: for example, even when knowing full
well what the schema is for UnixFS, a program may want to sometimes use that
schema (to traverse directories transparently, even when sharded, for example),
and also sometimes disregard it and see the raw structures (so it can fetch
the left-most leaves of a sharded file, for example -- which would effectively
translate to the "start" of a "file").


Other notes
-----------

Note that Advanced Layout types may not be used in nested type definitions in
schemas (e.g., a struct field may have a type `{String:String}<HAMT>`,
but *not* `[{String:String}<HAMT>]`).
There is no reason for this limitation other than limiting syntactic complexity.
(Imagine the readability of error messages regarding advanced layout contents
if there's no short clear type name for the root of the advanced layout -- terse
is critical usability feature here.)

Note that often Advanced Layout implementation will still use a simple map or
array rather than linked objects when the content of the value is small enough.
(This is a popular optimization based on the understanding that introducing more
intermediate nodes and blocks and link traversal necessarily carries some
additional coordination and overhead costs.)  However, not all Advanced Layouts
will do this, and it should not be relied upon in general.
