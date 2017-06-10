package restic

import (
	"io"
	"restic"
)

// Cache manages a local cache.
type Cache interface {
	// IsNotExist returns true if the error was caused by a non-existing file.
	IsNotExist(err error) bool

	// LoadIndex returns a reader that yields the contents of the index file with
	// the given id. rd must be closed after use. If an error is returned, the
	// ReadCloser is nil.
	LoadIndex(id ID) (io.ReadCloser, error)

	// ClearIndexes removes all indexs from the cache that are not contained in valid.
	ClearIndexes(valid restic.IDSet) error
}
