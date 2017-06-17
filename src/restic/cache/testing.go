package cache

import (
	"restic"
	"restic/crypto"
	"restic/test"
	"testing"
)

// TestNewCache returns a cache in a temporary directory which is removed when
// cleanup is called.
func TestNewCache(t testing.TB) (*Cache, func()) {
	dir, cleanup := test.TempDir(t)
	cache, err := New(restic.NewRandomID().String(), dir, crypto.NewRandomKey())
	if err != nil {
		t.Fatal(err)
	}
	return cache, cleanup
}
