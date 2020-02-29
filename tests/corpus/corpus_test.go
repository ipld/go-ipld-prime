package corpus

import (
	"encoding/json"
	"testing"

	"github.com/ipld/go-ipld-prime/must"
)

/*
	Sanity check that our corpuses are actually correct JSON in this package,
	before we start letting other packages find that out for us the hard way.
*/

func TestCorpusValidity(t *testing.T) {
	must.True(json.Valid([]byte(MapNStrInt(0))))
	must.True(json.Valid([]byte(MapNStrInt(1))))
	must.True(json.Valid([]byte(MapNStrInt(2))))
	must.True(json.Valid([]byte(MapNStrMap3StrInt(0))))
	must.True(json.Valid([]byte(MapNStrMap3StrInt(1))))
	must.True(json.Valid([]byte(MapNStrMap3StrInt(2))))
}
