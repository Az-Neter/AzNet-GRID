package main

import "fmt"

type World struct {
	characters []*Character
	rooms      []*Room
}

func NewWorld() *World {
	return &World{}
}

func (w *World) Init() {
	w.rooms = []*Room{
		{
			Id:   "0x00000000",
			Desc: "The space stretches out in all directions from this point, navigable only by your imagination.",
			Links: []*RoomLink{
				{
					Verb:   "east",
					RoomId: "0x0000000a",
				},
			},
		},
		{
			Id:   "0x00000001",
			Desc: "You cannot leave this area yet, but you can see a path to the east back to where you originated.",
			Links: []*RoomLink{
				{
					Verb:   "west",
					RoomId: "0x0000000b",
				},
			},
		},
	}
}

func (w *World) HandleCharacterJoined(character *Character) {
	w.rooms[0].AddCharacter(character)

	character.SendMessage("\n >Welcome to: \n ______  ______ _____ ______ \n|  ____ |_____/   |   |     \\ \n|_____| |    \\_ __|__ |_____/								 ")
	character.SendMessage("")
	character.SendMessage(character.Room.Desc)
}

func (w *World) GetRoomById(id string) *Room {
	for _, r := range w.rooms {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func (w *World) HandleCharacterInput(character *Character, input string) {
	room := character.Room
	for _, link := range room.Links {
		if link.Verb == input {
			target := w.GetRoomById(link.RoomId)
			if target != nil {
				w.MoveCharacter(character, target)
				return
			}
		}
	}

	character.SendMessage(fmt.Sprintf("You said, \"%s\"", input))

	for _, other := range character.Room.Characters {
		if other != character {
			other.SendMessage(fmt.Sprintf("%s said, \"%s\"", character.Name, input))
		}
	}
}

func (world *World) MoveCharacter(character *Character, to *Room) {
	character.Room.RemoveCharacter(character)
	to.AddCharacter(character)
	character.SendMessage(to.Desc)
}
