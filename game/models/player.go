package models

import (
	"container/list"
	"math/rand"
	"time"
)

const width = 800
const height = 600

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Position) RandomPos() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	p.X = r1.Intn(width)
	p.Y = r1.Intn(height)
}

type Player struct {
	Username string   `json:"username"`
	Position Position `json:"position"`
	Tracer *list.List `json:"tracer"`
}

