package game

import (
	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
)

func init() {
	game := gamemodels.NewGame()
	go game.Run()
}

