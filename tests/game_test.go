package tests

import (
	"github.com/go-park-mail-ru/2018_2_codeloft/game/game"
	"reflect"
	"testing"
)

func TestSingletonGame(t *testing.T) {
	firstGame := game.GetGame()
	secondGame := game.GetGame()
	if !reflect.DeepEqual(firstGame, secondGame) {
		t.Errorf("GetGame return different games: \n%v\n%v", firstGame, secondGame)
	}
}

func TestFindRoomSame(t *testing.T) {
	game := game.GetGame()
	r1 := game.FindRoom()
	r2 := game.FindRoom()
	switch {
	case game.MaxRooms < 2:
		return
	case game.MaxRooms > 2:
		if !reflect.DeepEqual(r1, r2) {
			t.Errorf("FindRoom find different rooms while size < game.MaxRooms: \n%v\n%v", r1, r2)
		}
	}
}


