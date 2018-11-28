package game

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(RoomsCount)
	game := GetGame()
	go game.Run()
}
