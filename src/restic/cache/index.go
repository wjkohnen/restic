package cache

import (
	"io"
	"restic"
)

// LoadIndex returns a reader that yields the contents of the index file with
// the given id. rd must be closed after use. If an error is returned, the
// ReadCloser is nil.
func (c *Cache) LoadIndex(id restic.ID) (io.ReadCloser, error) {
	return nil, nil
}

// ClearIndexes removes all indexs from the cache that are not contained in valid.
func (c *Cache) ClearIndexes(valid restic.IDSet) error {
	return nil
}

// ListIndexes returns a list of all index IDs in the cache.
func (c *Cache) ListIndexes() (restic.IDSet, error) {
	return nil, nil
}
