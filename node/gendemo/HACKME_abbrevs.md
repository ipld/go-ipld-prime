abbreviations
=============

- `n` -- **n**ode, of course -- the accessor functions on node implementations usually refer to their 'this' as 'n'.
- `w` -- **w**ork-in-progress node -- you'll see this in nearly every assembler.
- `ca` -- **c**hild **a**ssembler -- the thing embedded in key assemblers and value assemblers in recursive kinds.

inside nodes:

- `x` -- a placeholder for "the thing" for types that contain only one element of data (e.g., the string inside a codegen'd node of string kind).
- `t` -- **t**able -- the slice inside most map nodes that is used for alloc amortizations and maintaining order.
- `m` -- **m**ap -- the actual map inside most map nodes (seealso `t`, which is usually a sibling).

deprecated abbrevs:

- `ta` -- **t**yped **a**ssembler -- but this should probably subjected to `s/ta/na/g`; it's a weird distinction to make.
	- maybe worth keeping just so we can have `ra` for reprassembler.
