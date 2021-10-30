/*
	This package contains several useful readymade sharding functions,
	which should plug nicely into most storage implementations.

	The API contract for a sharding function is:

		func(key string, shards *[]string)

	In other words, the return is actually by a pointer to a slice which will be mutated.
	This API allows the calling code to hand in a slice with existing capacity,
	and thus allows for sharding functions to work without allocations.

	There is not a named type for this contract, because we prefer that packages
	implementing the storage APIs should be possible to write without
	being required to import any code from the go-ipld-prime module.
	However, the function type definition above can be seen in many packages.

	Not all packages use this API convention.  The `fsstore` package does;
	some other storage implementations don't use sharding functions because they don't need them;
	most of the adapter packages which target older code do not,
	because those modules have their own sharding APIs already.
*/
package sharding

// Shard_r133 is a sharding function which will return three hunks,
// the last of which is the full original key,
// and the first two of which are three bytes long.
// The prefix hunks are taken from the end of the original key,
// after skipping one byte.
// If the key is too short, padding of the ascii "0" character is used.
//
// (This somewhat odd-sounding procedure is a useful one in practice,
// because if applying it on a base32 string that's a CID or multihash (which is the typical usage),
// it avoids the uneven distribution of the trailing characters of a base32 string,
// and also avoids the uneven distribution of the prefixes of CIDs or mulithashes.)
//
// If the shards parameter is a pointer to a slice that starts at zero length
// and a capacity of at least 3, this function will operate with no allocations.
//
// Supposing the key is a base32 string (where each byte effectively contains 2^5 bits),
// if a sufficient range of keys is present that all shards are seen,
// each group of shards will contain (2^5)^3=32768 entries.
func Shard_r133(key string, shards *[]string) {
	l := len(key)
	switch {
	case l > 6:
		*shards = append(*shards, key[l-7:l-4], key[l-4:l-1], key)
	case l > 3:
		*shards = append(*shards, "000", key[l-4:l-1], key)
	default:
		*shards = append(*shards, "000", "000", key)
	}
}

// Shard_r133 is a sharding function which will return three hunks.
// It is very similar to Shard_r133, but with shorter hunks.
// The last hunk is the full original key,
// and the first two hunks are two bytes long each.
// The prefix hunks are taken from the end of the original key,
// after skipping one byte.
// If the key is too short, padding of the ascii "0" character is used.
//
// If the shards parameter is a pointer to a slice that starts at zero length
// and a capacity of at least 3, this function will operate with no allocations.
//
// Supposing the key is a base32 string (where each byte effectively contains 2^5 bits),
// if a sufficient range of keys is present that all shards are seen,
// each group of shards will contain (2^5)^2=1024 entries.
// (This is often a useful number in practice, because if one is mapping shards
// onto filesystem directories, 1024 entries is almost certainly going to fit
// efficiently within any filesystem format you're likely to encounter;
// 1024-within-1024 also means you'll see about a billion entries before
// directories on the second layer of sharding will contain more than 1024 files.
// (If we're assuming 1MB blocks of data asthe actual contents, that would be quite
// a few terabytes of storage, so this is a very nice balanced trade for
// most practical systems.))
func Shard_r122(key string, shards *[]string) {
	l := len(key)
	switch {
	case l > 4:
		*shards = append(*shards, key[l-5:l-3], key[l-3:l-1], key)
	case l > 2:
		*shards = append(*shards, "00", key[l-3:l-1], key)
	default:
		*shards = append(*shards, "00", "00", key)
	}
}

// Shard_r12 is a sharding function which will return two hunks.
// The last hunk is the full original key,
// and the first hunk is two bytes long.
// The prefix is are taken from the end of the original key,
// after skipping one byte.
// If the key is too short, the first hunk is just the ascii characters "00" instead.
//
// If the shards parameter is a pointer to a slice that starts at zero length
// and a capacity of at least 2, this function will operate with no allocations.
//
// Shard_r122 is functionally equivalent to "flatfs/shard/v1/next-to-last/2",
// as it's known in some other code -- it may be familiar as the default
// for block storage in go-ipfs.
func Shard_r12(key string, shards *[]string) {
	l := len(key)
	switch {
	case l > 2:
		*shards = append(*shards, key[l-3:l-1], key)
	default:
		*shards = append(*shards, "00", key)
	}
}
