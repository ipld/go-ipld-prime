package repose

import (
	"bytes"
	"context"
	"fmt"
	"io"

	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/traversal"
	multihash "github.com/multiformats/go-multihash"
)

func ComposeLinkLoader(
	actualLoader ActualLoader,
	nodeBuilderChooser NodeBuilderChooser,
	multicodecTable MulticodecDecodeTable,
	// there is no multihashTable, because those are effectively hardcoded in advance.
) traversal.LinkLoader {
	return func(ctx context.Context, lnk cid.Cid, lnkCtx traversal.LinkContext) (ipld.Node, error) {
		// Pick what kind of ipld.Node implementation we want to produce.
		//  (In some cases, this may be a nearly constant choice; in case of use
		//   with schema types, this might return complex information based on
		//    the LinkContext, which can be read for type system hints..!)
		nb, err := nodeBuilderChooser(ctx, lnkCtx)
		if err != nil {
			return nil, err
		}
		// Open the byte reader.
		r, err := actualLoader(ctx, lnk, lnkCtx)
		if err != nil {
			return nil, err
		}
		// Tee into hash checking and unmarshalling.
		mcDecoder, exists := multicodecTable.Table[lnk.Prefix().Codec]
		if !exists {
			return nil, fmt.Errorf("no decoder registered for multicodec %d", lnk.Prefix().Codec)
		}
		var hasher bytes.Buffer // multihash only exports bulk use, which is... really inefficient and should be fixed.
		node, decodeErr := mcDecoder(nb, io.TeeReader(r, &hasher))
		// Error checking order here is tricky.
		//  If decoding errored out, we should still run the reader to the end, to check the hash.
		//  (We still don't implement this by running the hash to the end first, because that would increase the high-water memory requirement.)
		//   ((Which we experience right now anyway because multihash's interface is silly, but we're acting as if that's fixed or will be soon.))
		//  If the hash is rejected, we should return that error (and even if there was a decodeErr, it becomes irrelevant).
		if decodeErr != nil {
			_, err := io.Copy(&hasher, r)
			if err != nil {
				return nil, err
			}
		}
		hash, err := multihash.Sum(hasher.Bytes(), lnk.Prefix().MhType, lnk.Prefix().MhLength)
		if err != nil {
			return nil, err
		}
		if hash.B58String() != lnk.Hash().B58String() {
			return nil, fmt.Errorf("hash mismatch!")
		}
		if decodeErr != nil {
			return nil, decodeErr
		}
		return node, nil
	}
}

type NodeBuilderChooser func(context.Context, traversal.LinkContext) (ipld.NodeBuilder, error)

type ActualLoader func(context.Context, cid.Cid, traversal.LinkContext) (io.Reader, error)

func ComposeLinkBuilder(
	actualStorer ActualStorer,
	multicodecTable MulticodecEncodeTable,
	// there is no multihashTable, because those are effectively hardcoded in advance.
) traversal.LinkBuilder {
	return func(
		ctx context.Context,
		node ipld.Node, lnkCtx traversal.LinkContext,
		multicodecType uint64, multihashType uint64, multihashLength int,
	) (cid.Cid, error) {
		// Open the byte writer.
		w, commit, err := actualStorer(ctx, lnkCtx)
		if err != nil {
			return cid.Undef, err
		}
		// Marshal, teeing into the storage writer and the hasher.
		mcEncoder, exists := multicodecTable.Table[multicodecType]
		if !exists {
			return cid.Undef, fmt.Errorf("no encoder registered for multicodec %d", multicodecType)
		}
		var hasher bytes.Buffer // multihash only exports bulk use, which is... really inefficient and should be fixed.
		w = io.MultiWriter(&hasher, w)
		err = mcEncoder(node, w)
		if err != nil {
			return cid.Undef, err
		}
		hash, err := multihash.Sum(hasher.Bytes(), multihashType, multihashLength)
		// FIXME finish making a CID out of this.
		// the cid package is a maze of twisty little passages all alike and I don't honestly know what's up where why.
		_ = hash
		if err := commit(cid.Undef); err != nil {
			return cid.Undef, err
		}
		panic("TODO")
	}
}

type ActualStorer func(context.Context, traversal.LinkContext) (io.Writer, StoreCommitter, error)

type StoreCommitter func(cid.Cid) error
