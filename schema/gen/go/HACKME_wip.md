misc notes during refresh
=========================

(This document will be deleted as soon as this update work cycle is complete;
the contents are a bit "top-of-the-head".)

corners needing mention in docs
-------------------------------

// todo:reorg: these are both **semantics** so they belong in ... schema package docs?  and godocs at that, not readmes.

- iterating a type-level node with optional fields can yield the field and a maybe containing absent.
	- ...which is funny because if you feed that into a type-level builder, it doesn't like that absent.  it wants you to not feed that field.
	- alternatively, accepting explicit puts of absent: worse: we'd have to keep state track that it's been put, but to none, and reject future puts.
		- REVIEW: maybe this isn't as bad as first thought.  I think we end up with that state bit anyway.
		- REVIEW: maybe we should turn this on its head entirely: would it be clearer and more consistent if building a struct without explicitly assigning undef to any optional fields is actually *rejected* when using the type-level assemblers?
	- alternatively, not yielding them on iterate: worse: generic printer for structs would end up not reporting fields, and that would be both wrong and hard to hack your way out of without writing a metric ton more code that inspects the type info, which would ruin the point of the monomorphic methods in the first place to a much higher degree than this need to handle undefined/absent does.

- There is no promise of nice-to-read errors if you over-hold child assemblers beyond their valid lifespan.
	- We **do** care about making things fail hard and fast rather than potentially leak inappropriate mutability.
	- We do **not** care about making these errors pretty (it's high cost to do so, and code that hits this path is almost certainly statically (and hopefully fairly obviously) wrong).

