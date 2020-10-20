package codectools

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
)

type Token struct {
	Kind TokenKind

	Length int       // Present for MapOpen or ListOpen.  May be -1 for "unknown" (e.g. a json tokenizer will yield this).
	Bool   bool      // Value.  Union: only has meaning if Kind is TokenKind_Bool.
	Int    int64     // Value.  Union: only has meaning if Kind is TokenKind_Int.
	Float  float64   // Value.  Union: only has meaning if Kind is TokenKind_Float.
	Str    string    // Value.  Union: only has meaning if Kind is TokenKind_String.  ('Str' rather than 'String' to avoid collision with method.)
	Bytes  []byte    // Value.  Union: only has meaning if Kind is TokenKind_Bytes.
	Link   ipld.Link // Value.  Union: only has meaning if Kind is TokenKind_Link.

	Node ipld.Node // Direct pointer to the original data, if this token is used to communicate data during a walk of existing in-memory data.  Absent when token is being used during deserialization.

	// TODO: position info?  We seem to want this basically everywhere the token goes, so it might as well just live here.
	//  Putting this position info into the token would require writing those fields many times, though;
	//   hopefully we can also use them as the primary accounting position then, or else this might be problematic for speed.
}

func (tk Token) String() string {
	switch tk.Kind {
	case TokenKind_MapOpen:
		return fmt.Sprintf("<%c:%d>", tk.Kind, tk.Length)
	case TokenKind_MapClose:
		return fmt.Sprintf("<%c>", tk.Kind)
	case TokenKind_ListOpen:
		return fmt.Sprintf("<%c:%d>", tk.Kind, tk.Length)
	case TokenKind_ListClose:
		return fmt.Sprintf("<%c>", tk.Kind)
	case TokenKind_Null:
		return fmt.Sprintf("<%c>", tk.Kind)
	case TokenKind_Bool:
		return fmt.Sprintf("<%c:%v>", tk.Kind, tk.Bool)
	case TokenKind_Int:
		return fmt.Sprintf("<%c:%v>", tk.Kind, tk.Int)
	case TokenKind_Float:
		return fmt.Sprintf("<%c:%v>", tk.Kind, tk.Float)
	case TokenKind_String:
		return fmt.Sprintf("<%c:%q>", tk.Kind, tk.Str)
	case TokenKind_Bytes:
		return fmt.Sprintf("<%c:%x>", tk.Kind, tk.Bytes)
	case TokenKind_Link:
		return fmt.Sprintf("<%c:%v>", tk.Kind, tk.Link)
	default:
		return "<INVALID>"
	}
}

type TokenKind uint8

const (
	TokenKind_MapOpen   = '{'
	TokenKind_MapClose  = '}'
	TokenKind_ListOpen  = '['
	TokenKind_ListClose = ']'
	TokenKind_Null      = '0'
	TokenKind_Bool      = 'b'
	TokenKind_Int       = 'i'
	TokenKind_Float     = 'f'
	TokenKind_String    = 's'
	TokenKind_Bytes     = 'x'
	TokenKind_Link      = '/'
)

type ErrMalformedTokenSequence struct {
	Detail string
}

func (e ErrMalformedTokenSequence) Error() string {
	return "malformed token sequence: " + e.Detail
}
