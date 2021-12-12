package tests

import (
	"strconv"
)

// Gen is a func which should generate key-value pairs.
// It's used to configure benchmarks.
//
// How exactly to use this is up to you, but
// a good gen function should probably return a wide variety of keys,
// and some known distribution of key and content sizes.
// If it returns the same key frequently, it should be documented,
// because key collision rates will affect benchmark results.
type Gen func() (key string, content []byte)

// NewCounterGen returns a Gen func which yields a unique value on each subsequent call,
// which is simply the base-10 string representation of an incrementing integer.
// The content and the key are the same.
func NewCounterGen(start int64) Gen {
	return func() (key string, content []byte) {
		k := strconv.FormatInt(start, 10)
		start++
		return k, []byte(k)
	}
}
