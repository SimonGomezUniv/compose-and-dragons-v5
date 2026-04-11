package main

import "codeberg.org/rpg/dungeon/dungeon"

type GoldCoins struct {
	Id       string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	patterns dungeon.ObjectPatterns
	Amount   int
	RoomId   int
}

// NewGoldCoins creates and adds a new gold coins object to the dungeon
// Default X,Y positions are centered according to the shape:
// Shape1Char: (6,2), Shape2H: (6,2), Shape2V: (6,1), Shape3H: (5,2), Shape3V: (6,1), Shape2x2: (6,1), Shape3x3: (5,1)
func NewGoldCoins(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs) *GoldCoins {

	var goldCoinsPatterns = dungeon.ObjectPatterns{
		Simple: '*',
	}

	gold := &GoldCoins{
		Id:       itemArgs.Id,
		Type:     dungeon.Gold,
		dng:      dng,
		patterns: goldCoinsPatterns,
		Amount:   100,
		RoomId:   itemArgs.RoomNumber,
	}

	gold.dng.AddObjectWithShape(itemArgs, gold.Type, dungeon.ColorBrightRed, gold.patterns)
	return gold
}

func (gold *GoldCoins) Remove() bool {
	return gold.dng.RemoveCharObjectFromRoom(gold.RoomId, gold.Id)
}
