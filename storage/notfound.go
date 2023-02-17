package storage

import (
	"errors"

	"github.com/ipfs/go-cid"
)

// ErrNotFound is a 404, but for block storage systems. It is returned when
// a block is not found. The Key is typically the binary form of a CID
// (CID#KeyString()).
//
// ErrNotFound implements `interface{NotFound() bool}`, which makes it roughly
// compatible with the legacy github.com/ipfs/go-ipld-format#ErrNotFound.
// The IsNotFound() function here will test for this and therefore be compatible
// with this ErrNotFound, and the legacy ErrNotFound. The same is not true for
// the legacy github.com/ipfs/go-ipld-format#IsNotFound.
//
// errors.Is() should be preferred as the standard Go way to test for errors;
// however due to the move of the legacy ErrNotFound to this package, it may
// not report correctly where older block storage packages emit the legacy
// ErrNotFound. The IsNotFound() function provides a maximally compatible
// matching function that should be able to determine whether an ErrNotFound,
// either new or legacy, exists within a wrapped error chain.
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

// Is allows errors.Is to work with this error type. It is compatible with the
// legacy github.com/ipfs/go-ipld-format#ErrNotFound, and other related error
// types as it uses a feature-test on the NotFound() method.
//
// It is important to note that because errors.Is() performs a reverse match,
// whereby the Is() of the error being checked is called on the target,
// the legacy ErrNotFound#Is will perform a strict type match, which will fail
// where the original error is of the legacy type. Where compatibility is
// required across multiple block storage systems that may return legacy error
// types, use the IsNotFound() function instead.
func (e ErrNotFound) Is(err error) bool {
	if v, ok := err.(interface{ NotFound() bool }); ok && v.NotFound() {
		return v.NotFound()
	}
	return false
}

var enf = ErrNotFound{}

// IsNotFound returns true if the error is a ErrNotFound, or compatible with an
// ErrNotFound, or wraps such an error. Compatibility is determined by the
// type implementing the NotFound() method which returns true.
// It is compatible with the legacy github.com/ipfs/go-ipld-format#ErrNotFound,
// and other related error types.
//
// This is NOT the same as errors.Is(err, storage.ErrNotFound{}) which relies on
// the Is() of the original err rather than the target. IsNotFound() uses the
// Is() of storage.ErrNotFound to perform the check. The difference being that
// the err being checked doesn't need to have a feature-testing Is() method for
// this to succeed, it only needs to have a NotFound() method that returns true.
//
// Prefer this method for maximal compatibility, including wrapped errors, that
// implement the minimal interface{ NotFound() true }.
func IsNotFound(err error) bool {
	for {
		if enf.Is(err) {
			return true
		}
		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}
