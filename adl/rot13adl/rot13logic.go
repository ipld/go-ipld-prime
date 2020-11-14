package rot13adl

import (
	"strings"
)

var replaceTable = []string{
	"A", "N",
	"B", "O",
	"C", "P",
	"D", "Q",
	"E", "R",
	"F", "S",
	"G", "T",
	"H", "U",
	"I", "V",
	"J", "W",
	"K", "X",
	"L", "Y",
	"M", "Z",
	"N", "A",
	"O", "B",
	"P", "C",
	"Q", "D",
	"R", "E",
	"S", "F",
	"T", "G",
	"U", "H",
	"V", "I",
	"W", "J",
	"X", "K",
	"Y", "L",
	"Z", "M",
	"a", "n",
	"b", "o",
	"c", "p",
	"d", "q",
	"e", "r",
	"f", "s",
	"g", "t",
	"h", "u",
	"i", "v",
	"j", "w",
	"k", "x",
	"l", "y",
	"m", "z",
	"n", "a",
	"o", "b",
	"p", "c",
	"q", "d",
	"r", "e",
	"s", "f",
	"t", "g",
	"u", "h",
	"v", "i",
	"w", "j",
	"x", "k",
	"y", "l",
	"z", "m",
}
var unreplaceTable = func() []string {
	v := make([]string, len(replaceTable))
	for i := 0; i < len(replaceTable); i += 2 {
		v[i] = replaceTable[i+1]
		v[i+1] = replaceTable[i]
	}
	return v
}()

// rotate transforms from the logical content to the raw content.
func rotate(s string) string {
	return strings.NewReplacer(replaceTable...).Replace(s)
}

// unrotate transforms from the raw content to the logical content.
func unrotate(s string) string {
	return strings.NewReplacer(unreplaceTable...).Replace(s)
}
