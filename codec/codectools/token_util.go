package codectools

import (
	"strings"
)

// StringifyTokenSequence is utility function often handy for testing.
// (Doing a diff on strings of tokens gives very good reports for minimal effort.)
func StringifyTokenSequence(seq []Token) string {
	var sb strings.Builder
	for _, tk := range seq {
		sb.WriteString(tk.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}
