package main

import "codeberg.org/rpg/dungeon/dungeon"

type Dwarf struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewDwarf creates and adds a new dwarf to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewDwarf(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Dwarf {

	// Predefined patterns for different dwarf shapes
	var dwarfPatterns = dungeon.ObjectPatterns{
		Simple: 'D',

		Pattern2H: [1][2]rune{{'D', 'D'}},
		Pattern2V: [2][1]rune{
			{'D'},
			{'▬'},
		},

		Pattern3H: [1][3]rune{{'≡', 'D', '≡'}},
		Pattern3V: [3][1]rune{
			{'D'},
			{'█'},
			{'▬'},
		},

		Pattern2x2: [2][2]rune{
			{'≡', '≡'},
			{'▬', '▬'},
		},
		Pattern3x3: [3][3]rune{
			{' ', 'D', ' '},
			{'≡', '█', '≡'},
			{'▬', '▬', '▬'},
		},
	}

	d := &Dwarf{
		Id:       itemArgs.Id,
		Type:     dungeon.Dwarf,
		dng:      dng,
		patterns: dwarfPatterns,
		Health:   20,
		Strength: 10,
		RoomId:   itemArgs.RoomNumber,
	}

	d.dng.AddObjectWithShape(itemArgs, d.Type, dungeon.ColorBrightYellow, d.patterns)
	return d
}

// MoveToRoom moves a dwarf to a specific room
func (d *Dwarf) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	d.RoomId = args.RoomNumber
	return d.dng.MoveObjectToRoom(args, d.Type, d.patterns)
}

// Move moves a dwarf within the dungeon
func (d *Dwarf) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := d.dng.MoveObject(args, d.Type, d.patterns)
	d.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes a dwarf from a specific room
func (d *Dwarf) Remove() bool {
	return d.dng.RemoveObjectFromRoom(d.RoomId, d.Id, d.Type)
}
