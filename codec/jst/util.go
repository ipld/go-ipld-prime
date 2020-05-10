package jst

import (
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

func max(a, b int) int { // honestly golang
	if a > b {
		return a
	}
	return b
}

func mustFirstKeyAsString(mapNode ipld.Node) string {
	itr := mapNode.MapIterator()
	k, _, err := itr.Next()
	if err != nil {
		panic(err)
	}
	ks, err := k.AsString()
	if err != nil {
		panic(err)
	}
	return ks
}

func indexOf(list []columnName, cn columnName) int {
	for i, v := range list {
		if v == cn {
			return i
		}
	}
	return -1
}

type codecAborted struct {
	p ipld.Path
	e error
}

func (e codecAborted) Error() string {
	return fmt.Sprintf("codec aborted: %s", e)
}

// attempt to record where the codec efforts encountered an error.
// this is currently a bit slapdash.
// a better system might also count byte and line offsets in the serial form (but this would require also integrating with the serial code pretty closely).
// a spectacularly general system might be ready to count serial offsets *twice* *or* path offsets *twice* (but this is simply waiting for someone's enthusiasm).
func recordErrorPosition(ctx *state, e error) error {
	return codecAborted{ipld.Path{ /*TODO*/ }, e}
}

// not yet used, but you'd probably want this for better error position purposes.
type writerNanny struct {
	totalOffset  int
	lineOffset   int
	offsetInLine int
	w            io.Writer
}
