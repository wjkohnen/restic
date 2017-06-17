package restic

import "io"

// Cache manages a local cache.
type Cache interface {
	// IsNotExist returns true if the error was caused by a non-existing file.
	IsNotExist(err error) bool

	// LoadIndex returns a reader that yields the contents of the index file with
	// the given id. rd must be closed after use. If an error is returned, the
	// ReadCloser is nil. The index is still encrypted.
	LoadIndex(id ID) (io.ReadCloser, error)

	// SaveIndex saves an index in the cache.
	SaveIndex(id ID, rd io.Reader) error

	// ClearIndexes removes all indexs from the cache that are not contained in valid.
	ClearIndexes(valid IDSet) error

	// HasIndex returns true if the index is cached.
	HasIndex(id ID) bool
}
