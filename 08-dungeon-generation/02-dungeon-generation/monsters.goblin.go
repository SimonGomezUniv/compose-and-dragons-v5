package main

import "codeberg.org/rpg/dungeon/dungeon"

type Goblin struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewGoblin creates and adds a new goblin to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewGoblin(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Goblin {

	// Predefined patterns for different goblin shapes
	var goblinPatterns = dungeon.ObjectPatterns{
		Simple: 'G',

		Pattern2H: [1][2]rune{{'<', '>'}},
		Pattern2V: [2][1]rune{
			{'▲'},
			{'▼'},
		},

		Pattern3H: [1][3]rune{{'<', 'o', '>'}},
		Pattern3V: [3][1]rune{
			{'o'},
			{'I'},
			{'-'},
		},

		Pattern2x2: [2][2]rune{
			{'o', 'o'},
			{'[', ']'},
		},
		Pattern3x3: [3][3]rune{
			{' ', 'o', ' '},
			{'>', 'T', '<'},
			{' ', '=', ' '},
		},
	}

	g := &Goblin{
		Id:       itemArgs.Id,
		Type:     dungeon.Goblin,
		dng:      dng,
		patterns: goblinPatterns,
		Health:   10,
		Strength: 5,
		RoomId:   itemArgs.RoomNumber,
	}

	g.dng.AddObjectWithShape(itemArgs, g.Type, dungeon.ColorBrightRed, g.patterns)
	return g
}

// MoveToRoom moves a goblin to a specific room
func (g *Goblin) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	g.RoomId = args.RoomNumber
	return g.dng.MoveObjectToRoom(args, g.Type, g.patterns)
}

// Move moves a goblin within the dungeon
func (g *Goblin) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := g.dng.MoveObject(args, g.Type, g.patterns)
	g.RoomId = newRoomId
	return newRoomId, err
}

// RemoveFromRoom removes a goblin from a specific room
func (g *Goblin) Remove() bool {
	return g.dng.RemoveObjectFromRoom(g.RoomId, g.Id, g.Type)
}

