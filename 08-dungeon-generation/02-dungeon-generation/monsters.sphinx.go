package main

import "codeberg.org/rpg/dungeon/dungeon"

type Sphinx struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewSphinx creates and adds a new sphinx to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewSphinx(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Sphinx {

	// Predefined patterns for different sphinx shapes
	var sphinxPatterns = dungeon.ObjectPatterns{
		Simple: 'X',

		Pattern2H: [1][2]rune{{'Ø', 'Ø'}},
		Pattern2V: [2][1]rune{
			{'Ø'},
			{'≈'},
		},

		Pattern3H: [1][3]rune{{'«', 'Ø', '»'}},
		Pattern3V: [3][1]rune{
			{'Ø'},
			{'▓'},
			{'≈'},
		},

		Pattern2x2: [2][2]rune{
			{'«', '»'},
			{'≈', '≈'},
		},
		Pattern3x3: [3][3]rune{
			{' ', 'Ø', ' '},
			{'«', '▓', '»'},
			{'≈', '≈', '≈'},
		},
	}

	x := &Sphinx{
		Id:       itemArgs.Id,
		Type:     dungeon.Sphinx,
		dng:      dng,
		patterns: sphinxPatterns,
		Health:   10,
		Strength: 5,
		RoomId:   itemArgs.RoomNumber,
	}

	x.dng.AddObjectWithShape(itemArgs, x.Type, dungeon.ColorBrightCyan, x.patterns)
	return x
}

// MoveToRoom moves a sphinx to a specific room
func (x *Sphinx) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	x.RoomId = args.RoomNumber
	return x.dng.MoveObjectToRoom(args, x.Type, x.patterns)
}

// Move moves a sphinx within the dungeon
func (x *Sphinx) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := x.dng.MoveObject(args, x.Type, x.patterns)
	x.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes a sphinx from a specific room
func (x *Sphinx) Remove() bool {
	return x.dng.RemoveObjectFromRoom(x.RoomId, x.Id, x.Type)
}
