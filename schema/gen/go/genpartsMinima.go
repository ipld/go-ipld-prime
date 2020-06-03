package gengo

import (
	"fmt"
	"io"

	wish "github.com/warpfork/go-wish"
)

// EmitInternalEnums creates a file with enum types used internally.
// For example, the state machine values used in map and list builders.
// These always need to exist exactly once in each package created by codegen.
func EmitInternalEnums(packageName string, w io.Writer) {
	fmt.Fprint(w, wish.Dedent(`
		package `+packageName+`

		`+doNotEditComment+`

		import (
			"github.com/ipld/go-ipld-prime/schema"
		)

	`))

	// The 'Maybe' enum does double-duty in this package as a state machine for assembler completion.
	//
	// The 'Maybe_Absent' value gains the additional semantic of "clear to assign (but not null)"
	//  (which works because if you're *in* a value assembler, "absent" as a final result is already off the table).
	// Additionally, we get a few extra states that we cram into the same area of bits:
	//   - `midvalue` is used by assemblers of recursives to block AssignNull after BeginX.
	//   - `allowNull` is used by parent assemblers when initializing a child assembler to tell the child a transition to Maybe_Null is allowed in this context.
	fmt.Fprint(w, wish.Dedent(`
		const (
			midvalue = schema.Maybe(4)
			allowNull = schema.Maybe(5)
		)

	`))

	fmt.Fprint(w, wish.Dedent(`
		type maState uint8

		const (
			maState_initial     maState = iota
			maState_midKey
			maState_expectValue
			maState_midValue
			maState_finished
		)

		type laState uint8

		const (
			laState_initial  laState = iota
			laState_midValue
			laState_finished
		)
	`))
}
