package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	config := &TableConfig{
		MasterTables: []MasterTable{
			{
				SqlName: "players",
				Statdir: "/players/{plydbid}",
			},
		},
		SlaveTables: []SlaveTable{
			{
				SqlName: "rolcampwall",
				Statdir: "/players/{plydbid}/campwall/{wallpos}",
			},
		},
	}

	cache := NewCache(config)

	// Test Set and Get
	key := "/players/123"
	data := []byte("test player data")

	// Test Set
	err := cache.Set(key, data)
	if err != nil {
		t.Errorf("unexpected error setting data: %v", err)
	}

	// Test Get
	result, found, err := cache.Get(key)
	if err != nil {
		t.Errorf("unexpected error getting data: %v", err)
	}
	if !found {
		t.Errorf("data should be found")
	}
	if string(result) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(result))
	}

	// Test Get non-existent key
	_, found, err = cache.Get("/players/456")
	if err != nil {
		t.Errorf("unexpected error getting non-existent data: %v", err)
	}
	if found {
		t.Errorf("data should not be found")
	}

	// Test Delete
	err = cache.Delete(key)
	if err != nil {
		t.Errorf("unexpected error deleting data: %v", err)
	}

	// Verify deletion
	_, found, err = cache.Get(key)
	if err != nil {
		t.Errorf("unexpected error getting deleted data: %v", err)
	}
	if found {
		t.Errorf("data should not be found after deletion")
	}

	// Test invalid key
	err = cache.Set("/invalid/key", data)
	if err == nil {
		t.Errorf("expected error for invalid key")
	}

	// Test Size and Keys
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("cache should be empty after clear")
	}

	cache.Set("/players/123", data)
	cache.Set("/players/456/campwall/1", []byte("wall data"))

	if cache.Size() != 2 {
		t.Errorf("expected cache size 2, got %d", cache.Size())
	}

	keys := cache.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}