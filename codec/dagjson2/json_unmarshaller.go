package dagjson

import (
	"io"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/codectools"
	"github.com/ipld/go-ipld-prime/codec/dagjson2/token"
)

// Unmarshal reads data from input, parses it as DAG-JSON,
// and unfolds the data into the given NodeAssembler.
//
// The strict interpretation of DAG-JSON is used.
// Use a ReusableMarshaller and set its DecoderConfig if you need
// looser or otherwise customized decoding rules.
//
// This function is the same as the function found for DAG-JSON
// in the default multicodec registry.
func Unmarshal(into ipld.NodeAssembler, input io.Reader) error {
	// FUTURE: consider doing a whole sync.Pool jazz around this.
	r := ReusableUnmarshaller{}
	r.SetDecoderConfig(jsontoken.DecoderConfig{
		AllowDanglingComma:  false,
		AllowWhitespace:     false,
		AllowEscapedUnicode: false,
		ParseUtf8C8:         true,
	})
	r.SetInitialBudget(1 << 20)
	return r.Unmarshal(into, input)
}

// ReusableUnmarshaller has an Unmarshal method, and also supports
// customizable DecoderConfig and resource budgets.
//
// The Unmarshal method may be used repeatedly (although not concurrently).
// Keeping a ReusableUnmarshaller around and using it repeatedly may allow
// the user to amortize some allocations (some internal buffers can be reused).
type ReusableUnmarshaller struct {
	d jsontoken.Decoder

	InitialBudget int
}

func (r *ReusableUnmarshaller) SetDecoderConfig(cfg jsontoken.DecoderConfig) {
	r.d.DecoderConfig = cfg
}
func (r *ReusableUnmarshaller) SetInitialBudget(budget int) {
	r.InitialBudget = budget
}

func (r *ReusableUnmarshaller) Unmarshal(into ipld.NodeAssembler, input io.Reader) error {
	r.d.Init(input)
	return codectools.TokenAssemble(into, r.d.Step, r.InitialBudget)
}
