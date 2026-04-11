package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PlayerInventoryMetadata stores player inventory information
type PlayerInventoryMetadata struct {
	GoldCoins int      `json:"gold_coins"`
	GoldIDs   []string `json:"gold_ids"`
	Potions   int      `json:"potions"`
	PotionIDs []string `json:"potion_ids"`
}

// EntityMetadata stores health and strength for a single entity
type EntityMetadata struct {
	ID        string                   `json:"id"`
	Type      string                   `json:"type"` // "player", "goblin", "skeleton", "vampire", "sphinx"
	Name      string                   `json:"name,omitempty"`
	Race      string                   `json:"race,omitempty"`
	Class     string                   `json:"class,omitempty"`
	Health    int                      `json:"health"`
	Strength  int                      `json:"strength"`
	Inventory *PlayerInventoryMetadata `json:"inventory,omitempty"` // Only for player
}

// DungeonMetadata stores all entity metadata
type DungeonMetadata struct {
	Entities     []EntityMetadata `json:"entities"`
	RiddleSolved bool             `json:"riddle_solved"`
}

// SaveMetadata saves the dungeon metadata to a JSON file
func SaveMetadata(metadata *DungeonMetadata, filepath string) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// LoadMetadata loads the dungeon metadata from a JSON file
func LoadMetadata(filepath string) (*DungeonMetadata, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var metadata DungeonMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// GetEntityMetadata retrieves metadata for a specific entity by ID
func (m *DungeonMetadata) GetEntityMetadata(id string) *EntityMetadata {
	for i := range m.Entities {
		if m.Entities[i].ID == id {
			return &m.Entities[i]
		}
	}
	return nil
}

// AddEntityMetadata adds or updates metadata for an entity
func (m *DungeonMetadata) AddEntityMetadata(entity EntityMetadata) {
	// Check if entity already exists
	for i := range m.Entities {
		if m.Entities[i].ID == entity.ID {
			m.Entities[i] = entity
			return
		}
	}
	// Add new entity
	m.Entities = append(m.Entities, entity)
}

// GenerateTimestampedFilename generates a filename with timestamp
func GenerateTimestampedFilename(baseName string) string {
	timestamp := time.Now().Format("20060102_150405")
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	return fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
}

// FindMostRecentSave finds the most recent save file with timestamp
// Returns the timestamped filename if found, otherwise returns the default filename
func FindMostRecentSave(dataDir, baseFilename string) string {
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := strings.TrimSuffix(baseFilename, ext)
	pattern := filepath.Join(dataDir, nameWithoutExt+"_*"+ext)

	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		// No timestamped files found, return default
		return filepath.Join(dataDir, baseFilename)
	}

	// Sort matches to find the most recent (lexicographic sort works with our timestamp format)
	sort.Strings(matches)
	return matches[len(matches)-1] // Return the last one (most recent)
}
