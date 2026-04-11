package main

import "codeberg.org/rpg/dungeon/dungeon"

type Vampire struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewVampire creates and adds a new vampire to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewVampire(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Vampire {

	// Predefined patterns for different vampire shapes
	var vampirePatterns = dungeon.ObjectPatterns{
		Simple: 'V',

		Pattern2H: [1][2]rune{{'╬', '╬'}},
		Pattern2V: [2][1]rune{
			{'Ω'},
			{'╬'},
		},

		Pattern3H: [1][3]rune{{'∧', 'Ω', '∧'}},
		Pattern3V: [3][1]rune{
			{'Ω'},
			{'╬'},
			{'█'},
		},

		Pattern2x2: [2][2]rune{
			{'∧', '∧'},
			{'╬', '╬'},
		},
		Pattern3x3: [3][3]rune{
			{'∧', 'Ω', '∧'},
			{'═', '╬', '═'},
			{' ', '█', ' '},
		},
	}

	v := &Vampire{
		Id:       itemArgs.Id,
		Type:     dungeon.Vampire,
		dng:      dng,
		patterns: vampirePatterns,
		Health:   10,
		Strength: 5,
		RoomId:   itemArgs.RoomNumber,
	}

	v.dng.AddObjectWithShape(itemArgs, v.Type, dungeon.ColorBrightMagenta, v.patterns)
	return v
}

// MoveToRoom moves a vampire to a specific room
func (v *Vampire) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	v.RoomId = args.RoomNumber
	return v.dng.MoveObjectToRoom(args, v.Type, v.patterns)
}

// Move moves a vampire within the dungeon
func (v *Vampire) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := v.dng.MoveObject(args, v.Type, v.patterns)
	v.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes a vampire from a specific room
func (v *Vampire) Remove() bool {
	return v.dng.RemoveObjectFromRoom(v.RoomId, v.Id, v.Type)
}
