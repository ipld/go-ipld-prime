IPLD Schema system
==================

// Disclaimer: This document will eventually move to the `ipld/specs` repo.
It's currently incubating here in the same repo as the first implementation
for practical developmental flow reasons.

Goals
-----

1. Provide a reasonable, easy to use, general type system for declaring
  useful properties about data.
2. Compose nicely over IPLD Data Model.
3. Be efficient to verify (e.g. roughly linear in the size of the data
  and schema; and absolutely not turing complete).
4. Be language-agnostic: many compatible implementations of the schema
  checker tools should exist, as well as bindings for every language.
5. Assist rather than obstruct migration.  (We expect data to exist from
  *before* the current schemas; we need to work well on this case.)

Features
--------

Well-understood basics (sum types, product types, some specific recursive types):

- typed maps (e.g. `{Foo:Bar}`)
- typed arrays (e.g. `[Foo]`)
- typedef'd primitives (e.g. `type Foo int`)
- typed structs (e.g. `struct{ f Foo; b Bar }`)
- typed unions (several styles)
- enums (over strings)
- "advanced layouts" (more on that later)

Bonus features:

- simple syntax for 'maybe' (e.g. `struct{ f nullable Foo }` or `[nullable Foo]`)
- clear syntax for non-required fields (e.g. `struct { f optional Foo }`)

Some non-features:

- [dependent types](https://en.wikipedia.org/wiki/Dependent_type)
- anything that would introduce turing-completeness or otherwise unpredictable
  computational complexity to the checking process.

(n.b. It's not that we don't think dependent types are cool -- they are -- just
that this is not the place for them.  A dependently typed system which can
reason about and produce declarations in the IPLD Schema system would be
extremely neat!)

All types are defacto *structural*, rather than *nominative*.
Because we expect data to exist independently of the schema, it follows that
every part of the system has to start by checking that the data matches the
schema.  This doesn't mean you can't have nominative-style behaviors: it's
just that you get them by using features like typed unions, which effectively
give you a nominative-style behavior while leaving the configuration of it
clearly in the hands of your schema.

Implementation
--------------

There are three major components of the implementation:
the schema representation,
the reified schema,
and schema adjuncts.

### Schema Representation

// n.b. sometimes called `ast.*` -- but unsure if "AST" is really the best
term for something that's a full serializable thing itself, not just an IR.

### Reified Schema

Reified schema refers to parts of the code which handle the fully processed
schema info -- it's distinct from the Schema Representation code because it's
allowed to contain pointers (including cyclic references), etc.

The Reified schema can be computed purely from the schema representation.

### Schema Amendments/Adjuncts

Schema Adjuncts are more properties which can be attached to a Reified schema.
Adjuncts are -- as the name suggests -- adjoined to the schema, but technically
not entirely a part of it.

// todo: document behavior tree patterns for handling concepts like "defaults",
schema stack probing, options for programmatic callbacks (for fancy migrations),
etc

### Slurping

// particular to ipldcbor.Node (and other serializables)?

Excursions
----------

Some literature in type theory refers to "open" vs "closed" types, particularly
in regard to unions and enumerations (sum types).

Roughly, "closed" types are those where *all* values are known and countable;
"open" types allow a "default" case for handling unknown values.

Aside from this introduction, we won't use the terms "open" and "closed" much.
All schema types are *like* "closed"; but they're also inherently "open" since
we are of course handling data which may have existed outside of the schema.

In IPLD schema tooling, we always coerce data from its "open" nature to our
"closed" treatment of it by frontloading that check when handling the whole
document.  As such, the distinction isn't particularly useful to make.

// todo: more detailed behavior trees
