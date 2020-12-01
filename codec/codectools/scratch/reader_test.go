package scratch

import (
	"bytes"
	"io"
	"runtime"
	"testing"

	. "github.com/warpfork/go-wish"
)

func someBytes(lo, hi byte) []byte {
	bs := make([]byte, hi-lo)
	for i := lo; i < hi; i++ {
		bs[i-lo] = byte(i)
	}
	return bs
}

var m1, m2 runtime.MemStats

func init() {
	runtime.GOMAXPROCS(1) // necessary for asking precise questions about allocation.
	runtime.GC()          // magic.  shouldn't be necessary.  is.
}
func memCheckpoint() {
	runtime.GC()
	runtime.ReadMemStats(&m1)
}
func newMallocs() int {
	runtime.GC()
	runtime.ReadMemStats(&m2)
	return int(m2.Mallocs - m1.Mallocs)
}

func TestReader(t *testing.T) {
	tests := func(t *testing.T, r *Reader, streaming bool) {
		// This is looks like one big nasty long walk, and so it is, but...
		//  the point is to test stateful interactions, so to write it this way is telling the truth.

		// Force one more GC before we begin our accounting.
		// (This doesn't really seem like it should be necessary, because indeed, memCheckpoint is also doing a forced GC;
		//  but empirically (at least as of go1.15), it is; flakey numbers for the very first test will result without this double-GC to prime things.)
		runtime.GC()

		// Fixed length short read should work.
		memCheckpoint()
		bs, err := r.Readnzc(4)
		Wish(t, newMallocs(), ShouldEqual, 0) // in either mode, should be hitting the scratch slice.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, bs, ShouldEqual, []byte{0, 1, 2, 3})
		Wish(t, r.NumRead(), ShouldEqual, int64(4))

		// Another read should work.
		memCheckpoint()
		bs, err = r.Readnzc(4)
		Wish(t, newMallocs(), ShouldEqual, 0) // in either mode, should be hitting the scratch slice.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, bs, ShouldEqual, []byte{4, 5, 6, 7})
		Wish(t, r.NumRead(), ShouldEqual, int64(8))

		// Single-byte reads should work.
		memCheckpoint()
		b, err := r.Readn1()
		Wish(t, newMallocs(), ShouldEqual, 0) // in either mode, should be hitting the scratch slice.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, b, ShouldEqual, byte(8))
		Wish(t, r.NumRead(), ShouldEqual, int64(9))

		// Unread should be valid at this point.
		memCheckpoint()
		r.Unreadn1()
		Wish(t, newMallocs(), ShouldEqual, 0) // no reason at all for this to ever alloc.
		Wish(t, r.NumRead(), ShouldEqual, int64(8))

		// Single byte reads should re-read the unread byte.
		memCheckpoint()
		b, err = r.Readn1()
		Wish(t, newMallocs(), ShouldEqual, 0) // shouldn't allocate no matter how it's implemented.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, b, ShouldEqual, byte(8))
		Wish(t, r.NumRead(), ShouldEqual, int64(9))

		// Unread again to set up for next check.
		r.Unreadn1()
		Wish(t, r.NumRead(), ShouldEqual, int64(8))

		// Any length of read should get the unread byte.
		bs = make([]byte, 6)
		memCheckpoint()
		n, err := r.Readb(bs)
		Wish(t, newMallocs(), ShouldEqual, 0) // alloc was done by caller; shouldn't be any more.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, 6)
		Wish(t, bs, ShouldEqual, []byte{8, 9, 10, 11, 12, 13})
		Wish(t, r.NumRead(), ShouldEqual, int64(14))

		// Tracking should work.
		memCheckpoint()
		r.Track()
		Wish(t, newMallocs(), ShouldEqual, 0)

		// Continue reading ahead while tracking.
		memCheckpoint()
		n, err = r.Readb(bs)
		if streaming {
			Wish(t, newMallocs(), ShouldEqual, 1) // streaming has to start a new buffer here.
		} else {
			Wish(t, newMallocs(), ShouldEqual, 0) // if in buf mode, it's all subslicing.
		}
		Wish(t, err, ShouldEqual, nil)
		Wish(t, n, ShouldEqual, 6)
		Wish(t, bs, ShouldEqual, []byte{14, 15, 16, 17, 18, 19})
		Wish(t, r.NumRead(), ShouldEqual, int64(20))

		// More reading while tracking.
		memCheckpoint()
		b, err = r.Readn1()
		Wish(t, newMallocs(), ShouldEqual, 0) // if in buf mode, subslicing; if in stream mode, happens to be amortized.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, b, ShouldEqual, byte(20))
		Wish(t, r.NumRead(), ShouldEqual, int64(21))
		b, err = r.Readn1()
		Wish(t, err, ShouldEqual, nil)
		Wish(t, b, ShouldEqual, byte(21))
		Wish(t, r.NumRead(), ShouldEqual, int64(22))

		// And an unread while tracking.
		r.Unreadn1()
		Wish(t, r.NumRead(), ShouldEqual, int64(21))

		// Poke every variation of read while tracking.
		memCheckpoint()
		bs, err = r.Readnzc(5)
		if streaming {
			Wish(t, newMallocs(), ShouldEqual, 1) // streaming happens to stride over an amortization threshhold here, causing an alloc during growing the tracking slice.
		} else {
			Wish(t, newMallocs(), ShouldEqual, 0) // if in buf mode, it's all subslicing.
		}
		Wish(t, err, ShouldEqual, nil)
		Wish(t, bs, ShouldEqual, []byte{21, 22, 23, 24, 25})
		Wish(t, r.NumRead(), ShouldEqual, int64(26))

		// StopTrack should yield all that in duplicate.
		memCheckpoint()
		bs = r.StopTrack()
		Wish(t, newMallocs(), ShouldEqual, 0) // in either mode, all the allocations (if any!) were already done.
		Wish(t, bs, ShouldEqual, someBytes(14, 26))
		Wish(t, r.NumRead(), ShouldEqual, int64(26))

		// Read outside of tracking.
		memCheckpoint()
		bs, err = r.Readnzc(4)
		Wish(t, newMallocs(), ShouldEqual, 0) // in either mode, should be hitting the scratch slice.
		Wish(t, err, ShouldEqual, nil)
		Wish(t, bs, ShouldEqual, someBytes(26, 30))
		Wish(t, r.NumRead(), ShouldEqual, int64(30))

		// Subsequent tracks should still work right.
		memCheckpoint()
		r.Track()
		r.Readnzc(4)
		r.Readnzc(2)
		bs = r.StopTrack()
		Wish(t, newMallocs(), ShouldEqual, 0) // amount read is smaller than previous tracking buffer, so shouldn't alloc even in stream mode.
		Wish(t, bs, ShouldEqual, someBytes(30, 36))
		Wish(t, r.NumRead(), ShouldEqual, int64(36))

		// Read us past the end.  Should get an error warning so.
		memCheckpoint()
		bs, err = r.Readnzc(300)
		if streaming {
			Wish(t, newMallocs(), ShouldEqual, 1) // this read is bigger than scratchByteArrayLen, so it allocs.
		} else {
			Wish(t, newMallocs(), ShouldEqual, 0) // if in buf mode, it's STILL all subslicing.
		}
		Wish(t, err, ShouldEqual, io.ErrUnexpectedEOF)
		Wish(t, bs, ShouldEqual, someBytes(36, 64))
		Wish(t, r.NumRead(), ShouldEqual, int64(64))
	}
	t.Run("StreamMode", func(t *testing.T) {
		var r Reader
		r.InitReader(bytes.NewBuffer(someBytes(0, 64)))
		tests(t, &r, true)
	})
	t.Run("SliceMode", func(t *testing.T) {
		var r Reader
		r.InitSlice(someBytes(0, 64))
		tests(t, &r, false)
	})
}

// TestTechniqueSliceExtension is just a quick demo on how slice extension works, meant for human reading.
// Long story short: yeah, you can actually use subslicing syntax to just tell a slice to get longer from its current position if it still has cap;
// and you can use this to peek back into memory that you had previously removed from view by subslicing.
func TestTechniqueSliceExtension(t *testing.T) {
	slice := make([]int, 10, 10)
	slice[5] = 5
	slice[9] = 9
	t.Logf("slice len = %d cap = %d", len(slice), cap(slice))
	subslice := slice[3:5] // if we *really* wanted to drop ability to ever see "5" and "9" again, we'd want to do `[3:5:5]` here.
	t.Logf("subslice len = %d cap = %d", len(subslice), cap(subslice))
	t.Logf("subslice: %v", subslice) // should show up just zeros -- the "5" and "9" are both in the cap area but in reach of len.
	subslice2 := subslice[0:5]
	t.Logf("subslice2 len = %d cap = %d", len(subslice2), cap(subslice2))
	t.Logf("subslice2: %v", subslice2) // should show up the "5" again
}
