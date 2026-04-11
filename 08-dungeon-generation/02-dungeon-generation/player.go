package main

import "codeberg.org/rpg/dungeon/dungeon"

type Player struct {
	Id       string
	Name     string
	Race     string
	Class    string
	Type     dungeon.Item
	dng      *dungeon.Dungeon
	RoomId   int
	Health   int
	Strength int
}

func NewPlayer(dng *dungeon.Dungeon, itemArgs dungeon.ItemArgs, name, race, class string) *Player {
	p := &Player{
		Id:       itemArgs.Id,
		Name:     name,
		Race:     race,
		Class:    class,
		Type:     dungeon.Player,
		dng:      dng,
		RoomId:   itemArgs.RoomNumber,
		Health:   50,
		Strength: 5,
	}
	p.dng.AddPlayer(itemArgs)

	return p
}

func (p *Player) MoveToRoom(args dungeon.MoveToRoomArgs) bool {
	p.RoomId = args.RoomNumber
	return p.dng.MovePlayerToRoom(args)
}

// Move moves a goblin within the dungeon
func (p *Player) Move(args dungeon.MoveArgs) (int, error) {
	newRoomId, err := p.dng.MovePlayer(args)
	p.RoomId = newRoomId
	return newRoomId, err
}

func (p *Player) Remove() bool {
	return p.dng.RemovePlayerFromRoom(p.RoomId, p.Id)
}
