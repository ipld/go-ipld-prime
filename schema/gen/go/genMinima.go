package gengo

import (
	"fmt"
	"io"

	wish "github.com/warpfork/go-wish"
)

// EmitInternalEnums creates a file with enum types used internal.
// For example, the state machine values used in map and list builders.
func EmitInternalEnums(packageName string, w io.Writer) {
	fmt.Fprint(w, wish.Dedent(`
		package `+packageName+`

		// Code generated go-ipld-prime DO NOT EDIT.

		type maState uint8

		const (
			maState_initial     maState = iota
			maState_midKey
			maState_expectValue
			maState_midValue
			maState_finished
		)
	`))
}
