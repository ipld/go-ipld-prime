package sharding

import (
	"fmt"
	"path"
)

func printShard(fn func(string, *[]string), key string) {
	v := make([]string, 0, 4)
	fn(key, &v)
	fmt.Printf("%s => %s\n", key, path.Join(v...))
}

func Example_shard_R133() {
	printShard(Shard_r133, "abcdefgh")
	printShard(Shard_r133, "abcdefg")
	printShard(Shard_r133, "abcdef")
	printShard(Shard_r133, "abcde")
	printShard(Shard_r133, "abcd")
	printShard(Shard_r133, "abc")

	// Output:
	// abcdefgh => bcd/efg/abcdefgh
	// abcdefg => abc/def/abcdefg
	// abcdef => 000/cde/abcdef
	// abcde => 000/bcd/abcde
	// abcd => 000/abc/abcd
	// abc => 000/000/abc

}

func Example_shard_r122() {
	printShard(Shard_r122, "abcdefgh")
	printShard(Shard_r122, "abcdefg")
	printShard(Shard_r122, "abcdef")
	printShard(Shard_r122, "abcde")
	printShard(Shard_r122, "abcd")
	printShard(Shard_r122, "abc")

	// Output:
	// abcdefgh => de/fg/abcdefgh
	// abcdefg => cd/ef/abcdefg
	// abcdef => bc/de/abcdef
	// abcde => ab/cd/abcde
	// abcd => 00/bc/abcd
	// abc => 00/ab/abc
}

func Example_shard_r12() {
	printShard(Shard_r12, "abcde")
	printShard(Shard_r12, "abcd")
	printShard(Shard_r12, "abc")
	printShard(Shard_r12, "ab")

	// Output:
	// abcde => cd/abcde
	// abcd => bc/abcd
	// abc => ab/abc
	// ab => 00/ab
}
