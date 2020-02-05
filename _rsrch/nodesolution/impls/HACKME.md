hackme
======

Design rational are documented here.

This doc is not necessary reading for users of this package,
but if you're considering submitting patches -- or just trying to understand
why it was written this way, and check for reasoning that might be dated --
then it might be useful reading.

### scalars are just typedefs

This is noteworthy because in codegen, this is typically *not* the case:
in codegen, even scalar types are boxed in a struct, such that it prevents
casting values into those types.

This casting is not a concern for ipldfree types, because
A) we don't have any kind of validation rules to make such casting worrying; and
B) since our types are unexported, casting is still blocked by this anyway.

### pointer-vs-value inhabitant consistency

Builders always return interfaces inhabited by pointers.
Some constructor functions (e.g. `String()`) return an interface inhabited by
the bare typedef and no pointer.

This means you can get `*plainString` and `plainString` inhabitants of `Node`.

(In contrast, codegen systems usually make a stance of consistently using
only pointer inhabitants -- so you can do type cast checks more easily.)

This hasn't been a problem yet, but it's something to keep an eye on.
If this causes the slightest irritation, we'll standardize to pointer inhabitants.

### about nodestyles

A `Node`'s `Style()` property varies based on how it was constructed.

If something was *made* as a `plainString` -- i.e.,
either via `String()` or via `Style__String{}.NewBuilder()` --
then its `Style()` will be `Style__String`.

If something was made as an "any" -- i.e.,
via `Style__Any{}.NewBuilder()` which then *happened* to be assigned a string value --
then while the builder will still return something functionally equivalent to `plainString`,
it will carry a `Style()` property that returns `Style__Any`.

This is important because the `traversal` package may want to do a
transformation on some data, and use the `Style()` property to find out
which builder to use for the new data -- and if the value allowed in a position
is "any", then that's the style and builder we should get!

### nodestyles meet recursive assemblers

Asking for a NodeStyle in a recursive assembly process tells you about what
kind of node would be accepted in an `AssignNode(Node)` call.
It does *not* make any remark on the fact it's a key assembler or value assembler
and might be wrapped with additional rules (such as map key uniqueness, field
name expectations, etc).

(Note that it's also not an exclusive statement about what `AssignNode(Node)` will
accept; e.g. in many situations, while a `MyStringType` might be the style
returned, any string kinded node can be used in `AssignNode(Node)` and will be
appropriately converted.)

Any of these paths counts as "recursive assembly process":

- `MapNodeAssembler.KeyStyle()`
- `MapNodeAssembler.ValueStyle()`
- `ListNodeAssembler.ValueStyle()`
- `MapNodeAssembler.AssembleKey().Style()`
- `MapNodeAssembler.AssembleValue().Style()`
- `ListNodeAssembler.AssembleValue().Style()`

### about builders for scalars

The assembler types for scalars (string, int, etc) are pretty funny-looking.
You might wish to make them work without any state at all!

The reason this doesn't fly is that we have to keep the "wip" value in hand
just long enough to return it from the `NodeBuilder.Build` method -- the
`NodeAssembler` contract for `Assign*` methods doesn't permit just returning
their results immediately.

(Another possible reason is if we expected to use these assemblers on
slab-style allocations (say, `[]plainString`)...
however, this is inapplicable at present, because
A) we don't (except places that have special-case internal paths anyway); and
B) the types aren't exported, so users can't either.)

Does this mean that using `NodeBuilder` for scalars has a completely
unnecessary second allocation, which is laughably inefficient?  Yes.
It's unfortunate the interfaces constrain us to this.
**But**: one typically doesn't actually use builders for scalars much;
they're just here for completeness.
So this is less of a problem in practice than it might at first seem.

More often, one will use the "any" builder (which is has a whole different set
of design constraints and tradeoffs);
or, if one is writing code and knows which scalar they need, the exported
direct constructor function for that kind
(e.g., `String("foo")` instead of `Style__String{}.NewBuilder().AssignString("foo")`)
will do the right thing and do it in one allocation (and it's less to type, too).

### maps and list keyAssembler and valueAssemblers have custom scalar handling

Related to the above heading.

Maps and lists in this package do their own internal handling of scalars,
using unexported features inside the package, because they can more efficient.
