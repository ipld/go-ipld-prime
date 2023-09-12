package storage_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/storage"
)

func TestNotFound(t *testing.T) {
	nf := storage.ErrNotFound{}
	if !storage.IsNotFound(nf) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if !errors.Is(nf, storage.ErrNotFound{}) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if nf.Error() != "ipld: could not find node" {
		t.Fatal("unexpected error message")
	}

	nf = storage.NewErrNotFound(cid.MustParse("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"))
	if !storage.IsNotFound(nf) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if !errors.Is(nf, storage.ErrNotFound{}) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if nf.Error() != "ipld: could not find bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi" {
		t.Fatal("unexpected error message")
	}

	wrappedNf := fmt.Errorf("wrapped outer: %w", fmt.Errorf("wrapped inner: %w", nf))
	if !storage.IsNotFound(wrappedNf) {
		t.Fatal("expected wrapped ErrNotFound to be a NotFound error")
	}
	if !errors.Is(wrappedNf, storage.ErrNotFound{}) {
		t.Fatal("expected wrapped ErrNotFound to be a NotFound error")
	}

	fmt.Println("WeirdNotFoundErr")
	wnf := weirdNotFoundError{}
	if !storage.IsNotFound(wnf) {
		t.Fatal("expected weirdNotFoundError to be a NotFound error")
	}
	if !errors.Is(wnf, storage.ErrNotFound{}) {
		t.Fatal("expected weirdNotFoundError to be a NotFound error")
	}

	// a weirder case, this one implements `NotFound()` but it returns false; but
	// it also implements the same Is() that will claim it's a not-found, so
	// it should work one way around, but not the other, when it's being asked
	// whether an error is or not
	wnnf := weirdNotNotFoundError{}
	if storage.IsNotFound(wnnf) {
		// this shouldn't be true because we test NotFound()==true
		t.Fatal("expected weirdNotNotFoundError to NOT be a NotFound error")
	}
	if !errors.Is(wnnf, storage.ErrNotFound{}) {
		// this should be true, because weirdNotNotFoundError.Is() performs the
		// check on storage.ErrNotFound{}.NotFound() which does return true.
		t.Fatal("expected weirdNotNotFoundError to be a NotFound error")
	}
	if errors.Is(nf, weirdNotNotFoundError{}) {
		// switch them around and we get the same result as storage.IsNotFound, but
		// won't work with wrapped weirdNotNotFoundError errors.
		t.Fatal("expected weirdNotNotFoundError to NOT be a NotFound error")
	}
}

type weirdNotFoundError struct{}

func (weirdNotFoundError) NotFound() bool {
	return true
}

func (weirdNotFoundError) Is(err error) bool {
	if v, ok := err.(interface{ NotFound() bool }); ok && v.NotFound() {
		return v.NotFound()
	}
	return false
}

func (weirdNotFoundError) Error() string {
	return "weird not found error"
}

type weirdNotNotFoundError struct{}

func (weirdNotNotFoundError) NotFound() bool {
	return false
}

func (weirdNotNotFoundError) Is(err error) bool {
	if v, ok := err.(interface{ NotFound() bool }); ok && v.NotFound() {
		return v.NotFound()
	}
	return false
}

func (weirdNotNotFoundError) Error() string {
	return "weird not NOT found error"
}
