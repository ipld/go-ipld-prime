package rot13adl

import (
	"strings"
)

func rot13(r rune) rune {
	if r >= 'a' && r <= 'z' {
		if r >= 'm' {
			return r - 13
		} else {
			return r + 13
		}
	} else if r >= 'A' && r <= 'Z' {
		if r >= 'M' {
			return r - 13
		} else {
			return r + 13
		}
	}
	return r
}

// rotate transforms from the logical content to the raw content.
func rotate(s string) string {
	return strings.Map(rot13, s)
}

// unrotate transforms from the raw content to the logical content.
func unrotate(s string) string {
	return strings.Map(rot13, s)
}
