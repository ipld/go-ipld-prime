HACKME
======

This package is for reusable tests and benchmarks.
These test and benchmark functions work over the Node and NodeBuilder interfaces,
so they should work to test compatibility and compare performance of various implementations of Node.

This is easier said than done.


Naming conventions
------------------

### name prefix

All reusable test functions start with the name prefix `TestSpec_`.

All reusable benchmarks start with the name prefix `BenchmarkSpec_`.

The "Test" and "Benchmark" prefixes are as per the requirements of the
golang standard `testing` package.  They take `*testing.T` and `*testing.B`
arguments respectively.  They also take at least one interface argument
which is how you give your Node implementation to the test spec.

The word "Spec" reflects on the fact that these are reusable/standardized tests.

We recommend you copy-paste these method names outright into the package of your Node implementation.
It's not necessary, but it's nice for consistency.
(In the future, there may be tooling to help make automated comparisons
of different Node implementation's relative performance; this would
necessarily rely on consistent names across packages.)

If your Node implementation package has *more* tests and benchmarks that
*are not* from this reusable set, that's great -- but don't use the "Spec"
word as a segment of their name; it'll make processing bulk output easier.

### full pattern

The full pattern is:

`BenchmarkSpec_{Application}_{FixtureCohort}/codec={codec}/n={size}`

- `{Application}` means what feature or big-picture behavior we're testing.
  Examples include "Marshal", "Unmarshal", "Walk", etc.
- `{FixtureCohort}` means... well, see the names from the 'corpus' subpackage;
  it should be literally one of those strings.
- `n={size}` will be present for variable-scale benchmarks.
  You'll have to consider the Application and FixtureCohort to understand the
  context of what part of the data is being varied in size, though.
- `codec={codec}` is an example of extra info that might exist for some applications.
  For example, it might include "json" and "cbor" for "Marshal" and "Unmarshal",
  but will not be seen at all in other applications like "Walk".

The parts after the slash are those which are handled internally.
For those, you call the `BenchmarkSpec_*` function name (stopping before the first slash),
and that function will call `b.Run` to make sub-tests for all the variations.
For example, when you call `BenchmarkSpec_Walk_MapNStrMap3StrInt`, that one call
will result in a suite of tests for various sizes, each of which will be denoted
in the output by `BenchmarkSpec_Walk_MapNStrMap3StrInt/n=1`, then `.../n=2`, etc.

### variable scale benchmarks

Some corpuses have fixed sizes.  Some are variable.

With fixed-size corpuses, you'll see an integer in the "FixtureCohort" name.
For variable-size corpuses, you'll see the letter "N" in place of an integer.

See the docs in the 'corpus' subpackage for more discussion of this.
