package main

import "codeberg.org/rpg/dungeon/dungeon"

type Human struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewHuman creates and adds a new human to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewHuman(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Human {

	// Predefined patterns for different human shapes
	var humanPatterns = dungeon.ObjectPatterns{
		Simple: 'H',

		Pattern2H: [1][2]rune{{'H', 'H'}},
		Pattern2V: [2][1]rune{
			{'H'},
			{'│'},
		},

		Pattern3H: [1][3]rune{{'─', 'H', '─'}},
		Pattern3V: [3][1]rune{
			{'H'},
			{'┼'},
			{'│'},
		},

		Pattern2x2: [2][2]rune{
			{'─', '─'},
			{'│', '│'},
		},
		Pattern3x3: [3][3]rune{
			{' ', 'H', ' '},
			{'─', '┼', '─'},
			{'│', '│', '│'},
		},
	}

	h := &Human{
		Id:       itemArgs.Id,
		Type:     dungeon.Human,
		dng:      dng,
		patterns: humanPatterns,
		Health:   15,
		Strength: 8,
		RoomId:   itemArgs.RoomNumber,
	}

	h.dng.AddObjectWithShape(itemArgs, h.Type, dungeon.ColorBrightBlue, h.patterns)
	return h
}

// MoveToRoom moves a human to a specific room
func (h *Human) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	h.RoomId = args.RoomNumber
	return h.dng.MoveObjectToRoom(args, h.Type, h.patterns)
}

// Move moves a human within the dungeon
func (h *Human) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := h.dng.MoveObject(args, h.Type, h.patterns)
	h.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes a human from a specific room
func (h *Human) Remove() bool {
	return h.dng.RemoveObjectFromRoom(h.RoomId, h.Id, h.Type)
}
