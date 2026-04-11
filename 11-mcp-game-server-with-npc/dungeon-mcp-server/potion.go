package main

import "codeberg.org/rpg/dungeon/dungeon"

type MagicPotion struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Health   int
	RoomId   int
}

// NewMagicPotion creates and adds a new magic potion object to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewMagicPotion(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *MagicPotion {

	var magicPotionPatterns = dungeon.ObjectPatterns{
		Simple: 'Y',
	}

	potion := &MagicPotion{
		Id:       itemArgs.Id,
		Type:     dungeon.Potion,
		dng:      dng,
		patterns: magicPotionPatterns,
		Health:   100,
		RoomId:   itemArgs.RoomNumber,
	}

	potion.dng.AddObjectWithShape(itemArgs, potion.Type, dungeon.ColorBrightRed, potion.patterns)
	return potion
}

func (potion *MagicPotion) Remove() bool {
	return potion.dng.RemoveCharObjectFromRoom(potion.RoomId, potion.Id)
}
