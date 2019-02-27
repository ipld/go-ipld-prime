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


the Node interface
------------------

TODO


Node implementations
--------------------

TODO


using NodeBuilder
-----------------

TODO


Typed Nodes
-----------

TODO

Everything in the IPLD *Schema* layer is *also* a node.
This sometimes means you'll jump over *several* nodes in the raw Data Model
representation of the data when traversing one link in the Schema layer.
