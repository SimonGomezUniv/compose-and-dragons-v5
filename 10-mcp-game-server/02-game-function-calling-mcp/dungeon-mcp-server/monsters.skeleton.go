package main

import "codeberg.org/rpg/dungeon/dungeon"

type Skeleton struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	Strength int
	RoomId   int
}

// NewSkeleton creates and adds a new skeleton to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewSkeleton(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *Skeleton {

	// Predefined patterns for different skeleton shapes
	var skeletonPatterns = dungeon.ObjectPatterns{
		Simple: 'S',

		Pattern2H: [1][2]rune{{'◊', '◊'}},
		Pattern2V: [2][1]rune{
			{'◊'},
			{'║'},
		},

		Pattern3H: [1][3]rune{{'◊', '║', '◊'}},
		Pattern3V: [3][1]rune{
			{'◊'},
			{'║'},
			{'Ψ'},
		},

		Pattern2x2: [2][2]rune{
			{'◊', '◊'},
			{'║', '║'},
		},
		Pattern3x3: [3][3]rune{
			{' ', '◊', ' '},
			{'~', '║', '~'},
			{' ', 'Ψ', ' '},
		},
	}

	s := &Skeleton{
		Id:       itemArgs.Id,
		Type:     dungeon.Skeleton,
		dng:      dng,
		patterns: skeletonPatterns,
		Health:   10,
		Strength: 5,
		RoomId:   itemArgs.RoomNumber,
	}

	s.dng.AddObjectWithShape(itemArgs, s.Type, dungeon.ColorBrightWhite, s.patterns)
	return s
}

// MoveToRoom moves a skeleton to a specific room
func (s *Skeleton) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	s.RoomId = args.RoomNumber
	return s.dng.MoveObjectToRoom(args, s.Type, s.patterns)
}

// Move moves a skeleton within the dungeon
func (s *Skeleton) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := s.dng.MoveObject(args, s.Type, s.patterns)
	s.RoomId = newRoomId
	return newRoomId, err
}

// Remove removes a skeleton from a specific room
func (s *Skeleton) Remove() bool {
	return s.dng.RemoveObjectFromRoom(s.RoomId, s.Id, s.Type)
}
