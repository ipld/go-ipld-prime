// Storage contains some simple implementations for the
// ipld.BlockReadOpener and ipld.BlockWriteOpener interfaces,
// which are typically used by composition in a LinkSystem.
//
// These are provided as simple "batteries included" storage systems.
// They are aimed at being quickly usable to build simple demonstrations.
// For heavy usage (large datasets, with caching, etc) you'll probably
// want to start looking for other libraries which go deeper on this subject.
package storage
