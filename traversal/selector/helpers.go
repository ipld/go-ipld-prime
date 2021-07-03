package selector

import (
	"strings"

	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/basic"
)

// REVIEW: I feel sketchy about CompileJSONSelector because it brings in json as a transitive dependency.
// This has consequences as far as dragging go-cid in because our json shares code with dagjson...
// and then go-cid ends up in the transitives of the traversal package, since traversal depends on selector.
// That's... not great.
// On the other hand: trying to keep up that paper wall is creating a lot of strife.  Maybe it's time to give up.

func CompileJSONSelector(jsonStr string) (Selector, error) {
	na := basicnode.Prototype.Any.NewBuilder()
	if err := json.Decode(na, strings.NewReader(jsonStr)); err != nil {
		return nil, err
	}
	if s, err := CompileSelector(na.Build()); err != nil {
		return nil, err
	} else {
		return s, nil
	}
}

func must(s Selector, e error) Selector {
	if e != nil {
		panic(e)
	}
	return s
}

var CommonSelector_MatchPoint = must(CompileJSONSelector(`{".":{}}`))
