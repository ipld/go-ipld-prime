package codectools

import (
	"fmt"

	"github.com/ipld/go-ipld-prime"
)

type Token struct {
	Kind TokenKind

	Length int64     // Present for MapOpen or ListOpen.  May be -1 for "unknown" (e.g. a json tokenizer will yield this).
	Bool   bool      // Value.  Union: only has meaning if Kind is TokenKind_Bool.
	Int    int64     // Value.  Union: only has meaning if Kind is TokenKind_Int.
	Float  float64   // Value.  Union: only has meaning if Kind is TokenKind_Float.
	Str    string    // Value.  Union: only has meaning if Kind is TokenKind_String.  ('Str' rather than 'String' to avoid collision with method.)
	Bytes  []byte    // Value.  Union: only has meaning if Kind is TokenKind_Bytes.
	Link   ipld.Link // Value.  Union: only has meaning if Kind is TokenKind_Link.

	Node ipld.Node // Direct pointer to the original data, if this token is used to communicate data during a walk of existing in-memory data.  Absent when token is being used during deserialization.

	// The following fields all track position and progress:
	// (These may be useful to copy into any error messages if errors arise.)
	// (Implementations may assume token reuse and treat these as state keeping;
	// you may experience position accounting accuracy problems if *not* reusing tokens or if zeroing these fields.)

	pth          []ipld.PathSegment // Set by token producers (whether marshallers or deserializers) to track logical position.
	offset       int64              // Set by deserializers (for both textual or binary formats alike) to track progress.
	lineOffset   int64              // Set by deserializers that work with textual data.  May be ignored by binary deserializers.
	columnOffset int64              // Set by deserializers that work with textual data.  May be ignored by binary deserializers.
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
	TokenKind_MapOpen   TokenKind = '{'
	TokenKind_MapClose  TokenKind = '}'
	TokenKind_ListOpen  TokenKind = '['
	TokenKind_ListClose TokenKind = ']'
	TokenKind_Null      TokenKind = '0'
	TokenKind_Bool      TokenKind = 'b'
	TokenKind_Int       TokenKind = 'i'
	TokenKind_Float     TokenKind = 'f'
	TokenKind_String    TokenKind = 's'
	TokenKind_Bytes     TokenKind = 'x'
	TokenKind_Link      TokenKind = '/'
)

type ErrMalformedTokenSequence struct {
	Detail string
}

func (e ErrMalformedTokenSequence) Error() string {
	return "malformed token sequence: " + e.Detail
}
