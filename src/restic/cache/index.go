package cache

import (
	"io"
	"os"
	"path/filepath"
	"restic"
	"restic/errors"
)

func (c *Cache) indexFilename(id restic.ID) string {
	return filepath.Join(c.Path, "index", id.String())
}

// LoadIndex returns a reader that yields the contents of the index file with
// the given id. rd must be closed after use. If an error is returned, the
// ReadCloser is nil. The index is still encrypted.
func (c *Cache) LoadIndex(id restic.ID) (io.ReadCloser, error) {
	f, err := os.Open(c.indexFilename(id))
	if err != nil {
		return nil, errors.Wrap(err, "Open")
	}

	return f, nil
}

// SaveIndex saves an index in the cache.
func (c *Cache) SaveIndex(id restic.ID, rd io.Reader) error {
	f, err := os.Create(c.indexFilename(id))
	if err != nil {
		return errors.Wrap(err, "Create")
	}

	if _, err = io.Copy(f, rd); err != nil {
		f.Close()
		return errors.Wrap(err, "Copy")
	}

	if err = f.Close(); err != nil {
		return errors.Wrap(err, "Close")
	}

	return nil
}

// ClearIndexes removes all indexs from the cache that are not contained in valid.
func (c *Cache) ClearIndexes(valid restic.IDSet) error {
	list, err := c.ListIndexes()
	if err != nil {
		return err
	}

	for id := range list {
		if valid.Has(id) {
			continue
		}

		if err = os.Remove(c.indexFilename(id)); err != nil {
			return err
		}
	}

	return nil
}

// ListIndexes returns a list of all index IDs in the cache.
func (c *Cache) ListIndexes() (restic.IDSet, error) {
	d, err := os.Open(filepath.Join(c.Path, "index"))
	if err != nil {
		return nil, err
	}

	entries, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	if err = d.Close(); err != nil {
		return nil, err
	}

	list := restic.NewIDSet()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		id, err := restic.ParseID(e.Name())
		if err != nil {
			continue
		}

		list.Insert(id)
	}

	return list, nil
}

// HasIndex returns true if the index is cached.
func (c *Cache) HasIndex(id restic.ID) bool {
	_, err := os.Stat(c.indexFilename(id))
	if err == nil {
		return true
	}

	return false
}
