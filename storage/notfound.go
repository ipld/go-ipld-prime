package storage

import "github.com/ipfs/go-cid"

// compatible with the go-ipld-format ErrNotFound, match against
// interface{NotFound() bool}
// this could go into go-ipld-prime, but for now we'll just exercise the
// feature-test pattern

// ErrNotFound is a 404, but for block storage systems. It is returned when
// a block is not found. The Key is typically the binary form of a CID
// (CID#KeyString()).
//
// ErrNotFound implements `interface{NotFound() bool}`, which makes it roughly
// compatible with the legacy github.com/ipfs/go-ipld-format#ErrNotFound.
// The IsNotFound() function here will test for this and therefore be compatible
// with this ErrNotFound, and the legacy ErrNotFound. The same is not true for
// the legacy github.com/ipfs/go-ipld-format#IsNotFound.
type ErrNotFound struct {
	Key string
}

// NewErrNotFound is a convenience factory that creates a new ErrNotFound error
// from a CID.
func NewErrNotFound(c cid.Cid) ErrNotFound {
	return ErrNotFound{Key: c.KeyString()}
}

func (e ErrNotFound) Error() string {
	if c, err := cid.Cast([]byte(e.Key)); err == nil && c != cid.Undef {
		return "ipld: could not find " + c.String()
	}
	return "ipld: could not find " + e.Key
}

// NotFound always returns true, and is used to feature-test for ErrNotFound
// errors.
func (e ErrNotFound) NotFound() bool {
	return true
}

// Is allows errors.Is to work with this error type.
func (e ErrNotFound) Is(err error) bool {
	switch err.(type) {
	case ErrNotFound:
		return true
	default:
		return false
	}
}

// IsNotFound returns true if the error is a ErrNotFound. As it uses a
// feature-test, it is also compatible with the legacy
// github.com/ipfs/go-ipld-format#ErrNotFound.
func IsNotFound(err error) bool {
	if nf, ok := err.(interface{ NotFound() bool }); ok {
		return nf.NotFound()
	}
	return false
}
