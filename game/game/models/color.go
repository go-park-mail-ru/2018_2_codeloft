package models

import (
	"math/rand"
	"time"
)

const (
	COLOR_RED    = "#ff2d00"
	COLOR_BLUE   = "#006cff"
	COLOR_ORANGE = "#ff9300"
	COLOR_GREEN  = "#32ff00"
	COLOR_LIME   = "#00fff3"
	COLOR_YELLOW = "#ffff00"
	COLOR_WHITE  = "#ffffff"
	COLOR_VIOLET = "#c908b4"
	COLOR_BLACK  = "#000000"
)

var colorMap = map[int]string{
	0: COLOR_RED,
	1: COLOR_BLUE,
	2: COLOR_ORANGE,
	3: COLOR_GREEN,
	4: COLOR_LIME,
	5: COLOR_YELLOW,
	6: COLOR_WHITE,
	7: COLOR_VIOLET,
}

func GetRandomColor() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(8)
	color := colorMap[num]
	return color
}
