// The gendemo package contains some what codegen output code,
// so that it can demonstrate what schema-based codegen looks and acts like.
//
// The main purpose is to benchmark things,
// and to provide an easy-to-look-at _thing_ for prospective users
// who want to lay eyes on generated code without needing to get up-and-running with the generator themselves.
//
// This package is absolutely _not_ full of general purpose node implementations
// that you should use in _any_ application.
//
// The input info for the code generation is in `gen.go` file.
// (This is currently wired directly in code; in the future, the same instructions
// will be extracted to an IPLD Schema file and standard tools will be used to process it.)
// The code generation is triggered by `go:generate` comments in the `doc.go` file.

//go:generate go run gen.go

package gendemo
