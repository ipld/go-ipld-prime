package gengo

import (
	"github.com/ipld/go-ipld-prime/schema"
)

type generateKindStruct struct {
	schema.TypeStruct
	// FUTURE: probably some adjunct config data should come with here as well.
	// FUTURE: perhaps both a global one (e.g. output package name) and a per-type one.
}
