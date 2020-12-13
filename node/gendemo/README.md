hackme
======

This package demonstrates what codegen output code looks and acts like.

The main purpose is to benchmark things,
and to provide an easy-to-look-at _thing_ for prospective users
who want to lay eyes on generated code without needing to get up-and-running with the generator themselves.

This package is absolutely _not_ full of general purpose node implementations
that you should use in _any_ application.

The input info for the code generation is in `gen.go` file.
(This'll be extracted to be its own schema file, etc, later --
but you'll have to imagine that part; at present, it's wired directly in code.)
The code generation is triggered by `go:generate` comments in `gen_trigger.go`.
