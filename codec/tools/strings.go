package codectools

import (
	"bytes"
	"io"
	"unicode/utf8"
)

var hex = "0123456789abcdef"

// EscapeString is as per WriteEscapedString, but returns the result rather than pushing it into an io.Writer.
func EscapeString(s string) string {
	var buf bytes.Buffer
	WriteEscapedString(s, &buf)
	return buf.String()
}

func QuoteEscapeString(s string) string {
	var buf bytes.Buffer
	WriteQuotedEscapedString(s, &buf)
	return buf.String()
}

func WriteQuotedEscapedString(s string, w io.Writer) error {
	if _, err := w.Write([]byte{'"'}); err != nil {
		return err
	}
	if err := WriteEscapedString(s, w); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'"'}); err != nil {
		return err
	}
	return nil
}

// WriteEscapedString emits a string, surrounded by quotation marks, that's escaped in what's probably "DWIM"-compliant.
//
// More specifically:
//
//   - The output always begins and ends in `"`.
//   - Any `"` marks in the string are escaped as `\"`.
//   - The common escape sequences for familiar unprintables is used: `\r`, `\n`, `\t`, etc.
//   - Unprintables are encoded as `\x__`, where `__` is replaced by the two lowercase hex characters for that byte.
//   - UTF-8 runes will be printed as they are (no escaping; e.g., å is å).
//   - UTF invalid sequences will be encoded as if they were unprintables.
//   - Any `\` characters in the string are escaped as `\\`.
//
// The resulting output will be a single line, and can be parsed back to the exact same string as the original.
// The resulting output is parsable as a JSON string if it does not include any `\x__` sequences,
// or if your json parser supports hex escapes.
//
func WriteEscapedString(s string, w io.Writer) error {
	scratch := [4]byte{'\\', 0, 0, 0}
	start := 0 // 'start' will track the last place we flushed; 'i' will track what byte offset we're examining.
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != 127 {
				i++
				continue
			}
			if start < i { // flush what we've scanned but not flushed leading up to this
				if _, err := w.Write([]byte(s[start:i])); err != nil {
					return err
				}
			}
			switch b {
			case '\\', '"':
				scratch[1] = b
				if _, err := w.Write(scratch[0:2]); err != nil {
					return err
				}
			case '\n':
				scratch[1] = 'n'
				if _, err := w.Write(scratch[0:2]); err != nil {
					return err
				}
			case '\r':
				scratch[1] = 'r'
				if _, err := w.Write(scratch[0:2]); err != nil {
					return err
				}
			case '\t':
				scratch[1] = 't'
				if _, err := w.Write(scratch[0:2]); err != nil {
					return err
				}
			case 127:
				fallthrough
			default:
				// This encodes DEL and bytes < 0x20 except for \t, \n and \r.
				scratch[1] = 'x'
				scratch[2] = hex[b>>4]
				scratch[3] = hex[b&0xF]
				if _, err := w.Write(scratch[0:4]); err != nil {
					return err
				}
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i { // flush what we've scanned but not flushed leading up to this
				if _, err := w.Write([]byte(s[start:i])); err != nil {
					return err
				}
			}
			// Encode the next byte as hex, then continue;
			//  we'll hope the stream resynchronizes back to UTF with the next byte.
			scratch[1] = 'x'
			scratch[2] = hex[s[i]>>4]
			scratch[3] = hex[s[i]&0xF]
			if _, err := w.Write(scratch[0:4]); err != nil {
				return err
			}
			i++
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		if _, err := w.Write([]byte(s[start:])); err != nil {
			return err
		}
	}
	return nil
}
