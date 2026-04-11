package main

import (
	"fmt"

	"codeberg.org/rpg/dungeon/dungeon"
)

type Inventory struct {
	GoldCoins     int
	GoldIDs       []string // IDs of collected gold
	Potions       int
	PotionIDs     []string // IDs of collected potions
}

type Player struct {
	Id        string
	Name      string
	Race      string
	Class     string
	Type      dungeon.Item
	dng       *dungeon.Dungeon
	RoomId    int
	Health    int
	Strength  int
	Inventory Inventory
}

func NewPlayer(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs, name, race, class string) *Player {
	// Ensure player has a color - use ColorBrightBlue by default
	if itemArgs.ForeColor == "" {
		itemArgs.ForeColor = dungeon.ColorBrightBlue
	}

	p := &Player{
		Id:       itemArgs.Id,
		Name:     name,
		Race:     race,
		Class:    class,
		Type:     dungeon.Player,
		dng:      dng,
		RoomId:   itemArgs.RoomNumber,
		Health:   50,
		Strength: 5,
		Inventory: Inventory{
			GoldCoins: 0,
			GoldIDs:   []string{},
			Potions:   0,
			PotionIDs: []string{},
		},
	}
	p.dng.AddPlayer(itemArgs)

	return p
}

// CollectGold adds gold to player's inventory
func (p *Player) CollectGold(amount int, id string) {
	p.Inventory.GoldCoins += amount
	p.Inventory.GoldIDs = append(p.Inventory.GoldIDs, id)
}

// CollectPotion adds a potion to player's inventory
func (p *Player) CollectPotion(id string) {
	p.Inventory.Potions++
	p.Inventory.PotionIDs = append(p.Inventory.PotionIDs, id)
}

// DrinkPotion consumes a potion and increases health by 5
func (p *Player) DrinkPotion() bool {
	if p.Inventory.Potions <= 0 {
		fmt.Println("\n❌ You don't have any potions to drink!")
		return false
	}

	// Remove one potion from inventory
	p.Inventory.Potions--
	if len(p.Inventory.PotionIDs) > 0 {
		// Remove the last potion ID
		p.Inventory.PotionIDs = p.Inventory.PotionIDs[:len(p.Inventory.PotionIDs)-1]
	}

	// Increase health by 5
	p.Health += 5
	fmt.Printf("\n🧪 You drank a magic potion! Health increased by 5.\n")
	fmt.Printf("💚 Current Health: %d ❤️\n", p.Health)

	return true
}

// ShowInventory displays the player's inventory
func (p *Player) ShowInventory() {
	fmt.Println("\n=== Your Inventory ===")
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Race: %s\n", p.Race)
	fmt.Printf("Class: %s\n", p.Class)
	fmt.Printf("Health: %d ❤️\n", p.Health)
	fmt.Printf("Strength: %d 💪\n", p.Strength)
	fmt.Printf("Gold Coins: %d\n", p.Inventory.GoldCoins)
	if len(p.Inventory.GoldIDs) > 0 {
		fmt.Println("  Gold IDs:")
		for _, id := range p.Inventory.GoldIDs {
			fmt.Printf("    - %s\n", id)
		}
	}
	fmt.Printf("Potions: %d\n", p.Inventory.Potions)
	if len(p.Inventory.PotionIDs) > 0 {
		fmt.Println("  Potion IDs:")
		for _, id := range p.Inventory.PotionIDs {
			fmt.Printf("    - %s\n", id)
		}
	}
	fmt.Println("=====================")
}

func (p *Player) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	p.RoomId = args.RoomNumber
	return p.dng.MovePlayerToRoom(args)
}

// Move moves a goblin within the dungeon
func (p *Player) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := p.dng.MovePlayer(args)
	p.RoomId = newRoomId
	return newRoomId, err
}

func (p *Player) Remove() bool {
	return p.dng.RemovePlayerFromRoom(p.RoomId, p.Id)
}
