package gengo

import (
	"github.com/ipld/go-ipld-prime/schema"
)

func mungeTypeNodeIdent(t schema.Type) string {
	return string(t.Name())
}

// future: something return a "_x__Node" might also make an appearance,
//  which would address the fun detail that a gen'd struct type might not actually implement Node directly.

func mungeTypeNodebuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__NodeBuilder"
}

func mungeTypeReprNodeIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__Repr"
}

func mungeTypeReprNodebuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ReprBuilder"
}

// itr

// MapBuilder

// ListBuilder

// reprItr

// ReprMapBuilder

// ReprListBuilder

// okay these are getting a little out of hand given how formulaic they are.
// the repr ones could all be glued onto the end of mungeTypeReprNodeIdent safely;
//  problem with that is, the pattern doesn't always hold:
//   things that are on the main node can't, because that symbol is exported.
