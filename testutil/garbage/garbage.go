package garbage

import (
	"math"
	mathrand "math/rand"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/must"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/multiformats/go-multihash"
)

type Options struct {
	initialWeights map[datamodel.Kind]int
	weights        map[datamodel.Kind]int
	blockSize      uint64
}

type generator func(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node)

type hasher struct {
	code   uint64
	length int
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`~!@#$%^&*()-_=+[]{}|\\:;'\",.<>?/ \t\n☺💩"

var (
	codecs     = []uint64{0x55, 0x70, 0x71, 0x0129}
	hashes     = []hasher{{0x12, 256}, {0x16, 256}, {0x1b, 256}, {0xb220, 256}, {0x13, 512}, {0x15, 384}, {0x14, 512}}
	kinds      = append(datamodel.KindSet_Scalar, datamodel.KindSet_Recursive...)
	runes      = []rune(charset)
	generators map[datamodel.Kind]generator
)

// Generate produces random Nodes which can be useful for testing and benchmarking. By default, the
// Nodes produced are relatively small, averaging near the 1024 byte range when encoded
// (very roughly, with a wide spread).
//
// Options can be used to adjust the average size and weights of occurances of different kinds
// within the complete Node graph.
//
// Care should be taken when using a random source to generate garbage for testing purposes, that
// the randomness is stable across test runs, or a seed is captured in such a way that a failure
// can be reproduced (e.g. by printing it to stdout during the test run so it can be captured in
// CI for a failure).
func Generate(rand *mathrand.Rand, opts ...Option) datamodel.Node {
	options := applyOptions(opts...)
	_, n := generate(rand, options.blockSize, options)
	return n
}

func generate(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	weights := opts.weights
	if opts.initialWeights != nil {
		weights = opts.initialWeights
		opts = Options{weights: opts.weights}
	}
	totWeight := 0
	for _, kind := range kinds {
		totWeight += weights[kind]
	}
	r := rand.Float64() * float64(totWeight)
	var wacc int
	for _, kind := range kinds {
		wacc += weights[kind]
		if float64(wacc) >= r {
			return generators[kind](rand, count, opts)
		}
	}
	panic("bad options")
}

func rndSize(rand *mathrand.Rand, bias uint64) uint64 {
	if bias == 0 {
		panic("size shouldn't be zero")
	}
	mean := float64(bias)
	stdev := mean / 10
	for {
		s := math.Abs(rand.NormFloat64())*stdev + mean
		if s >= 1 {
			return uint64(s)
		}
	}
}

func rndRune(rand *mathrand.Rand) rune {
	return runes[rand.Intn(len(runes))]
}

func listGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	len := rndSize(rand, 10)
	lb := basicnode.Prototype.List.NewBuilder()
	la, err := lb.BeginList(int64(len))
	if err != nil {
		panic(err)
	}
	size := uint64(0)
	for i := uint64(0); i < len && size < count; i++ {
		c, n := generate(rand, count-size, opts)
		err := la.AssembleValue().AssignNode(n)
		if err != nil {
			panic(err)
		}
		size += c
	}
	err = la.Finish()
	if err != nil {
		panic(err)
	}
	return size, lb.Build()
}

func mapGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	length := rndSize(rand, 10)
	mb := basicnode.Prototype.Map.NewBuilder()
	ma, err := mb.BeginMap(int64(length))
	if err != nil {
		panic(err)
	}
	size := uint64(0)
	keys := make(map[string]struct{})
	for i := uint64(0); i < length && size < count; i++ {
		var key string
		for {
			c, k := stringGenerator(rand, 5, opts)
			key = must.String(k)
			if _, ok := keys[key]; !ok && len(key) > 0 {
				keys[key] = struct{}{}
				size += c
				break
			}
		}
		sz := count - size
		if size >= count { // the case where we've blown our budget already on the key
			sz = 5
		}
		c, value := generate(rand, sz, opts)
		size += c
		err := ma.AssembleKey().AssignString(key)
		if err != nil {
			panic(err)
		}
		err = ma.AssembleValue().AssignNode(value)
		if err != nil {
			panic(err)
		}
	}
	err = ma.Finish()
	if err != nil {
		panic(err)
	}
	return size, mb.Build()
}

func stringGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	len := rndSize(rand, count/2+1)
	sb := strings.Builder{}
	for i := uint64(0); i < len; i++ {
		sb.WriteRune(rndRune(rand))
	}
	return len, basicnode.NewString(sb.String())
}

func bytesGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	len := rndSize(rand, count/2+1)
	ba := make([]byte, len)
	_, err := rand.Read(ba)
	if err != nil {
		panic(err)
	}
	return len, basicnode.NewBytes(ba)
}

func boolGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	return 0, basicnode.NewBool(rand.Float64() > 0.5)
}

func intGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	i := rand.Int63()
	if rand.Float64() > 0.5 {
		i = -i
	}
	return 0, basicnode.NewInt(i)
}

func floatGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	return 0, basicnode.NewFloat(math.Tan((rand.Float64() - 0.5) * math.Pi))
}

func nullGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	return 0, datamodel.Null
}

func linkGenerator(rand *mathrand.Rand, count uint64, opts Options) (uint64, datamodel.Node) {
	hasher := hashes[rand.Intn(len(hashes))]
	codec := codecs[rand.Intn(len(codecs))]
	ba := make([]byte, hasher.length/8)
	rand.Read(ba)
	mh, err := multihash.Encode(ba, hasher.code)
	if err != nil {
		panic(err)
	}
	return uint64(hasher.length / 8), basicnode.NewLink(cidlink.Link{Cid: cid.NewCidV1(codec, mh)})
}

type Option func(*Options)

func applyOptions(opt ...Option) Options {
	opts := Options{
		blockSize:      1024,
		initialWeights: DefaultInitialWeights(),
		weights:        DefaultWeights(),
	}
	for _, o := range opt {
		o(&opts)
	}
	return opts
}

// DefaultInitialWeights provides the default map of weights that can be
// overridden by the InitialWeights option. The default is an equal weighting
// of 1 for every scalar kind and 10 for the recursive kinds.
func DefaultInitialWeights() map[datamodel.Kind]int {
	return map[datamodel.Kind]int{
		datamodel.Kind_List:   10,
		datamodel.Kind_Map:    10,
		datamodel.Kind_Bool:   1,
		datamodel.Kind_Bytes:  1,
		datamodel.Kind_Float:  1,
		datamodel.Kind_Int:    1,
		datamodel.Kind_Link:   1,
		datamodel.Kind_Null:   1,
		datamodel.Kind_String: 1,
	}
}

// DefaultWeights provides the default map of weights that can be overridden by
// the Weights option. The default is an equal weighting of 1 for every kind.
func DefaultWeights() map[datamodel.Kind]int {
	return map[datamodel.Kind]int{
		datamodel.Kind_List:   1,
		datamodel.Kind_Map:    1,
		datamodel.Kind_Bool:   1,
		datamodel.Kind_Bytes:  1,
		datamodel.Kind_Float:  1,
		datamodel.Kind_Int:    1,
		datamodel.Kind_Link:   1,
		datamodel.Kind_Null:   1,
		datamodel.Kind_String: 1,
	}
}

// InitialWeights sets a per-kind weighting for the root node. That is, the weights
// set here will determine the liklihood of the returned Node's direct .Kind().
// These weights are ignored after the top-level Node (for recursive kinds,
// obviously for scalar kinds there is only a top-level Node).
//
// The default initial weights bias toward Map and List kinds, by a ratio of
// 10:1—i.e. the recursive kinds are more likely to appear at the top-level.
func InitialWeights(initialWeights map[datamodel.Kind]int) Option {
	return func(o *Options) {
		o.initialWeights = initialWeights
	}
}

// Weights sets a per-kind weighting for nodes appearing throughout the returned
// graph. When assembling a graph, these weights determine the liklihood that
// a given kind will be selected for that node.
//
// A weight of 0 will turn that kind off entirely. So, for example, if you
// wanted output data with no maps or bytes, then set both of those weights to
// zero, leaving the rest >0 and do the same for InitialWeights.
//
// The default weights are set to 1—i.e. there is an equal liklihood that any of
// the valid kinds will be selected for any point in the graph.
//
// This option is overridden by InitialWeights (which also has a default even
// if not set explicitly) for the top-level node.
func Weights(weights map[datamodel.Kind]int) Option {
	return func(o *Options) {
		o.weights = weights
	}
}

// TargetBlockSize sets a very rough bias in number of bytes that the resulting
// Node may consume when encoded (i.e. the block size). This is a very
// approximate measure, but over enough repeated Generate() calls, the resulting
// Nodes, once encoded, should have a median that is somewhere in this vicinity.
//
// The default target block size is 1024. This should be tuned in accordance with
// the anticipated average block size of the system under test.
func TargetBlockSize(blockSize uint64) Option {
	return func(o *Options) {
		o.blockSize = blockSize
	}
}

func init() {
	// can't be declared statically because of some cycles through list & map to generate()
	generators = map[datamodel.Kind]generator{
		datamodel.Kind_List:   listGenerator,
		datamodel.Kind_Map:    mapGenerator,
		datamodel.Kind_String: stringGenerator,
		datamodel.Kind_Bytes:  bytesGenerator,
		datamodel.Kind_Bool:   boolGenerator,
		datamodel.Kind_Int:    intGenerator,
		datamodel.Kind_Float:  floatGenerator,
		datamodel.Kind_Null:   nullGenerator,
		datamodel.Kind_Link:   linkGenerator,
	}
}
