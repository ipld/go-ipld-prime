package gengo

import (
	"github.com/ipld/go-ipld-prime/schema"
)

func mungeTypeNodeIdent(t schema.Type) string {
	return string(t.Name())
}

// future: something return a "_x__Node" might also make an appearance,
//  which would address the fun detail that a gen'd struct type might not actually implement Node directly.

func mungeTypeNodeItrIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__Itr"
}

func mungeTypeNodebuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__NodeBuilder"
}

func mungeTypeNodeMapBuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__MapBuilder"
}

func mungeTypeNodeListBuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ListBuilder"
}

func mungeTypeReprNodeIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__Repr"
}

func mungeTypeReprNodeItrIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ReprItr"
}

func mungeTypeReprNodebuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ReprBuilder"
}

func mungeTypeReprNodeMapBuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ReprMapBuilder"
}

func mungeTypeReprNodeListBuilderIdent(t schema.Type) string {
	return "_" + string(t.Name()) + "__ReprListBuilder"
}

func mungeNodebuilderConstructorIdent(t schema.Type) string {
	return string(t.Name()) + "__NodeBuilder"
}

func mungeReprNodebuilderConstructorIdent(t schema.Type) string {
	return string(t.Name()) + "__ReprBuilder"
}
