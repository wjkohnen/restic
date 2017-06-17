package cache

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"restic"
	"restic/test"
	"testing"
	"time"
)

func generateRandomIndexes(t testing.TB, c *Cache) restic.IDSet {
	ids := restic.NewIDSet()
	for i := 0; i < rand.Intn(15)+10; i++ {
		buf := test.Random(rand.Int(), 1<<19)
		id := restic.Hash(buf)

		if c.HasIndex(id) {
			t.Errorf("index %v present before save", id)
		}

		err := c.SaveIndex(id, bytes.NewReader(buf))
		if err != nil {
			t.Fatal(err)
		}
		ids.Insert(id)
	}
	return ids
}

// randomID returns a random ID from s.
func randomID(s restic.IDSet) restic.ID {
	for id := range s {
		return id
	}
	panic("set is empty")
}

func loadIndex(t testing.TB, c *Cache, id restic.ID) []byte {
	rd, err := c.LoadIndex(id)
	if err != nil {
		t.Fatal(err)
	}

	if rd == nil {
		t.Fatalf("LoadIndex() returned nil reader")
	}

	buf, err := ioutil.ReadAll(rd)
	if err != nil {
		t.Fatal(err)
	}

	if err = rd.Close(); err != nil {
		t.Fatal(err)
	}

	return buf
}

func listIndexes(t testing.TB, c *Cache) restic.IDSet {
	list, err := c.ListIndexes()
	if err != nil {
		t.Errorf("listing failed: %v", err)
	}

	return list
}

func clearIndexes(t testing.TB, c *Cache, valid restic.IDSet) {
	if err := c.ClearIndexes(valid); err != nil {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	seed := time.Now().Unix()
	t.Logf("seed is %v", seed)
	rand.Seed(seed)

	c, cleanup := TestNewCache(t)
	defer cleanup()

	ids := generateRandomIndexes(t, c)
	id := randomID(ids)

	id2 := restic.Hash(loadIndex(t, c, id))

	if !id.Equal(id2) {
		t.Errorf("wrong data returned, want %v, got %v", id.Str(), id2.Str())
	}

	if !c.HasIndex(id) {
		t.Errorf("cache thinks index %v isn't present", id.Str())
	}

	list := listIndexes(t, c)
	if !ids.Equals(list) {
		t.Errorf("wrong list of index IDs returned, want:\n  %v\ngot:\n  %v", ids, list)
	}

	clearIndexes(t, c, restic.NewIDSet(id))
	list2 := listIndexes(t, c)
	ids.Delete(id)
	want := restic.NewIDSet(id)
	if !list2.Equals(want) {
		t.Errorf("ClearIndexes removed indexes, want:\n  %v\ngot:\n  %v", list2, want)
	}

	clearIndexes(t, c, restic.NewIDSet())
	want = restic.NewIDSet()
	list3 := listIndexes(t, c)
	if !list3.Equals(want) {
		t.Errorf("ClearIndexes returned a wrong list, want:\n  %v\ngot:\n  %v", want, list3)
	}
}
