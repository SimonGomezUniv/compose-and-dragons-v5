package main

import (
	"context"

	"github.com/snipwise/nova/nova-sdk/agents"
	"github.com/snipwise/nova/nova-sdk/agents/structured"
	"github.com/snipwise/nova/nova-sdk/messages"
	"github.com/snipwise/nova/nova-sdk/messages/roles"
	"github.com/snipwise/nova/nova-sdk/models"
	"github.com/snipwise/nova/nova-sdk/ui/display"
)

/*
This program is used to generate a room in a dungeon, with its name, description,
and the items it contains (treasures, potions, ...).
A room can be empty, contain multiple items, or just one...
*/

type ItemType string

const (
	ItemTypeTreasure ItemType = "treasure"
	ItemTypePotion   ItemType = "potion"
	ItemTypeWeapon   ItemType = "weapon"
)

type Item struct {
	Type     ItemType `json:"type"`
	Quantity int      `json:"quantity"`
}

type Room struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Items            []Item `json:"items"`
}

func main() {

	nemoModelConfig := models.Config{
		Name:        "huggingface.co/quantfactory/nemotron-mini-4b-instruct-gguf:q4_k_m",
		Temperature: models.Float64(1.0),
		TopP:        models.Float64(0.9),
		TopK:        models.Int(40),
	}
	janModelConfig := models.Config{
		Name:        "hf.co/menlo/jan-nano-gguf:q4_k_m",
		Temperature: models.Float64(1.0),
		TopP:        models.Float64(0.9),
		TopK:        models.Int(40),
	}
	lucyModelConfig := models.Config{
		Name:        "hf.co/menlo/lucy-gguf:q4_k_m",
		Temperature: models.Float64(1.0),
		TopP:        models.Float64(0.9),
		TopK:        models.Int(40),
	}
	_ = janModelConfig
	_ = nemoModelConfig
	_ = lucyModelConfig

	ctx := context.Background()
	agent, err := structured.NewAgent[Room](
		ctx,
		agents.Config{
			EngineURL: "http://localhost:12434/engines/llama.cpp/v1",
			SystemInstructions: `
			You are an expert dungeon master creating rooms for a text-based adventure game.
			Generate a room with a name, description. a short description, 
			and a list of items (type and quantity) it contains.
			A room can be empty, contain multiple items, or just one...
			The items can be of type: treasure, potion, weapon.
			Each item should have a type and a quantity.
			The quantity should vary between 0 and 10.
			Make sure the room is unique and interesting.
			Provide the response in valid JSON format.
			`,
			KeepConversationHistory: true,
		},
		janModelConfig,
	)
	if err != nil {
		panic(err)
	}

	room, finishReason, err := agent.GenerateStructuredData([]messages.Message{
		{Role: roles.User, Content: "Generate a new room in the dungeon."},
	})
	_ = finishReason

	if err != nil {
		panic(err)
	}

	display.Title("Room Generated")
	display.KeyValue("Name", room.Name)
	display.KeyValue("Description", room.Description)
	display.KeyValue("Short Description", room.ShortDescription)
	display.Separator()
	if len(room.Items) == 0 {
		display.Info("No items in this room.")
	} else {
		for _, item := range room.Items {
			display.Println(string(item.Type) + ": " + string(rune(item.Quantity+'0')))
		}
	}

	

}
