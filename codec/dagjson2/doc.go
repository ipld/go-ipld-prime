// Several groups of exported symbols are available at different levels of abstraction:
//
//   - You might just want the multicodec registration!  Then never deal with this package directly again.
//   - You might want to use the `Encode(Node,Writer)` and `Decode(NodeAssembler,Reader)` functions directly.
//   - You might want to use `ReusableEncoder` and `ReusableDecoder` types and their configuration options,
//     then use their Encode and Decode methods with that additional control.
//   - You might want to use the lower-level TokenReader and TokenWriter tools to process the serial data
//     as a stream, without necessary creating ipld Nodes at all.
//   - (this is a stretch) You might want to use some of the individual token processing functions,
//     perhaps as part of a totally new codec that just happens to share some behaviors with this one.
//
// The first three are exported from this package.
// The last two can be found in the "./token" subpackage.
package dagjson
