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

Excursions
----------

Some literature in type theory refers to "open" vs "closed" types, particularly
in regard to unions and enumerations (sum types).

Roughly, "closed" types are those where *all* values are known and countable;
"open" types allow a "default" case for handling unknown values.

Aside from this introduction, we won't use the terms "open" and "closed" much.
All schema types are *like* "closed"; but they're also inherently "open" since
we are of course handling data which may have existed outside of the schema.

In IPLD, the data at the Data Model layer is always "open" in nature; and
at the Schema layer we treat it as "closed".  As such, we don't spend much
futher time with the "open"/"closed" distinction; it's simply "does this data
match the schema or not?".

Most go-ipld-prime APIs for handling typed data will frontload the schema match
checking -- by the time a handle to the document has been returned, the entire
piece of data is verified to match the schema.
There are also some optional ways to use the library which
defer the open->closed mapping until midway through your handling of data,
in exchange for the schema mismatch error becoming something that needs handling
at that point in your code rather than up front.

Schemas and Migration
---------------------

Fundamental to our approach to schemas is an understanding:

> Data Never Changes.  Only our interpretation varies.

Data can be created under one schema, and interpreted later under another.
Data may predate or be created without any kind of schema at all.
All of this needs to be fine.

Moreover, before talking about migration, it's important to note that we
don't allow the comforting, easy notion that migration is a one-way process,
or can be carried out atomically at one magically instantaneous point in time.
Because data is immutable, and producing updated versions of it doesn't make
the older version of the data go away, migration is less a thing that you do;
and more a state of mind.  Migration has to be seamless at any time.

### Using Schema Match checking as Version Detection

We don't include any built-in/blessed concepts of versioning in IPLD Schemas.
It's not necessary: we have rich primitives which can be used to build
either explicit versioning or version detection, at your option.

Since it's easy to check if a schema fits over a piece of data, it's
easy to simply probe a series of schemas until finding one that fits.
Therefore, any constraint a schema makes has the potential to be used
for version detection!

There are a handful of recognizable patterns that are used frequently:

- Using a union to get nominative typing at the document root.
  - e.g. `{"foo": {...}}`, using "foo" as the type+version hint.
  - See the schema-schema for an example of this!
  - Any union representation will do.
- Using a "version" field, plus manual unpacking.
  - e.g. `{"version": "1.2.3", "data":{...}}`
  - This can be implemented using unions of either envelope or inline representation!
    - However, it might not be best to do so: this requires that the multiple
	  versions be implemented *within* your one schema!  Typically it's more
	  composable and maintainable to have a separate schema per version.
  - This can be implemented by double-unpacking.  E.g., match once with a struct
    with fields for version (keep it) and content (dev/null it); and match again
	with a more complete schema chosen based on the version.

(Currently, this probing is left to the library user.  More built-in features
around this are expected to come in the future.)

(In the future, we may also be able to construct some specialized schemas that
suggest jumping to another schema specifically and directly (rather than
linear probing); some research required.  (Ideally this would work consistently
regardless of the ordering of fields in the arriving data, but there's some
tension between that and performance.)  It might also be possible to construct
these as a user already!)

### Some comments on versioning-theory

There are different philosophies of versioning: namely, explicit versioning
 which to use is a choice.

In short, explicit versioning tends towards fragility and is not particularly
fork/community/decentralization-friendly.  Version detection -- also known as
its generalized cousin, *Feature* detection -- is strictly more powerful, but
tends to require more thought to deploy effectively.

Explicit versioning tends to treat version numbers as a junk drawer, upon
which we can heap unbounded amounts of not-necessarily-related semantics.
This is a temptation which can be migitated through diligence, but the
fundamental incentive is always there: like global variables in programming,
a document-global explicit version allows lazy coding and fosters presumptions.

Version/feature detection has the potential to become a fractal.
Using it well thus *also* requires diligence.  However, there is no built-in
siren temptation to misuse them in the same way as explicit versioning; the
trade-offs in complexity tend to be make themselves fairly pronounced and
as such are relatively easily communicated.

It's impossible to make a blanket prescription of how to associate version
information with data; IPLD Schemas makes either choice viable.

### Strongly linked Schemas

It is possible to have a document which links directly to its own Schema!
Since IPLD Schemas are themselves representable in IPLD, it's outright trivial
to make an object containing a CID linking to a Schema.

This may be useful -- in particular, it certainly solves any issue of chosing
unique version strings in using explicit versioning! -- but it is also useless,
by definion, to *migration*.

Migration means wanting to treat old data as new data matching a new schema.
Knowing which other schema is stated to match the data can be a useful input
to deciding how to treat that data, but -- unless you're okay using that *exact*
schema, and it's what your application logic is already built against -- that
knowledge doesn't fully specify what to do to turn that data into what you want.

### Actually Migrating!

