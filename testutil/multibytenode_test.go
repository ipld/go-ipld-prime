package testutil_test

import (
	"io"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/testutil"
)

func TestMultiByteNode(t *testing.T) {
	mbn := testutil.NewMultiByteNode(
		[]byte("foo"),
		[]byte("bar"),
		[]byte("baz"),
		[]byte("!"),
	)
	// Sanity check that the readseeker works.
	// (This is a test of the test, not the code under test.)

	for _, rl := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		t.Run("readseeker works with read length "+qt.Format(rl), func(t *testing.T) {
			rs, err := mbn.AsLargeBytes()
			qt.Assert(t, err, qt.IsNil)
			acc := make([]byte, 0, mbn.TotalLength())
			buf := make([]byte, rl)
			for {
				n, err := rs.Read(buf)
				if err == io.EOF {
					qt.Check(t, n, qt.Equals, 0)
					break
				}
				qt.Assert(t, err, qt.IsNil)
				acc = append(acc, buf[0:n]...)
			}
			qt.Assert(t, string(acc), qt.DeepEquals, "foobarbaz!")
		})
	}

	t.Run("readseeker can seek and read middle bytes", func(t *testing.T) {
		rs, err := mbn.AsLargeBytes()
		qt.Assert(t, err, qt.IsNil)
		_, err = rs.Seek(2, io.SeekStart)
		qt.Assert(t, err, qt.IsNil)
		buf := make([]byte, 2)
		acc := make([]byte, 0, 5)
		for len(acc) < 5 {
			n, err := rs.Read(buf)
			qt.Assert(t, err, qt.IsNil)
			acc = append(acc, buf[0:n]...)
		}
		qt.Assert(t, string(acc), qt.DeepEquals, "obarba")
	})

	t.Run("readseeker can seek and read last byte", func(t *testing.T) {
		rs, err := mbn.AsLargeBytes()
		qt.Assert(t, err, qt.IsNil)
		_, err = rs.Seek(-1, io.SeekEnd)
		qt.Assert(t, err, qt.IsNil)
		buf := make([]byte, 1)
		n, err := rs.Read(buf)
		qt.Assert(t, err, qt.IsNil)
		qt.Check(t, n, qt.Equals, 1)
		qt.Check(t, string(buf[0]), qt.Equals, "!")
	})
}
