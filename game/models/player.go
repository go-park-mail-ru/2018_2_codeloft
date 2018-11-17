package models

import (
	"math/rand"
	"time"
)

//const width = 320
//const height = 180

const (
	FIELD_WIDTH   = 16
	FIELD_HEIGHT  = 9
	DEFAULT_SPEED = 1
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//type Direction struct {
//	X int
//	Y int
//}

var Directions = map[string]Position{"BOTTOM": {0, 1}, "TOP": {0, -1}, "RIGHT": {1, 0}, "LEFT": {-1, 0}}

func (p *Position) RandomPos() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	p.X = r1.Intn(FIELD_WIDTH)
	p.Y = r1.Intn(FIELD_HEIGHT)
}

type Player struct {
	Username      string     `json:"username"`
	Position      Position   `json:"position"`
	Tracer        []Position `json:"-"`
	Speed         int        `json:"speed"`
	MoveDirection string     `json:"move_direction"`
	Score         int        `json:"score"`
}

func (p *Player) ChangeDirection(direction string) {
	if _, exist := Directions[direction]; !exist {
		return
	}
	p.MoveDirection = direction
}

func (p *Player) Move() {
	p.Position.X += p.Speed * Directions[p.MoveDirection].X
	p.Position.Y += p.Speed * Directions[p.MoveDirection].Y
	if p.Position.X >= FIELD_WIDTH-1 {
		p.Position.X = 1
	}
	if p.Position.Y >= FIELD_HEIGHT-1 {
		p.Position.Y = 1
	}
	if p.Position.X < 0 {
		p.Position.X = FIELD_WIDTH - 2
	}
	if p.Position.Y < 0 {
		p.Position.Y = FIELD_HEIGHT - 2
	}
}
