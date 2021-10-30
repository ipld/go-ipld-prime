package fsstore

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Store is implements storage.ReadableStorage and storage.WritableStorage,
// as well as quite a few of the other extended storage feature interfaces,
// backing it with simple filesystem operations.
//
// This implementation uses golang's usual `os` package for IO,
// so it should be highly portable.
//
// Both the sharding and escaping functions are configurable,
// but a typical recommended setup is to use base32 encoding,
// and a sharding function that returns two shards of two characters each.
// The escaping and sharding functions should be chosen with regard to each other --
// the sharding function is applied to the escaped form.
type Store struct {
	basepath     string
	escapingFunc func(string) string
	shardingFunc func(key string, shards *[]string)
}

func (store *Store) InitDefaults(basepath string) error {
	return store.Init(
		basepath,
		func(raw string) string {
			return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(raw))
		},
		func(key string, shards *[]string) {
			l := len(key)
			switch {
			case l > 6:
				*shards = append(*shards, key[l-7:l-4], key[l-4:l-1], key)
			case l > 3:
				*shards = append(*shards, "000", key[l-4:l-1], key)
			default:
				*shards = append(*shards, "000", "000", key)
			}
		},
	)
}

func (store *Store) Init(
	basepath string,
	escapingFunc func(string) string,
	shardingFunc func(key string, shards *[]string),
) error {
	// Simple args and state check.
	if basepath == "" {
		return fmt.Errorf("fsstore: invalid setup args: need a path")
	}
	if store.basepath != "" {
		return fmt.Errorf("fsstore: cannot init: is already initialized")
	}
	store.basepath = basepath
	store.escapingFunc = escapingFunc
	store.shardingFunc = shardingFunc

	// Make sure basepath is a dir, and make sure the staging and content dirs exist.
	if err := CheckAndMakeBasepath(basepath); err != nil {
		return err
	}

	// That's it for setup on this one.
	return nil
}

// pathForKey applies sharding funcs as well as adds the basepath prefix,
// returning a string ready to use as a filesystem path.
func (store *Store) pathForKey(key string) string {
	shards := make([]string, 1, 4) // future work: would be nice if we could reuse this rather than fresh allocating.
	shards[0] = store.basepath     // not part of the path shard, but will be a param to Join, so, practical to put here.
	//shards[1] = storageDir       // not part of the path shard, but will be a param to Join, so, practical to put here.
	store.shardingFunc(key, &shards)
	return filepath.Join(shards...)
}

