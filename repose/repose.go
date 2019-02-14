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

// An example use of LinkContext in LinkLoader logic might be inspecting the
// LinkNode, and if it's using the type system, inspecting its Type property;
// then deciding on whether or not we want to load objects of that Type.
// This might be used to do a traversal which looks at all directory objects,
// but not file contents, for example.
type LinkContext struct {
	LinkPath   traversal.Path // whoops, nice cycle
	LinkNode   ipld.Node      // has the cid again, but also might have type info // always zero for writing new nodes, for obvi reasons.
	ParentNode ipld.Node
}

type LinkLoader func(context.Context, cid.Cid, LinkContext) (ipld.Node, error)

// One presumes implementations will also *save* the content somewhere.
// The LinkContext parameter is a nod to this -- none of those parameters
// are relevant to the generation of the Cid itself, but perhaps one might
// want to use them when deciding where to store some files, etc.
type LinkBuilder func(
	ctx context.Context,
	node ipld.Node, lnkCtx LinkContext,
	multicodecType uint64, multihashType uint64, multihashLength int,
) (cid.Cid, error)

func ComposeLinkLoader(
	actualLoader ActualLoader,
	multicodecTable MulticodecDecodeTable,
	// there is no multihashTable, because those are effectively hardcoded in advance.
) LinkLoader {
	return func(ctx context.Context, lnk cid.Cid, lnkCtx LinkContext) (ipld.Node, error) {
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
		node, decodeErr := mcDecoder(io.TeeReader(r, &hasher))
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

type ActualLoader func(context.Context, cid.Cid, LinkContext) (io.Reader, error)

func ComposeLinkBuilder(
	actualStorer ActualStorer,
	multicodecTable MulticodecEncodeTable,
	// there is no multihashTable, because those are effectively hardcoded in advance.
) LinkBuilder {
	return func(
		ctx context.Context,
		node ipld.Node, lnkCtx LinkContext,
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

type ActualStorer func(context.Context, LinkContext) (io.Writer, StoreCommitter, error)

type StoreCommitter func(cid.Cid) error
