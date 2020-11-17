package codectools

import (
	"strings"
)

// Normalize sets any value in the token to its zero value if it's not applicable for the token's kind.
// E.g., if the token kind is string, the float, bytes, and etc fields are all zero'd.
// Path and offset progress information is left unmodified.
// This is sometimes helpful in writing test fixtures and equality assertions.
func (tk *Token) Normalize() {
	if tk.Kind != TokenKind_MapOpen && tk.Kind != TokenKind_ListOpen {
		tk.Length = 0
	}
	if tk.Kind != TokenKind_Bool {
		tk.Bool = false
	}
	if tk.Kind != TokenKind_Int {
		tk.Int = 0
	}
	if tk.Kind != TokenKind_Float {
		tk.Float = 0
	}
	if tk.Kind != TokenKind_String {
		tk.Str = ""
	}
	if tk.Kind != TokenKind_Bytes {
		tk.Bytes = nil
	}
	if tk.Kind != TokenKind_Link {
		tk.Link = nil
	}
}

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
