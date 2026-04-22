package dagcbor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/ipld/go-ipld-prime/datamodel"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func TestFunBlocks(t *testing.T) {
	t.Run("zero length link", func(t *testing.T) {
		// This fixture has a zero length link -- not even the multibase byte (which dag-cbor insists must be zero) is there.
		buf := strings.NewReader("\x8d\x8d\x97\xd8*@")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.Equals, ErrInvalidMultibase)
	})
	t.Run("fuzz001", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit runtime
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("fuzz002", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9f\x9a\xff000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("fuzz003", func(t *testing.T) {
		// This fixture might cause an overly large allocation if you aren't careful to have resource budgets.
		buf := strings.NewReader("\x9f\x9f\x9f\x9f\x9f\x9f\x9f\xbb00000000")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		if runtime.GOARCH == "386" {
			// TODO: fix refmt to properly handle 64-bit ints on 32-bit
			qt.Assert(t, err.Error(), qt.Equals, "cbor: positive integer is out of length")
		} else {
			qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
		}
	})
	t.Run("undef", func(t *testing.T) {
		// This fixture tests that we tolerate cbor's "undefined" token (even though it's noncanonical and you shouldn't use it),
		// and that it becomes a null in the data model level.
		buf := strings.NewReader("\xf7")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_Null)
	})
	t.Run("extra bytes", func(t *testing.T) {
		buf := strings.NewReader("\xa0\x00")
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, buf)
		qt.Assert(t, err, qt.Equals, ErrTrailingBytes)
	})
}

func cborMapHeader(length uint32) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0xBA)
	binary.Write(&buf, binary.BigEndian, length)
	return buf.Bytes()
}

func cborArrayHeader(length uint32) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x9A)
	binary.Write(&buf, binary.BigEndian, length)
	return buf.Bytes()
}

// TestDecodeOptions_AllocationBudget verifies that the configurable allocation
// budget is respected, both with defaults and custom values.
func TestDecodeOptions_AllocationBudget(t *testing.T) {
	t.Run("default budget rejects oversized structure", func(t *testing.T) {
		// A map header claiming more entries than the default budget allows
		payload := cborMapHeader(20_000_000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})

	t.Run("custom budget accepts within limit", func(t *testing.T) {
		// Build a small valid map: {"a": 1}
		var buf bytes.Buffer
		buf.Write(cborMapHeader(1))
		buf.WriteByte(0x61) // text(1)
		buf.WriteByte('a')
		buf.WriteByte(0x01) // uint(1)

		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 100}.Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)
		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_Map)
	})

	t.Run("custom budget rejects when exhausted", func(t *testing.T) {
		// A map claiming 50 entries with a budget of only 10 should be rejected
		payload := cborMapHeader(50)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 10}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})

	t.Run("budget accounts for declared collection sizes", func(t *testing.T) {
		// A list claiming 1000 entries consumes budget even if no entries follow
		payload := cborArrayHeader(1000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 500}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})
}

// TestDecodeOptions_MaxCollectionPrealloc verifies that the preallocation cap
// is respected and that large collections still decode correctly.
func TestDecodeOptions_MaxCollectionPrealloc(t *testing.T) {
	t.Run("large map decodes correctly with default cap", func(t *testing.T) {
		const numEntries = 5_000
		var buf bytes.Buffer
		buf.Write(cborMapHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			key := fmt.Sprintf("k%05d", i)
			buf.WriteByte(0x66) // text(6)
			buf.WriteString(key)
			if i < 24 {
				buf.WriteByte(byte(i))
			} else if i < 256 {
				buf.WriteByte(0x18)
				buf.WriteByte(byte(i))
			} else {
				buf.WriteByte(0x19)
				binary.Write(&buf, binary.BigEndian, uint16(i))
			}
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_Map)
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))

		v, err := node.LookupByString("k00000")
		qt.Assert(t, err, qt.IsNil)
		vi, err := v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(0))

		v, err = node.LookupByString("k04999")
		qt.Assert(t, err, qt.IsNil)
		vi, err = v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(4999))
	})

	t.Run("large list decodes correctly with default cap", func(t *testing.T) {
		const numEntries = 5_000
		var buf bytes.Buffer
		buf.Write(cborArrayHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			if i < 24 {
				buf.WriteByte(byte(i))
			} else if i < 256 {
				buf.WriteByte(0x18)
				buf.WriteByte(byte(i))
			} else {
				buf.WriteByte(0x19)
				binary.Write(&buf, binary.BigEndian, uint16(i))
			}
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))

		v, err := node.LookupByIndex(0)
		qt.Assert(t, err, qt.IsNil)
		vi, err := v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(0))

		v, err = node.LookupByIndex(numEntries - 1)
		qt.Assert(t, err, qt.IsNil)
		vi, err = v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(numEntries-1))
	})

	t.Run("custom prealloc cap with valid data", func(t *testing.T) {
		// 100-entry list with a prealloc cap of 10 should still decode fine
		const numEntries = 100
		var buf bytes.Buffer
		buf.Write(cborArrayHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			buf.WriteByte(byte(i % 24))
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, MaxCollectionPrealloc: 10}.Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))
	})
}