// Has implements go-ipld-prime/storage.Storage.Has.
func (store *Store) Has(ctx context.Context, key string) (bool, error) {
	_, err := os.Stat(store.pathForKey(key))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Get implements go-ipld-prime/storage.ReadableStorage.Get.
func (store *Store) Get(ctx context.Context, key string) ([]byte, error) {
	f, err := store.GetStream(ctx, key)
	if err != nil {
		return nil, err
	}
	defer f.(io.Closer).Close()
	return ioutil.ReadAll(f)
}

// Put implements go-ipld-prime/storage.WritableStorage.Put.
func (store *Store) Put(ctx context.Context, key string, content []byte) error {
	// We can't improve much on what we get by wrapping the stream interface;
	//  we always end up using a streaming action on the very bottom because that's how file writing works
	//   (especially since we care about controlling the write flow enough to be able to do the atomic move at the end).
	wr, wrCommitter, err := store.PutStream(ctx)
	if err != nil {
		return err
	}
	// Write, all at once.
	// Note we can ignore the size return, because the contract of io.Writer states "Write must return a non-nil error if it returns n < len(p)".
	_, err = wr.Write(content)
	if err != nil {
		wrCommitter("")
		return err
	}
	// Commit.
	return wrCommitter(key)
}

// GetStream implements go-ipld-prime/storage.StreamingReadableStorage.GetStream.
//
// Note that the returned reader will also be an io.Closer;
// if the caller does not check for that, and call Close appropriately, as StreamingReadableStorage documents is appropriate,
// then there may be resource leaks.
func (store *Store) GetStream(ctx context.Context, key string) (io.Reader, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Figure out where we expect it to be.
	destpath := store.pathForKey(key)

	// Open and return.
	// TODO: we should normalize things like "not exists" errors before hurling them up the stack.
	return os.OpenFile(destpath, os.O_RDONLY, 0)
}

// PutStream implements go-ipld-prime/storage.StreamingWritableStorage.PutStream.
func (store *Store) PutStream(ctx context.Context) (io.Writer, func(string) error, error) {
	for {
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		// Open a new file in the staging area, with a random name.
		var bs [8]byte
		rand.Read(bs[:])
		stagepath := filepath.Join(store.basepath, stagingDir, hex.EncodeToString(bs[:]))
		f, err := os.OpenFile(stagepath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("fsstore.BeginWrite: could not create a staging file: %w", err)
		}
		// Okay, got a handle.  Return it... and its commit closure.
		return f, func(key string) error {
			// Close the staging file.
			if err := f.Close(); err != nil {
				return err
			}
			if key == "" {
				return os.Remove(stagepath)
			}
			// n.b. there is a lack of fsync here.  I am going to choose to believe that a sane filesystem will not let me do a 'move' without flushing somewhere in between.

			// Figure out where we want it to go.
			destpath := store.pathForKey(key)

			// Get it there.
			return move(stagepath, destpath)
		}, nil
	}
}

const stagingDir = ".temp" // same as flatfs uses.

func CheckAndMakeBasepath(basepath string) error {
	// Is this basepath a dir?
	// (This is TOCTOU, obviously, but also it's nice to sanity check early and return error quickly because it's probably a setup error.)
	if fi, err := os.Stat(basepath); err != nil {
		return fmt.Errorf("fsstore: cannot init: basepath must be a directory: %w", err)
	} else {
		if !fi.IsDir() {
			return fmt.Errorf("fsstore: cannot init: basepath must be a directory")
		}
	}

	// Make sure the staging dir exists.
	err := os.Mkdir(filepath.Join(basepath, stagingDir), 0777)
	switch {
	case err == nil:
		// excellent.
	case os.IsExist(err):
		// sanity check it's a directory already.
		fi, err := os.Stat(filepath.Join(basepath, stagingDir))
		if err != nil {
			return fmt.Errorf("fsstore: failed to make staging dir: %w", err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("fsstore: staging dir path contains not a dir!")
		}
	default:
		return fmt.Errorf("fsstore: failed to make staging dir: %w", err)
	}

	return nil
}

// move file at stagepath to destpath.
// First, attempt to directly rename to the destination;
// if we get a ENOENT error code, that means the parent didn't exist, and we make that and then retry.
// If making the parent failed: recurse, and use similar logic.
//
// This optimistic approach should have fewer syscall RTTs when most of the parents exist
// than would be taken if we checked that each parent segment exists.
//
// (An alternative approach would be to blindly mkdir the parent segments every time,
// rather than do this backwards stepping.  Have not benchmarked these against each other.)
func move(stagepath, destpath string) error {
	err := os.Rename(stagepath, destpath)
	if os.IsNotExist(err) {
		// This probably means parent of destpath doesn't exist yet, so we'll make it.
		//  It's technically a race condition to assume that this is because destpath has no parents vs that stagepath hasn't been removed out from underneath us, but, alas; kernel ABIs.
		//   If we did this will all fds, it could be somewhat better.
		//    (This is certainly possible, at least in linux; but we'd have to import the syscall package and do it ourselves, which is not a rubicon we're willing to cross in this package.)
		//   In practice, this is probably not going to kerfuffle things.
		if err := haveDir(filepath.Dir(destpath)); err != nil {
			return err
		}
		// Now try again.
		//  (And don't return quite yet; there's one more check to do, because someone might've raced us.)
		err = os.Rename(stagepath, destpath)
	}
	if os.IsExist(err) {
		// Oh!  Some content is already there?
		//  We're a write-once (presumed-to-be-)content-addressable blob store -- that means *we keep what already exists*.
		//  FIXME: no, I wish this is how the Rename function worked, but it is not, actually.
		return os.Remove(stagepath)
	}
	return err
}

// haveDir tries to make sure a directory exists at pth.
// If this sounds a lot like os.MkdirAll: yes,
// except this function is going to assume if it exists, it's a dir,
// and that saves us some stat syscalls.
func haveDir(pth string) error {
	err := os.Mkdir(pth, 0777)
	if os.IsNotExist(err) {
		if err := haveDir(filepath.Dir(pth)); err != nil {
			return err
		}
		return os.Mkdir(pth, 0777)
	}
	return err
}
