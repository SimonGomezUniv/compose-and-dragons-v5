package main

import (
	"fmt"

	"codeberg.org/rpg/dungeon/dungeon"
)

func main() {
	// Create and generate dungeon with default colors
	dng := dungeon.NewDungeon(dungeon.DungeonColors{
		StartColor: dungeon.ColorBlue,  // Blue for R01
		EndColor:   dungeon.ColorGreen, // Green for last room
		EmptyColor: dungeon.ColorGray,  // Gray for empty spaces
	})

	// Generate with default number (1d6 + 5)
	dng.Generate()

	// Example: Add characters to rooms using the dungeon method
	// Add a player character '@' in room 1 (center position)
	dng.AddCharToRoom(1, '@', 6, 2, dungeon.ColorYellow, "")

	// Add a marker in room 5
	dng.AddCharToRoom(5, '▲', 6, 2, dungeon.ColorBrightGreen, "")

	// Add an enemy 'E' in room 2 (if it exists)
	dng.AddCharToRoom(2, 'E', 3, 1, dungeon.ColorBrightRed, "")

	// Add a treasure 'T' in the last room with background color
	if len(dng.Rooms) > 0 {
		lastRoomID := dng.Rooms[len(dng.Rooms)-1].ID
		dng.AddCharToRoom(lastRoomID, 'T', 10, 3, dungeon.ColorYellow, dungeon.BgBlue)
	}

	// Add some decorative elements in room 3
	dng.AddCharToRoom(3, '*', 0, 0, dungeon.ColorCyan, "")
	dng.AddCharToRoom(3, '*', 13, 0, dungeon.ColorCyan, "")
	dng.AddCharToRoom(3, '*', 0, 4, dungeon.ColorCyan, "")
	dng.AddCharToRoom(3, '*', 13, 4, dungeon.ColorCyan, "")

	// Display the dungeon
	fmt.Println(dng.GetDetailedGrid())

	// Show legend
	fmt.Println("\nLegend:")
	fmt.Println("  @ = Player")
	fmt.Println("  E = Enemy")
	fmt.Println("  T = Treasure")
	fmt.Println("  * = Decoration")

	// Example: Remove a character
	// fmt.Println("\nRemoving enemy from room 2...")
	// dng.RemoveCharFromRoom(2, 3, 1)
	// fmt.Println(dng.GetDetailedGrid())

	// Example: Get dungeon description
	fmt.Println("\nDungeon Description:")
	rooms := dng.GetDungeonDescription()
	for _, room := range rooms {
		fmt.Printf("Room #%d at (%d, %d)\n", room.Number, room.Coordinates.X, room.Coordinates.Y)
		for dir, connectedID := range room.Connections {
			fmt.Printf("  %s → Room #%d\n", dir, connectedID)
		}
	}
}
