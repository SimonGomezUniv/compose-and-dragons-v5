package main

import (
	"strings"

	"github.com/snipwise/nova/nova-sdk/toolbox/files"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

func loadNPCSheetFromFile(sheetFilePath string) (string, *NPCCharacter, error) {
	// Read the character sheet to use as resource
	characterSheetContent, err := files.ReadTextFile(sheetFilePath)
	if err != nil {
		display.Errorf("❌ Error reading character sheet: %v", err)
		return "", nil, err
	}
	// load the npc json data
	npcJSONFilePath := strings.TrimSuffix(sheetFilePath, ".md") + ".json"
	npc, errJson := files.LoadFromJsonFile[NPCCharacter](npcJSONFilePath)
	if errJson != nil {
		display.Errorf("❌ Error reading NPC JSON data: %v", errJson)
		return "", nil, errJson
	}

	return characterSheetContent, &npc, nil
}

func saveNPCSheetToFile(sheetFilePath, content string, npc *NPCCharacter) error {
	// === SAVE CHARACTER SHEET TO FILE ===
	// Write to file
	err := files.WriteTextFile(sheetFilePath, content)

	if err != nil {
		display.Errorf("❌ Error saving character sheet: %v", err)
		return err
	}

	// transform to json string the npc object
	// Save NPC JSON data to file
	jsonFilePath := strings.TrimSuffix(sheetFilePath, ".md") + ".json"

	errJson := files.SaveToJsonFile(jsonFilePath, *npc)
	if errJson != nil {
		display.Errorf("❌ Error saving NPC JSON data: %v", errJson)
		return errJson
	}

	return nil
}
