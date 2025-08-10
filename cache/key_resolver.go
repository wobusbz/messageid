package cache

import (
	"fmt"
	"regexp"
	"strings"
)

// KeyResolver handles the resolution of Statdir patterns to actual cache keys
type KeyResolver struct {
	config *TableConfig
}

// NewKeyResolver creates a new key resolver with the given configuration
func NewKeyResolver(config *TableConfig) *KeyResolver {
	return &KeyResolver{
		config: config,
	}
}

// ResolveKey resolves a Statdir pattern with given parameters to a cache key
// For example: "/players/{plydbid}" with params{"plydbid": "123"} -> "/players/123"
func (kr *KeyResolver) ResolveKey(statdir string, params map[string]string) (string, error) {
	if statdir == "" {
		return "", fmt.Errorf("statdir cannot be empty")
	}

	result := statdir
	
	// Find all {param} patterns in the statdir
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(statdir, -1)
	
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		
		placeholder := match[0] // e.g., "{plydbid}"
		paramName := match[1]   // e.g., "plydbid"
		
		value, exists := params[paramName]
		if !exists {
			return "", fmt.Errorf("missing parameter %s for statdir %s", paramName, statdir)
		}
		
		result = strings.Replace(result, placeholder, value, -1)
	}
	
	return result, nil
}

// ValidateKey checks if a key matches any known Statdir patterns
func (kr *KeyResolver) ValidateKey(key string) bool {
	// Check against master tables
	for _, master := range kr.config.MasterTables {
		if kr.matchesPattern(key, master.Statdir) {
			return true
		}
	}
	
	// Check against slave tables
	for _, slave := range kr.config.SlaveTables {
		if kr.matchesPattern(key, slave.Statdir) {
			return true
		}
	}
	
	return false
}

// matchesPattern checks if a key matches a Statdir pattern
func (kr *KeyResolver) matchesPattern(key, pattern string) bool {
	// Convert pattern to regex by replacing {param} with [^/]+ (non-slash characters)
	regexPattern := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(pattern, `[^/]+`)
	regexPattern = "^" + regexPattern + "$"
	
	matched, err := regexp.MatchString(regexPattern, key)
	if err != nil {
		return false
	}
	
	return matched
}

// GetTableByKey returns the table configuration that matches the given key
func (kr *KeyResolver) GetTableByKey(key string) (interface{}, bool) {
	// Check master tables first
	for _, master := range kr.config.MasterTables {
		if kr.matchesPattern(key, master.Statdir) {
			return master, true
		}
	}
	
	// Check slave tables
	for _, slave := range kr.config.SlaveTables {
		if kr.matchesPattern(key, slave.Statdir) {
			return slave, true
		}
	}
	
	return nil, false
}