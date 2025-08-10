package cache

import (
	"testing"
)

func TestKeyResolver(t *testing.T) {
	config := &TableConfig{
		MasterTables: []MasterTable{
			{
				SqlName: "players",
				Statdir: "/players/{plydbid}",
			},
			{
				SqlName: "guild",
				Statdir: "/guilds/{guildid}",
			},
		},
		SlaveTables: []SlaveTable{
			{
				SqlName: "friends2",
				Statdir: "/friends/{plydbid}/{plydbid2}",
			},
			{
				SqlName: "rolcampwall",
				Statdir: "/players/{plydbid}/campwall/{wallpos}",
			},
		},
	}

	resolver := NewKeyResolver(config)

	// Test ResolveKey
	tests := []struct {
		name     string
		statdir  string
		params   map[string]string
		expected string
		hasError bool
	}{
		{
			name:     "simple player key",
			statdir:  "/players/{plydbid}",
			params:   map[string]string{"plydbid": "123"},
			expected: "/players/123",
			hasError: false,
		},
		{
			name:     "multiple params",
			statdir:  "/friends/{plydbid}/{plydbid2}",
			params:   map[string]string{"plydbid": "123", "plydbid2": "456"},
			expected: "/friends/123/456",
			hasError: false,
		},
		{
			name:     "missing parameter",
			statdir:  "/players/{plydbid}",
			params:   map[string]string{},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolver.ResolveKey(tt.statdir, tt.params)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result)
				}
			}
		})
	}

	// Test ValidateKey
	validKeys := []string{
		"/players/123",
		"/guilds/456",
		"/friends/123/456",
		"/players/123/campwall/1",
	}

	for _, key := range validKeys {
		if !resolver.ValidateKey(key) {
			t.Errorf("key %s should be valid", key)
		}
	}

	invalidKeys := []string{
		"/invalid/123",
		"/players",
		"/guilds/456/invalid",
	}

	for _, key := range invalidKeys {
		if resolver.ValidateKey(key) {
			t.Errorf("key %s should be invalid", key)
		}
	}
}