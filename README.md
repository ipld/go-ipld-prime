go-ipld-prime
=============

`go-ipld-prime` is an implementation of the IPLD spec interfaces, a default "native" implementation of IPLD based on CBOR, and tooling for basic operations on IPLD objects.


API
---

- `github.com/ipld/go-ipld-prime` -- imported as just `ipld` -- contains interfaces for IPLD objects.  You can implement these interfaces too.  This package also provides a concrete implementation of merklepaths, and contains useful traversal functions.
- `github.com/ipld/go-ipld-prime/impl/cbor` -- imported as `ipldcbor` -- implements the IPLD interfaces backed with CBOR serialization.  There are lots of handy methods for you to convert other Go types to and from these nodes.
- `github.com/ipld/go-ipld-prime/cmd/ipld` -- provides a standalone command-line tool for useful operations, such as processing objects and producing CIDs (hashes) for their canonicalized IPLD forms.

### distinctions from go-ipld-interface&go-ipld-cbor

This library is a clean take on the IPLD interfaces and addresses several design decisions very differently than existing libraries:

- The Node interfaces are minimal;
- Many features known to be legacy are dropped;
- The Link implementations are purely CIDs;
- The Path implementations are provided in the same box;
- The CBOR implementation is provided in the same box;
- And several odd dependencies on blockstore and other interfaces from the rest of the IPFS ecosystem are removed.

Most of these changes are on the roadmap for the existing IPLD projects as well, but a clean break v2 simply seemed like a clearer project-management path to getting to the end.
Both the existing IPLD libraries and go-ipld-prime can co-exist on the same import path, and refer to the same kinds of serial data.
Projects wishing to migrate can do so smoothly and at their leisure.
