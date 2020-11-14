package scratch

import (
	"io"
)

// Reader provides zero-copy read methods and other helpful utilities commonly desired in parsers
// around either an io.Reader or a plain byte slice.
//
// Read methods with 'n' in the name take a size parameter for how much to read.
// Read methods with a number in the name read that fixed number of bytes.
// Read methods with 'b' in the name accept a byte slice parameter which will be used for output, allowing you to control memory reuse.
// Read methods with 'z' in the name will attempt to return zero-copy access to buffers controlled by the Reader --
// be careful when using these 'z' methods; it is not recommended to expose the zero-copy slices these methods yield,
// because the reader itself may also reuse them, and so the likelihood of spooky-action-at-a-distance bugs is high.
//
// While this Reader does some buffering, it's not much (and primarily oriented around reuse of scratch buffers rather than intentionly large batch readaheads);
// it may still be advisable to use a buffered reader to avoid small reads if reading streamingly from external IO like disk or network.
type Reader struct {
	stream     io.Reader                 // source to keep pumping, or may be nil if we're wrapping a single already-in-memory slice and we know it.
	buf        []byte                    // alternative to `stream`, if we have a single already-in-memory slice and we know it.
	cursor     int                       // position of start of next read, if we're using `buf`.
	scratch    [scratchByteArrayLen]byte // temp byte array re-used internally for efficiency during read.  'readz' methods return views into this.
	numRead    int64                     // aggregate number of bytes read (since last reset of numRead, anyway).
	tracked    []byte                    // bytes that have been read while we've been in tracking state.  a subslice of `buf` where possible, but may be its own alloc if we're in streaming mode.
	unread     byte                      // a byte tracked for the potential need for unreading.  only used if using `stream`; if using `buf`, we just adjusted `cursor`.
	isTracking bool                      // whether we're in the tracking state.
	canUnread  bool                      // whether unread is currently valid.
	haveUnread bool                      // whether we need to replay an unread byte.  only checked when in `stream` mode, because `buf` mode just adjusts `cursor`.
}

// You'll find many implementation methods have a large switch around `z.stream == nil`.
// This is effectively a toggle for whether we're operating in streaming mode or on already-in-memory byte slices.
// This would've been cleaner code with an interface and two implementations -- no doubt!
// However, it ends up less inliner-friendly if an interface is involved.
//
// Stylistically: I've allowed rightward drift in 'if' cases for stream vs buf mode,
// rather than using the usual golang rule of thumb about early returns.  I find this easier to read, given the semiotics.
//
// FUTURE: it may be worth reviewing the utility of this when go1.16 is out -- some of its features for optimization
//  through interfaces when concrete types can be inferred might change the consequences of this design quite a bit.

const (
	scratchByteArrayLen = 32
)

var (
	zeroByteSlice = []byte{}[:0:0]
)

// Init makes this Reader ready to consume the given io.Reader.
// If this Reader has been used before, all state is zeroed out cleanly.
//
// As a convenience, if the io.Reader looks like it can return all the bytes at once
// (e.g., it has a `Bytes() []byte` method -- as bytes.Buffer does, for example),
// then Init will access that and use InitSlice, which should lead to better performance.
func (z *Reader) Init(r io.Reader) {
	type BytesAccessor interface {
		Bytes() []byte
	}
	if ba, ok := r.(BytesAccessor); ok {
		z.InitSlice(ba.Bytes())
	} else {
		z.InitReader(r)
	}
}

// InitSlice makes this Reader ready to consume the given byte slice.
// If this Reader has been used before, all state is zeroed out cleanly.
//
// InitSlice is functionally equivalent to wrapping the byte slice in a reader and using Init,
// but will result in a Reader that generally operates somewhat faster and is able to deliver more zero-copy behaviors.
// (When we know we're working with a byte slice that's already entirely in memory,
// we never have to worry about read alignment, etc.)
func (z *Reader) InitSlice(bs []byte) {
	*z = Reader{}
	z.buf = bs
}

