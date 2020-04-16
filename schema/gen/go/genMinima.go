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

		import (
			"github.com/ipld/go-ipld-prime/schema"
		)

		// Code generated go-ipld-prime DO NOT EDIT.

		const (
			// The 'Maybe' enum does double-duty in this package as a state machine for assembler completion.
			// The 'Maybe_Absent' value gains the additional semantic of "clear to assign (but not null)"
			//  (which works because if you're *in* a value assembler, "absent" as a final result is already off the table).
			// Additionally, we get a few extra states that we cram into the same area of bits:

			midvalue = schema.Maybe(4) // used by assemblers of recursives to block AssignNull after BeginX.
			allowNull = schema.Maybe(5) // used by parent assemblers to tell child a transition to Maybe_Null is allowed.
		)

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
