package codectools

import (
	"testing"
	"unicode/utf8"
)

func TestEscapeString(t *testing.T) {
	for _, tc := range []struct {
		gostr string
		quote string
	}{
		{"", `""`},                      // empty string
		{"a", `"a"`},                    // simple ascii
		{"abc", `"abc"`},                // longer simple ascii
		{"\n", `"\n"`},                  // special control characters
		{"\t\t\r\\\n", `"\t\t\r\\\n"`},  // more special control characters
		{"\a", `"\x07"`},                // golang considers bell noteworthy; this EscapeString doesn't
		{"\x01\x00", `"\x01\x00"`},      // unprintables and null
		{"å", `"å"`},                    // unicode!
		{"åå", `"åå"`},                  // longer unicode!
		{"\xc3\x21", `"\xc3!"`},         // non utf-8!  resynchronizes!
		{"\xe2\x82\x28", `"\xe2\x82("`}, // longer non-utf-8 sequences.  resynchronizes!
		{"\xc3\xe2\x82\xc3\x21", `"\xc3\xe2\x82\xc3!"`}, // multiple sequential non-utf-8 sequences.  resynchronizes!
	} {
		result := QuoteEscapeString(tc.gostr)
		if result != tc.quote {
			t.Errorf("EscapeString fixture mismatch: expected %v, got %v", tc.quote, result)
		} else {
			t.Logf("%q -> `%v`", tc.gostr, result)
		}
	}
}

// TestStringSanity just reasserts some basics about how strings work,
// and what is widely understood as normal things to say about strings.
// It's not a test (it'll never change unless unicode does!) so much as it's documentation:
// unicode can be fiddly, and executable sanity checks make it easier to reason about things.
func TestStringSanity(t *testing.T) {
	for _, tc := range []struct {
		bytes  []byte
		gostr  string
		isUtf8 bool
	}{
		{[]byte{0x41}, "A", true},               // Printables are clear enough.
		{[]byte{0x0}, "\u0000", true},           // Yes, the "\uxxxx" format still unpacks to single bytes when enough of it is zeros.
		{[]byte{0xC3, 0xA5}, "å", true},         // Some unicode values are multibyte, of course.
		{[]byte{0xC3, 0xA5}, "\u00e5", true},    // You can say that same characer as a "\uxxxx" escape sequence.  (It's not at all the same as how you'd say the bytes as hex individually, though.)
		{[]byte{0xC3, 0x21}, "\xc3\x21", false}, // Not all byte sequences are valid unicode: a byte that says "unicode coming" followed by a plain ascii byte is invalid, for example.
		// Not present in this table: what "\xc3\x21" would look like in "\uxxxx" format... because it's not representable.
	} {
		if string(tc.bytes) != tc.gostr {
			t.Errorf("sanity table wrong: %v != %q; %q is %v", tc.bytes, tc.gostr, tc.gostr, []byte(tc.gostr))
		}
		if utf8.ValidString(tc.gostr) != tc.isUtf8 {
			t.Errorf("sanity table wrong: %q isUtf8 is %v", tc.gostr, utf8.ValidString(tc.gostr))
		}
		t.Logf("%v == %q ; utf8 = %v\n", tc.bytes, tc.gostr, tc.isUtf8)
	}
}