// InitReader makes this Reader ready to consume the given io.Reader.
// If this Reader has been used before, all state is zeroed out cleanly.
//
// Unlike Init, this initializer will not attempt to autodetect any interface
// which may provide direct access to underlying byte slices; it will always work in stream mode.
func (z *Reader) InitReader(r io.Reader) {
	*z = Reader{} // FUTURE: this could try to recycle any capacity in z.tracked.
	z.stream = r
}

// Readnzc read up to n bytes into a byte slice which may be shared and must not be reused after any additional calls to this reader.
// Readnzc will use the implementation scratch buffer if possible, (i.e. n < scratchByteArrayLen),
// or may return a view of the []byte being decoded from if the read is larger.
// If there is less than n bytes to be read, a shorter slice will be returned, and err will be ErrUnexpectedEOF.
// Requesting a zero length read will return `zeroByteSlice`, a len-zero cap-zero slice.
// If you know your read may be longer than scratchByteArrayLen and
// you already have an existing slice of sufficient size to reuse, prefer `Readb`.
func (z *Reader) Readnzc(n int) (bs []byte, err error) {
	if n == 0 {
		return zeroByteSlice, nil
	}
	z.canUnread = false
	if z.stream == nil { // in `buf` mode, we can just return subslices.
		remaining := len(z.buf) - z.cursor
		if n > remaining { // partial read from end of buf
			n = remaining             // mostly the same, just shorter
			err = io.ErrUnexpectedEOF // and give notice of the short read
		}
		bs = z.buf[z.cursor : z.cursor+n]
		z.cursor += n
		z.numRead += int64(n)
		if z.isTracking {
			z.tracked = z.tracked[:len(z.tracked)+n] // See TestTechniqueSliceExtension if this bewilders you.
		}
		return
	} else { // in `stream` mode, we'll set up buffers, then use Readb do to most of the work.
		if n < len(z.scratch) { // read from stream and fits in scratch
			bs = z.scratch[:n]
		} else { // read from stream and needs a new allocation
			bs = make([]byte, n) // this is a sadpath; you should've used Readb.
		}
		n, err = z.readStream(bs)
		return bs[:n], err
	}
}

// Readb reads up to `len(b)` bytes into the given slice, starting at its beginning,
// overwriting all values, and disregarding any extra capacity.
// If the there is less than `len(b)` bytes to be read, a partial read will be returned:
// some of the slice will be modified, n will be less than the slice length, and err will be ErrUnexpectedEOF.
// (If you're intentionally providing a larger slice than may be necessary in order to get a batch read,
// you will want to check for and discard ErrUnexpectedEOF!)
// If no error is returned, n will always be the length of the slice.
//
// Readb will never return a zero-copy subslice of an existing buffer;
// use one of the 'Read*z*' methods for that.
func (z *Reader) Readb(bs []byte) (n int, err error) {
	if len(bs) == 0 {
		return 0, nil
	}
	z.canUnread = false
	if z.stream == nil { // in `buf` mode, we can just return subslices.
		n = len(bs)
		remaining := len(z.buf) - z.cursor
		if n > remaining { // partial read from end of buf
			n = remaining             // mostly the same, just shorter
			err = io.ErrUnexpectedEOF // and give notice of the short read
		}
		copy(bs, z.buf[z.cursor:z.cursor+n])
		z.cursor += n
		z.numRead += int64(n)
		if z.isTracking {
			z.tracked = z.tracked[:len(z.tracked)+n] // See TestTechniqueSliceExtension if this bewilders you.
		}
		return
	} else {
		return z.readStream(bs)
	}
}

func (z *Reader) readStream(bs []byte) (n int, err error) {
	// fun note: a corresponding readBuf method turned out not useful to create,
	//  because the different return conventions of the exported methods actually matter to what shortcuts we can take when wrangling raw slices
	//   (whereas the impact of those return conventions turn out not to carry as far when we already have to handle extra slices as we do in `stream` mode).

	// In `stream` mode, we first handle replaying unreads if necessary; then, use io.ReadAtLeast to load as much data as requested.
	if z.haveUnread {
		bs[0] = z.unread
		z.haveUnread = false
		n, err = io.ReadAtLeast(z.stream, bs[1:], len(bs)-1)
		n++
	} else {
		n, err = io.ReadAtLeast(z.stream, bs, len(bs))
	}
	z.numRead += int64(n)
	if z.isTracking {
		z.tracked = append(z.tracked, bs[:n]...)
	}
	return
}

