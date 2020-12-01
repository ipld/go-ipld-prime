package jsontoken

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/ipld/go-ipld-prime/codec/codectools/scratch"
)

// License note: the string and numeric parsers here borrow
// heavily from the golang stdlib json parser scanner.
// That code is originally Copyright 2010 The Go Authors,
// and is governed by a BSD-style license.

// DecodeString will attempt to decode data in the format of a JSON string from the reader.
// If the first byte read is not `"`, it is not a string at all, and an error is returned.
// Any other parse errors of json strings also result in error.
func DecodeString(r *scratch.Reader) (string, error) {
	// Check that this actually begins like a string.
	majorByte, err := r.Readn1()
	if err != nil {
		return "", err
	}
	if majorByte != '"' {
		return "", fmt.Errorf("not a string: strings must begin with '\"', not %q", majorByte)
	}
	// Decode the string body.
	s, err := DecodeStringBody(r)
	if err != nil {
		return "", err
	}
	// Swallow the trailing `"` again (which DecodeStringBody has insured we have).
	r.Readn1()
	return s, nil
}

// DecodeStringBody will attempt to decode data in the format of a JSON string from the reader,
// except it assumes that the leading `"` has already been consumed,
// and will similarly leave the trailing `"` unread (although it will check for its presence).
//
// Implementation note: you'll find that this method is used in the Decoder's implementation,
// while DecodeString is actually not.  This is because when doing a whole document parse,
// the leading `"` is always already consumed because it's how we discovered it's time to parse a string.
func DecodeStringBody(r *scratch.Reader) (string, error) {
	// First `"` is presumed already eaten.
	// Start tracking the byte slice; real string starts here.
	r.Track()
	// Scan until scanner tells us end of string.
	for step := strscan_normal; step != nil; {
		majorByte, err := r.Readn1()
		if err != nil {
			return "", err
		}
		step, err = step(majorByte)
		if err != nil {
			return "", err
		}
	}
	// Unread one.  The scan loop consumed the trailing quote already,
	// which we don't want to pass onto the parser.
	r.Unreadn1()
	// Parse!
	s, ok := parseString(r.StopTrack())
	if !ok {
		panic("string parse failed") // this is a sanity check; our scan phase should've already excluded any data that would cause this.
	}
	return string(s), nil
}

// strscanStep steps are applied over the data to find how long the string is.
// A nil step func is returned to indicate the string is done.
// Actually parsing the string is done by 'parseString()'.
type strscanStep func(c byte) (strscanStep, error)

// The default string scanning step state.  Starts here.
func strscan_normal(c byte) (strscanStep, error) {
	if c == '"' { // done!
		return nil, nil
	}
	if c == '\\' {
		return strscan_esc, nil
	}
	if c < 0x20 { // Unprintable bytes are invalid in a json string.
		return nil, fmt.Errorf("invalid unprintable byte in string literal: 0x%x", c)
	}
	return strscan_normal, nil
}

// "esc" is the state after reading `"\` during a quoted string.
func strscan_esc(c byte) (strscanStep, error) {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		return strscan_normal, nil
	case 'u':
		return strscan_escU, nil
	}
	return nil, fmt.Errorf("invalid byte in string escape sequence: 0x%x", c)
}

// "escU" is the state after reading `"\u` during a quoted string.
func strscan_escU(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU1, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU1" is the state after reading `"\u1` during a quoted string.
func strscan_escU1(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU12, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU12" is the state after reading `"\u12` during a quoted string.
func strscan_escU12(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_escU123, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// "escU123" is the state after reading `"\u123` during a quoted string.
func strscan_escU123(c byte) (strscanStep, error) {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		return strscan_normal, nil
	}
	return nil, fmt.Errorf("invalid byte in \\u hexadecimal character escape: 0x%x", c)
}

// Convert a json serial byte sequence that is a complete string body (i.e., quotes from the outside excluded)
// into a natural byte sequence (escapes, etc, are processed).
//
// The given slice should already be the right length.
// A blithe false for 'ok' is returned if the data is in any way malformed.
//
// FUTURE: this is native JSON string parsing, and not as strict as DAG-JSON should be.
//
//   - this does not implement UTF8-C8 unescpaing; we may want to do so.
//   - this transforms invalid surrogates coming from escape sequences into uFFFD; we probably shouldn't.
//   - this transforms any non-UTF-8 bytes into uFFFD rather than erroring; we might want to think twice about that.
//   - this parses `\u` escape sequences at all, while also allowing UTF8 chars of the same content; we might want to reject variations.
//
// It might be desirable to implement these stricter rules as configurable.
func parseString(s []byte) (t []byte, ok bool) {
	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room?  Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	r, err := strconv.ParseUint(string(s[2:6]), 16, 64)
	if err != nil {
		return -1
	}
	return rune(r)
}
