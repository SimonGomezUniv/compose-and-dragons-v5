package main

import (
	"encoding/json"
	"os"
)

// EntityMetadata stores health and strength for a single entity
type EntityMetadata struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // "player", "goblin", "skeleton", "vampire", "sphinx"
	Name     string `json:"name,omitempty"`
	Race     string `json:"race,omitempty"`
	Class    string `json:"class,omitempty"`
	Health   int    `json:"health"`
	Strength int    `json:"strength"`
}

// DungeonMetadata stores all entity metadata
type DungeonMetadata struct {
	Entities []EntityMetadata `json:"entities"`
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
