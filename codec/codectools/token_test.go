package codectools

import (
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
)

var tokenFixtures = []struct {
	value    ipld.Node
	sequence []Token
}{
	{
		value: must.Node(fluent.Reflect(basicnode.Prototype.Any,
			"a scalar",
		)),
		sequence: []Token{
			{Kind: TokenKind_String, Str: "a scalar"},
		},
	},
	{
		value: must.Node(fluent.Reflect(basicnode.Prototype.Any,
			map[string]interface{}{
				"a": "b",
				"c": "d",
			},
		)),
		sequence: []Token{
			{Kind: TokenKind_MapOpen, Length: 2},
			/**/ {Kind: TokenKind_String, Str: "a"}, {Kind: TokenKind_String, Str: "b"},
			/**/ {Kind: TokenKind_String, Str: "c"}, {Kind: TokenKind_String, Str: "d"},
			{Kind: TokenKind_MapClose},
		},
	},
	{
		value: must.Node(fluent.Reflect(basicnode.Prototype.Any,
			map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"c": "d",
				},
			},
		)),
		sequence: []Token{
			{Kind: TokenKind_MapOpen, Length: 2},
			/**/ {Kind: TokenKind_String, Str: "a"}, {Kind: TokenKind_Int, Int: 1},
			/**/ {Kind: TokenKind_String, Str: "b"}, {Kind: TokenKind_MapOpen, Length: 1},
			/**/ /**/ {Kind: TokenKind_String, Str: "c"}, {Kind: TokenKind_String, Str: "d"},
			/**/ {Kind: TokenKind_MapClose},
			{Kind: TokenKind_MapClose},
		},
	},
	{
		value: must.Node(fluent.Reflect(basicnode.Prototype.Any,
			[]interface{}{
				"a",
				"b",
				"c",
			},
		)),
		sequence: []Token{
			{Kind: TokenKind_ListOpen, Length: 3},
			/**/ {Kind: TokenKind_String, Str: "a"},
			/**/ {Kind: TokenKind_String, Str: "b"},
			/**/ {Kind: TokenKind_String, Str: "c"},
			{Kind: TokenKind_ListClose},
		},
	},
}

// utility function for testing.  Doing a diff on strings of tokens gives very good reports for minimal effort.
func stringifyTokens(seq []Token) string {
	var sb strings.Builder
	for _, tk := range seq {
		sb.WriteString(tk.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}
