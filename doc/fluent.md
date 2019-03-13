Fluent APIs
===========

The "fluent" APIs are replicates of most of the core interfaces -- e.g. Node,
NodeBuilder, etc -- which return single values.  This makes things easier
to compose in a functional/point-free style.

Errors in the fluent interfaces are handled by panicking.  These errors are
boxed in a `fluent.Error` type, which can be unwrapped into the original error.
`fluent.Recover` can wrap any function with automatic recovery of these errors,
and returns them back to normal flow.  Thus, we can write large blocks of code
using the fluent APIs, and handle all the errors in one place.  Just as easily,
we can use nested sets of `fluent.Recover` as desired for granular handling.