// TestDecodeOptions_MaxDepth verifies that the configurable nesting-depth
// limit is respected, both with defaults and custom values.
func TestDecodeOptions_MaxDepth(t *testing.T) {
	// buildNestedArrays returns depth `0x81` bytes (array(1)) followed by a
	// single `0xF6` null, forming `depth` levels of nested single-element
	// arrays.
	buildNestedArrays := func(depth int) []byte {
		buf := make([]byte, 0, depth+1)
		for i := 0; i < depth; i++ {
			buf = append(buf, 0x81)
		}
		buf = append(buf, 0xF6)
		return buf
	}

	t.Run("default depth rejects deeply nested structure", func(t *testing.T) {
		payload := buildNestedArrays(2000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("structure at default depth decodes", func(t *testing.T) {
		payload := buildNestedArrays(1024)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_List)
	})

	t.Run("custom depth rejects when exceeded", func(t *testing.T) {
		payload := buildNestedArrays(10)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 5}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("custom depth accepts within limit", func(t *testing.T) {
		payload := buildNestedArrays(5)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 10}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, nb.Build().Kind(), qt.Equals, datamodel.Kind_List)
	})

	t.Run("nested maps also limited", func(t *testing.T) {
		// Build N nested single-entry maps each with key "x" wrapping a null.
		const depth = 2000
		buf := make([]byte, 0, 3*depth+1)
		for i := 0; i < depth; i++ {
			buf = append(buf, 0xA1) // map(1)
			buf = append(buf, 0x61) // text(1)
			buf = append(buf, 'x')
		}
		buf = append(buf, 0xF6) // null

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(buf))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("zero MaxDepth resolves to default", func(t *testing.T) {
		payload := buildNestedArrays(2000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{MaxDepth: 0}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})

	t.Run("indefinite-length collections also limited", func(t *testing.T) {
		// Stream of 0x9F (indefinite list open) markers then a null, then
		// matching 0xFF break bytes. Exercises the indefinite-length branch
		// of the decoder which has a separate code path to the definite one.
		const depth = 2000
		buf := make([]byte, 0, 2*depth+1)
		for i := 0; i < depth; i++ {
			buf = append(buf, 0x9F)
		}
		buf = append(buf, 0xF6)
		for i := 0; i < depth; i++ {
			buf = append(buf, 0xFF)
		}
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(buf))
		qt.Assert(t, err, qt.Equals, ErrDecodeDepthExceeded)
	})
}

// TestDecoderBoundaries asserts that decoder-layer and tokenizer-layer limits
// behave as expected for unusual or malformed inputs. These are sanity tests
// around boundaries that callers sometimes need to reason about.
func TestDecoderBoundaries(t *testing.T) {
	t.Run("oversized string declaration rejected by tokenizer", func(t *testing.T) {
		// Text header declaring 1 TiB; no following bytes. The underlying
		// refmt tokenizer caps string/bytes length before attempting to read.
		var buf bytes.Buffer
		buf.WriteByte(0x7B) // text(uint64 length)
		binary.Write(&buf, binary.BigEndian, uint64(1<<40))

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(buf.Bytes()))
		qt.Assert(t, err, qt.Not(qt.IsNil))
	})

	t.Run("stacked CBOR tags rejected", func(t *testing.T) {
		// CBOR permits tagging a value, but the tokenizer refuses to stack
		// multiple tags on a single item. Link handling relies on this.
		payload := []byte{
			0xD8, 42, // tag(42)
			0xD8, 42, // tag(42)
			0x42, 0x00, 0x01,
		}
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Not(qt.IsNil))
	})

	t.Run("indefinite collection exceeding budget rejected", func(t *testing.T) {
		// Indefinite-length arrays have no declared size, but per-entry
		// budget still applies as entries are read.
		var buf bytes.Buffer
		buf.WriteByte(0x9F) // indefinite array
		for i := 0; i < 3_000_000; i++ {
			buf.WriteByte(0x00)
		}
		buf.WriteByte(0xFF) // break

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(buf.Bytes()))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})
}
