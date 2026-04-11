package main

import "codeberg.org/rpg/dungeon/dungeon"

type Elf struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewElf creates and adds a new elf to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewElf(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Elf {

	// Predefined patterns for different elf shapes
	var elfPatterns = dungeon.ObjectPatterns{
		Simple: 'E',

		Pattern2H: [1][2]rune{{'E', 'E'}},
		Pattern2V: [2][1]rune{
			{'E'},
			{'╪'},
		},

		Pattern3H: [1][3]rune{{'~', 'E', '~'}},
		Pattern3V: [3][1]rune{
			{'E'},
			{'║'},
			{'╪'},
		},

		Pattern2x2: [2][2]rune{
			{'~', '~'},
			{'╪', '╪'},
		},
		Pattern3x3: [3][3]rune{
			{' ', 'E', ' '},
			{'~', '║', '~'},
			{'╪', '╪', '╪'},
		},
	}

	e := &Elf{
		Id:       itemArgs.Id,
		Type:     dungeon.Elf,
		dng:      dng,
		patterns: elfPatterns,
		Health:   15,
		Strength: 7,
		RoomId:   itemArgs.RoomNumber,
	}

	e.dng.AddObjectWithShape(itemArgs, e.Type, dungeon.ColorBrightGreen, e.patterns)
	return e
}

// MoveToRoom moves an elf to a specific room
func (e *Elf) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	e.RoomId = args.RoomNumber
	return e.dng.MoveObjectToRoom(args, e.Type, e.patterns)
}

// Move moves an elf within the dungeon
func (e *Elf) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := e.dng.MoveObject(args, e.Type, e.patterns)
	e.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes an elf from a specific room
func (e *Elf) Remove() bool {
	return e.dng.RemoveObjectFromRoom(e.RoomId, e.Id, e.Type)
}
