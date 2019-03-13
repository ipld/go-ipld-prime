IPLD Big Picture
================

IPLD is a system for describing, storing, sharing, and hashing data.


Layers
------

IPLD comes in two distinct layers:

1. The IPLD Data Model
2. The IPLD Schema system

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

You can use IPLD at either of these levels as you prefer.  We recommend having a
schema for your data, because it makes validity checking, sanitization,
migration, and documentation all easier in the long run -- but the Data Model
is also absolutely usable stand-alone.

Some of the most exciting features of IPLD come from the crossover between
the two layers: for example, it's possible to create two distinct Schema specs,
and easily flip data between them since all the data on both sides consistently
fits in the Data Model; migration logic can be written as functions operating
purely on the Data Model layer to handle even complex structural changes.
The ability to freely mix and match strong validation with generic data handling
makes programs possible which would otherwise be cumbersome with either alone.


Nodes, operations, encoding, & linking
--------------------------------------

Three (or four) big picture groups of concepts you'll want to keep in mind when
reading the rest of the docs are nodes, operations, encoding, and linking.

Nodes: every piece of data in IPLD can be handled as a "node".
Strings are nodes; ints are nodes; maps are nodes; etc.
[Nodes have their own full chapter of the documentation](./nodes.md).

Operations: some basic computation we can do *on* a node is an operation.
Traversals and simple updates have built-in functions; more advanced operations
can be built with some use of callbacks.
[Operations have their own full chapter of the documentation](./operations.md).

Encoding: every piece of data in IPLD is serializable.  Any node can be
converted to a tokenized representation, and several built-in serialization
systems are supported, plus a pluggable mechanism for implementing more.
[Encoding has its own full chapter of the documentation](./encoding.md)

Linking: every piece of data in IPLD is serializable ***and hashable***...
and can be linked by those hashes.  These links can be manipulated explicitly,
and also used transparently by some of the built-in operation support.
[Linking has its own full chapter of the documentation](./linking.md)

All of these features work on both the Data Model (unityped) layer
and also with Schema layer (typed) nodes.

