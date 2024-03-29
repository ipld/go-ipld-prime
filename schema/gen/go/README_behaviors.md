Behaviors of Generated Types
============================

Types generated by `gogen` obey the rules that already exist for the IPLD Data Model,
and the rules that exist for Schema typed nodes.

There are some further details of behavior that might be worth noting, though.

Some of these details aren't necessarily programmer-friendly(!),
and arise from prioritizing performance over other concerns;
so especially watch out for these as you develop against the code output by this tool.

### retaining assemblers beyond their intended lifespan is not guaranteed to be safe

There is no promise of nice-to-read errors if you over-hold child assemblers beyond their valid lifespan.
`NodeAssembler` values should not be retained for any longer than they're actively in use.

- We **do** care about making things fail hard and fast rather than potentially leak inappropriate mutability.
- We do **not** care about making these errors pretty (it's high cost to do so, and code that hits this path is almost certainly statically (and hopefully fairly obviously) wrong).

In some cases it may also be the case that a `NodeAssembler` that populates the internals of some large structure
may become invalid (because of state transitions that block inappropriate mutability),
and yet become possible to use again later (because of coincidences of how we reuse memory internally for efficiency reasons).
We don't reliably raise errors in some of these situations, for efficiency reasons, but wish we could.
Users of the generated code should not rely on these behaviors:
it results in difficult-to-read code in any case,
and such internal details should not be considered part of the intended public API (e.g., such details may be subject to change without notice).

### absent values

Iterating a type-level node with optional fields will yield the field key and the `datamodel.Absent` constant as a value.
Getting a such a field will also yield the `datamodel.Absent` constant as a value, and will not return a "not found" error.

Attempting to *assign* an `datamodel.Absent` value, however --
via the `NodeAssembler.AssignNode` function (none of the other function signatures permit expressing this) --
will result in an `datamodel.ErrWrongKind` error.

// Seealso some unresolved todos in the [HACKME_wip](HACKME_wip.md) document regarding how absent values are handled.
