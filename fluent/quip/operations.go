package quip

import (
	"github.com/ipld/go-ipld-prime"
)

func CopyRange(e *error, la ipld.ListAssembler, src ipld.Node, start, end int64) {
	if *e != nil {
		return
	}
	if start >= src.Length() {
		return
	}
	if end < 0 {
		end = src.Length()
	}
	if end < start {
		return
	}
	for i := start; i < end; i++ {
		n, err := src.LookupByIndex(i)
		if err != nil {
			*e = err
			return
		}
		if err := la.AssembleValue().AssignNode(n); err != nil {
			*e = err
			return
		}
	}
	return
}