// Readn reads up to n bytes into a new byte slice.
// If there is less than n bytes to be read, a shorter slice will be returned, and err will be ErrUnexpectedEOF.
// If zero-copy views into existing buffers are acceptable (e.g. you know you
// won't later mutate, reference or expose this memory again), prefer `Readnzc`.
// If you already have an existing slice of sufficient size to reuse, prefer `Readb`.
// Requesting a zero length read will return `zeroByteSlice`, a len-zero cap-zero slice.
//
// Readn will never return a zero-copy subslice of an existing buffer;
// use one of the 'Read*z*' methods for that.
// (Readn is purely a convenience method; you can always use Readb to equivalent effect.)
func (z *Reader) Readn(n int) (bs []byte, err error) {
	if n == 0 {
		return zeroByteSlice, nil
	}
	// This really is just a convenience method.  It's the same regardless of mode we're in.
	bs = make([]byte, n)
	n, err = z.Readb(bs)
	return bs[:n], err
}

// Readn1 reads a single byte.
func (z *Reader) Readn1() (byte, error) {
	// Just use Readnzc, which handles any tracking shifts, and also any unread replays, transparently.
	//  Hopefully the compiler is clever enough to make the assembly shorter than the source.
	//   REVIEW: may want to look especially at the benchmark and the assembly on this; it might be improvable by hand-rolling more of this specialization,
	//    and it's probably important to do so, considering how much of parsing for textual formats like json involves single-byte scanning.
	bs, err := z.Readnzc(1)
	if err != nil {
		return 0, err
	}
	z.canUnread = true
	z.unread = bs[0]
	return bs[0], nil
}

// Unreadn1 "unreads" a single byte which was previously read by Readn1.
// The result is that subsequent reads will include that byte,
// and applying the Track method will also cause the track result to include that byte.
//
// Unreadn1 can only be used when the previous call was Readn1, and may panic otherwise.
func (z *Reader) Unreadn1() {
	if !z.canUnread {
		panic("Unreadn1 can only be used following Readn1")
	}
	z.canUnread = false
	z.numRead--
	if z.isTracking {
		z.tracked = z.tracked[0 : len(z.tracked)-1]
	}
	if z.stream == nil {
		z.cursor--
	} else {
		z.haveUnread = true
	}
}

func (z *Reader) NumRead() int64 {
	return z.numRead
}

func (z *Reader) ResetNumRead() {
	z.numRead = 0
}

// Track causes the Reader to place a marker and accumulate all bytes into a single contiguous slice
// up until StopTrack is called; StopTrack will return a reference to this slice.
// Thus, StopTrack will yield bytes which have already also been seen via other read method calls.
//
// This can be useful when parsing logic requires scanning ahead to look for the end of an unknown-length segment of data, for example.
//
// Calling Track twice without an intervening StopTrack will result in panic.
func (z *Reader) Track() {
	if z.isTracking {
		panic("Track cannot be called again when already tracking")
	}
	z.isTracking = true
	if z.stream == nil {
		// save the start position.  we'll just extend the length of it over the cap of buf as we go forward.
		z.tracked = z.buf[z.cursor:z.cursor]
	} else {
		// nothing to do for stream mode; it'll just accumulate naturally through appends.
	}
}

// StopTrack returns the byte slice accumulated since Track was called, and drops the marker.
//
// Calling StopTrack when Track is not in effect will result in panic.
//
// The slice returned by StopTrack may be reused if Track is called again in the future;
// the caller should copy the contents to a new byte slice before the next call to Track
// they intend to either make this data available for a long time or to mutate it.
func (z *Reader) StopTrack() []byte {
	if !z.isTracking {
		panic("StopTrack cannot be called when not tracking")
	}
	z.isTracking = false
	answer := z.tracked
	z.tracked = z.tracked[0:0]
	return answer
}
