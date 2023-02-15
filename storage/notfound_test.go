package storage_test

import (
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/storage"
)

func TestNotFound(t *testing.T) {
	nf := storage.ErrNotFound{Key: "foo"}
	if !storage.IsNotFound(nf) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if nf.Error() != "ipld: could not find foo" {
		t.Fatal("unexpected error message")
	}

	nf = storage.NewErrNotFound(cid.MustParse("bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi"))
	if !storage.IsNotFound(nf) {
		t.Fatal("expected ErrNotFound to be a NotFound error")
	}
	if nf.Error() != "ipld: could not find bafybeigdyrzt5sfp7udm7hu76uh7y26nf3efuylqabf3oclgtqy55fbzdi" {
		t.Fatal("unexpected error message")
	}

	wnf := &weirdNotFoundError{}
	if !storage.IsNotFound(wnf) {
		t.Fatal("expected weirdNotFoundError to be a NotFound error")
	}

	// a weirder case, this one implements `NotFound()` but it returns false, so
	// this shouldn't be a NotFound error
	wnnf := &weirdNotNotFoundError{}
	if storage.IsNotFound(wnnf) {
		t.Fatal("expected weirdNotNotFoundError to NOT be a NotFound error")
	}
}

type weirdNotFoundError struct{}

func (weirdNotFoundError) NotFound() bool {
	return true
}

func (weirdNotFoundError) Error() string {
	return "weird not found error"
}

type weirdNotNotFoundError struct{}

func (weirdNotNotFoundError) NotFound() bool {
	return false
}

func (weirdNotNotFoundError) Error() string {
	return "weird not NOT found error"
}
