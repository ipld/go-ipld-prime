package testutil

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDedent(t *testing.T) {
	for _, tr := range []struct{ a, b string }{
		{"", ""},
		{"\t", ""},
		{"\t\t", ""},
		{"\n", ""},
		{"\n\t", ""},
		{"\n\t\t", ""},
		{"\n\n", "\n"},
		{"\n\t\n", "\n"},
		{"\n\t\t\n", "\n"},
		{"\n\n", "\n"},
		{"\n\n\t", "\n\t"},
		{"\n\n\t\t", "\n\t\t"},
		{"a\nb\n\tc\n", "a\nb\n\tc\n"},
		{"\ta\nb\n\tc\n", "a\nb\nc\n"},
		{"\t\ta\nb\n\tc\n", "a\nb\nc\n"},
		{"\ta\n\t\tb\n\tc\n", "a\n\tb\nc\n"},
		{"\ta\n\t\t\tb\n\tc\n", "a\n\t\tb\nc\n"},
		{"\n\t\t\ta\n\t\tb\n\t\t\t\n\t\t\t\tc\n\t\t", "a\nb\n\n\tc\n"},
	} {
		actual := Dedent(tr.a)
		qt.Assert(t, actual, qt.Equals, tr.b)
	}
}
