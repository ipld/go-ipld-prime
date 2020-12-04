CHANGELOG
=========

Here is collected some brief notes on major changes over time, sorted by tag in which they are first available.

Of course for the "detailed changelog", you can always check the commit log!  But hopefully this summary _helps_.

Note about version numbering: All release tags are in the "v0.${x}" range.  _We do not expect to make a v1 release._
Nonetheless, this should not be taken as a statement that the library isn't _usable_ already.
Much of this code is used in other libraries and products, and we do take some care about making changes.
(If you're ever wondering about stability of a feature, ask -- or contribute more tests ;))

- [Planned/Upcoming Changes](#planned-upcoming-changes)
- [Changes on master branch but not yet Released](#unreleased-on-master)
- [Released Changes Log](#released-changes)


Planned/Upcoming Changes
------------------------

Here are some outlines of changes we intend to make that affect the public API:

- Linking and LinkLoaders will be streamlined; and Codecs will gain more friendly APIs.
	- See https://github.com/ipld/go-ipld-prime/issues/55 for discussion and drafts.
	- This will probably be scheduled for around v0.9, which is anticipated to be somewhere in early 2021.
- Most uses of `int` will become `int64`.
	- Since the IPLD Data Model specifications state that integers must be at least 2^53, using golang `int` on a 32-bit architecture would not be spec compliant, and therefore we must use `int64`.
	- This is scheduled for v0.7, which is anticipated to be somewhere in December 2020.

This is not an exhaustive list of planned changes, and does not include any internal changes, new features, performance improvements, and so forth.
It's purely a list of things you might want to know about as a downstream consumer planning your update cycles.

We will make these changes "soon" (for some definition of "soon").
They are currently not written on the master branch.
The definition of "soon" may vary, in service of a goal to sequence any public API changes in a way that's smooth to migrate over, and make those changes appear at an overall bearable chronological frequency.
Tagged releases will be made when any of these changes land, so you can upgrade intentionally.


Unreleased on master
--------------------

Changes here are on the master branch, but not in any tagged release yet.
When a release tag is made, this block of bullet points will just slide down to the [Released Changes](#released-changes) section.

- Feature: codegen is a reasonably usable alpha!  We now encourage trying it out (but still only for those willing to experience an "alpha" level of friction -- UX still rough, and we know it).
	- Consult the feature table in the codegen package readme: many major features of IPLD Schemas are now supported.
		- Structs with tuple representations?  Yes.
		- Keyed unions?  Yes.
		- Structs with stringjoin representations?  Yes.  Including nested?  _Yes_.
		- Lots of powerful stuff is now available to use.
	- Many generated types now have more methods for accessing them in typed ways (in addition to the usual `ipld.Node` interfaces, which can access the same data, but lose explicit typing).
		- Maps and lists now have both lookup methods and iterators which know the type of the child keys and values explicitly.
	- Cool: when generating unions, you can choose between different implementation strategies (favoring either interfaces, or embedded values) by using Adjunct Config.  This lets you tune for either speed (reduced allocation count) or memory footprint (less allocation size, but more granular allocations).
	- Cyclic references in types are now supported.
		- ... mostly.  Some manual configuration may sometimes be required to make sure the generated structure wouldn't have an infinite memory size.  We'll keep working on making this smoother in the future.
	- Field symbol overrides now work properly.  (E.g., if you have a schema with a field called "type", you can make that work now.  Just needs a field symbol override in the Adjunct Config when doing codegen!)
	- Codegen'd link types now implemented the `schema.TypedLinkNode` interface where applicable.
	- Structs now actually validate all required fields are present before allowing themselves to finish building.
	- Much more testing.  And we've got a nice new declarative testcase system that makes it easier to write descriptions of how data should behave (at both the typed and representation view levels), and then just call one function to run exhaustive tests to make sure it looks the same from every inspectable API.
	- Change: codegen now outputs a fixed set of files.  (Previously, it output one file per type in your schema.)  This makes codegen much more managable; if you remove a type from your schema, you don't have to chase down the orphaned file.  It's also just plain less clutter to look at on the filesystem.
- Demo: as proof of the kind of work that can be done now with codegen, we've implemented the IPLD Schema schema -- the schema that describes IPLD Schema declarations -- using codegen.  It's pretty neat.
	- Future: we'll be replacing most of the current current `schema` package with code based on this generated stuff.  Not there yet, though.  Taking this slow.
- Feature: the `schema` typesystem info packages are improved.
	- Cyclic references in types are now supported.
		- (Mind that there are still some caveats about this when fed to codegen, though.)
	- Graph completeness is now validated (e.g. missing type references emit useful errors)!
- Feature: there's a `traversal.Get` function.  It's like `traversal.Focus`, but just returns the reached data instead of dragging you through a callback.  Handy.
- Feature/bugfix: the DAG-CBOR codec now includes resource budgeting limits.  This means it's a lot harder for a badly-formed (or maliciously formed!) message to cause you to run out of memory while processing it.
- Bugfix: several other panics from the DAG-CBOR codec on malformed data are now nice politely-returned errors, as they should be.
- Bugfix: in codegen, there was a parity break between the AssembleEntry method and AssembleKey+AssembleValue in generated struct NodeAssemblers.  This has been fixed.
- Minor: ErrNoSuchField now uses PathSegment instead of a string.  You probably won't notice (but this was important interally: we need it so we're able to describe structs with tuple representations).
- Bugfix: an error path during CID creation is no longer incorrectly dropped.  (I don't think anyone ever ran into this; it only handled situations where the CID parameters were in some way invalid.  But anyway, it's fixed now.)
- Performance: when `cidlink.Link.Load` is used, it will do feature detection on its `io.Reader`, and if it looks like an already-in-memory buffer, take shortcuts that do bulk operations.  I've heard this can reduce memory pressure and allocation counts nicely in applications where that's a common scenario.
- Feature: there's now a `fluent.Reflect` convenience method.  Its job is to take some common golang structs like maps and slices of primitives, and flip them into an IPLD Node tree.
	- This isn't very high-performance, so we don't really recommend using it in production code (certainly not in any hot paths where performance matters)... but it's dang convenient sometimes.
- Feature: there's now a `traversal.SelectLinks` convenience method.  Its job is to walk a node tree and return a list of all the link nodes.
	- This is both convenient, and faster than doing the same thing using general-purpose Selectors (we implemented it as a special case).
- Demo: you can now find a "rot13" ADL in the `adl/rot13adl` package.  This might be useful reference material if you're interested in writing an ADL and wondering what that entails.
- In progress: we've started working on some new library features for working with data as streams of "tokens".  You can find some of this in the new `codec/codectools` package.
	- Functions are available for taking a stream of tokens and feeding them into a NodeAssembler; and for taking a Node and reading it out as a stream of tokens.
	- The main goal in mind for this is to provide reusable components to make it easier to implement new codecs.  But maybe there will be other uses for this feature too!
	- These APIs are brand new and are _extremely subject to change_, much more so than any other packages in this repo.  If you work with them at this stage, _do_ expect to need to update your code when things shift.


Released Changes
----------------

### v0.5.0

v0.5.0 is a small release -- it just contains a bunch of renames.
There are _no_ semantic changes bundled with this (it's _just_ renames) so this should be easy to absorb.

- Renamed: `NodeStyle` -> `NodePrototype`.
	- Reason: it seems to fit better!  See https://github.com/ipld/go-ipld-prime/issues/54 for a full discussion.
	- This should be a "sed refactor" -- the change is purely naming, not semantics, so it should be easy to update your code for.
	- This also affects some package-scoped vars named `Style`; they're accordingly also renamed to `Prototype`.
	- This also affects several methods such as `KeyStyle` and `ValueStyle`; they're accordingly also renamed to `KeyPrototype` and `ValuePrototype`.
- Renamed: `(Node).Lookup{Foo}` -> `(Node).LookupBy{Foo}`.
	- Reason: The former phrasing makes it sound like the "{Foo}" component of the name describes what it returns, but in fact what it describes is the type of the param (which is necessary, since Golang lacks function overloading parametric polymorphism).  Adding the preposition should make this less likely to mislead (even though it does make the method name moderately longer).
	- This should be a "sed refactor" -- the change is purely naming, not semantics, so it should be easy to update your code for.
- Renamed: `(Node).Lookup` -> `(Node).LookupNode`.
	- Reason: The shortest and least-qualified name, 'Lookup', should be reserved for the best-typed variant of the method, which is only present on codegenerated types (and not present on the Node interface at all, due to Golang's limited polymorphism).
	- This should be a "sed refactor" -- the change is purely naming, not semantics, so it should be easy to update your code for.  (The change itself in the library was fairly literally `s/Lookup(/LookupNode(/g`, and then `s/"Lookup"/"LookupNode"/g` to catch a few error message strings, so consumers shouldn't have it much harder.)
	- Note: combined with the above rename, this method overall becomes `(Node).LookupByNode`.
- Renamed: `ipld.Undef` -> `ipld.Absent`, and `(Node).IsUndefined` -> `(Node).IsAbsent`.
	- Reason: "absent" has emerged as a much, much better description of what this value means.  "Undefined" sounds nebulous and carries less meaning.  In long-form prose docs written recently, "absent" consistently fits the sentence flow much better.  Let's just adopt "absent" consistently and do away with "undefined".
	- This should be a "sed refactor" -- the change is purely naming, not semantics, so it should be easy to update your code for.


### v0.4.0

v0.4.0 contains some misceleanous features and documentation improvements -- perhaps most notably, codegen is re-introduced and more featureful than previous rounds -- but otherwise isn't too shocking.
This tag mostly exists as a nice stopping point before the next version coming up (which is planned to include several API changes).

- Docs: several new example functions should now appear in the godoc for how to use the linking APIs.
- Feature: codegen is back!  Use it if you dare.
	- Generated code is now up to date with the present versions of the core interfaces (e.g., it's updated for the NodeAssembler world).
	- We've got a nice big feature table in the codegen package readme now!  Consult that to see which features of IPLD Schemas now have codegen support.
	- There are now several implemented and working (and robustly tested) examples of codegen for various representation strategies for the same types.  (For example, struct-with-stringjoin-representation.)  Neat!
	- This edition of codegen uses some neat tricks to not just maintain immutability contracts, but even prevent the creation of zero-value objects which could potentially be used to evade validation phases on objects that have validation rules.  (This is a bit experimental; we'll see how it goes.)
	- There are oodles and oodles of deep documentation of architecture design choices recorded in "HACKME_*" documents in the codegen package that you may enjoy if you want to contribute or understand why generated things are the way they are.
	- Testing infrastructure for codegen is now solid.  Running tests for the codegen package will: exercise the generation itself; AND make sure the generated code compiles; AND run behavioral tests against it: the whole gamut, all from regular `go test`.
	- The "node/gendemo" package contains a real example of codegen output... and it's connected to the same tests and benchmarks as other node implementations.  (Are the gen'd types fast?  yes.  yes they are.)
	- There's still lots more to go: interacting with the codegen system still requires writing code to interact with as a library, as we aren't shipping a CLI frontend to it yet; and many other features are still in development as well.  But you're welcome to take it for a spin if you're eager!
- Feature: introduce JSON Tables Codec ("JST"), in the `codec/jst` package.  This is a codec that emits bog-standard JSON, but leaning in on the non-semantic whitespace to produce aligned output, table-like, for pleasant human reading.  (If you've used `column -t` before in the shell: it's like that.)
	- This package may be a temporary guest in this repo; it will probably migrate to its own repo soon.  (It's a nice exercise of our core interfaces, though, so it incubated here.)
- I'm quietly shifting the versioning up to the 0.x range.  (Honestly, I thought it was already there, heh.)  That makes this this "v0.4".


### v0.0.3

v0.0.3 contained a massive rewrite which pivoted us to using NodeAssembler patterns.
Code predating this version will need significant updates to match; but, the performance improvements that result should be more than worth it.

- Constructing new nodes has a major pivot towards using "NodeAssembler" pattern: https://github.com/ipld/go-ipld-prime/pull/49
	- This was a massively breaking change: it pivoted from bottom-up composition to top-down assembly: allocating large chunks of structured memory up front and filling them in, rather than stitching together trees over fragmented heap memory with lots of pointers
- "NodeStyle" and "NodeBuilder" and "NodeAssembler" are all now separate concepts:
	- NodeStyle is more or less a builder factory (forgive me -- but it's important: you can handle these without causing allocations, and that matters).
	  Use NodeStyle to tell library functions what kind of in-memory representation you want to use for your data.  (Typically `basicnode.Style.Any` will do -- but you have the control to choose others.)
	- NodeBuilder allocates and begins the assembly of a value (or a whole tree of values, which may be allocated all at once).
	- NodeAssembler is the recursive part of assembling a value (NodeBuilder implements NodeAssembler, but everywhere other than the root, you only use the NodeAssembler interface).
- Assembly of trees of values now simply involves asking the assembler for a recursive node to give you assemblers for the keys and/or values, and then simply... using them.
	- This is much simpler (and also faster) to use than the previous system, which involved an awkward dance to ask about what kind the child nodes were, get builders for them, use those builders, then put the result pack in the parent, and so forth.
- Creating new maps and lists now accepts a size hint argument.
	- This isn't strictly enforced (you can provide zero, or even a negative number to indicate "I don't know", and still add data to the assembler), but may improve efficiency by reducing reallocation costs to grow structures if the size can be estimated in advance.
- Expect **sizable** performance improvements in this version, due to these interface changes.
- Some packages were renamed in an effort to improve naming consistency and feel:
	- The default node implementations have moved: expect to replace `impl/free` in your package imports with `node/basic` (which is an all around better name, anyway).
	- The codecs packages have moved: replace `encoding` with `codec` in your package imports (that's all there is to it; nothing else changed).
- Previous demos of code generation are currently broken / disabled / removed in this tag.
	- ...but they'll return in future versions, and you can follow along in branches if you wish.
- Bugfix: dag-cbor codec now correctly handles marshalling when bytes come after a link in the same object. [[53](https://github.com/ipld/go-ipld-prime/pull/53)]

### v0.0.2

- Many various performance improvements, fixes, and docs improvements.
- Many benchmarks and additional tests introduced.
- Includes early demos of parts of the schema system, and early demos of code generation.
- Mostly a checkpoint before beginning v0.0.3, which involved many large API reshapings.

### v0.0.1

- Our very first tag!
- The central `Node` and `NodeBuilder` interfaces are already established, as is `Link`, `Loader`, and so forth.
  You can already build generic data handling using IPLD Data Model concepts with these core interfaces.
- Selectors and traversals are available.
- Codecs for dag-cbor and dag-json are batteries-included in the repo.
- There was quite a lot of work done before we even started tagging releases :)
