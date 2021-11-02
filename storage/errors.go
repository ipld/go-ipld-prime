package storage

const (
	ErrTag_StorageEntryNotAvailable = "ipld-storage-entry-not-available" // This is probably what you're looking for if you're looking for "not exists".
	ErrTag_StorageCancelled         = "ipld-storage-cancelled"           // This error is to denote that the operation was aborted due to context cancellation.
	ErrTag_StorageDisconnected      = "ipld-storage-disconnected"        // This error can indicate the network broke, there's a permissions error, storage is full, etc -- it probably wraps some other error, and just generally means the read or write failed.
	ErrTag_StorageErrorUnknown      = "ipld-storage-error-unknown"       // The package-scope functions will filter any strange errors from a storage implementation into these.  They should not otherwise be seen (storage implementations should not return them).
)
