package models

import (
	"math/rand"
	"time"
)

//const width = 320
//const height = 180

const (
	scale         = 5
	FIELD_WIDTH   = 16 * scale
	FIELD_HEIGHT  = 9 * scale
	DEFAULT_SPEED = 100 //количество милисекунд до обновления координат игрока
)

//easyjson:/json
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//type Direction struct {
//	X int
//	Y int
//}

var Directions = map[string]Position{"DOWN": {0, 1}, "UP": {0, -1}, "RIGHT": {1, 0}, "LEFT": {-1, 0}}

func (p *Position) RandomPos() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	p.X = r1.Intn(FIELD_WIDTH)
	p.Y = r1.Intn(FIELD_HEIGHT)
}

//easyjson:json
type Player struct {
	Username      string       `json:"username"`
	Position      Position     `json:"position"`
	Tracer        []Position   `json:"-"`
	SpeedTicker   *time.Ticker `json:"-"`
	Speed         int          `json:"speed"`
	MoveDirection string       `json:"move_direction"`
	Score         int          `json:"score"`
	ID            int          `json:"-"`
	IsDead        bool         `json:"is_dead,omitempty"`
	Color         string       `json:"color"`
}

func (p *Player) ChangeDirection(direction string) {
	if _, exist := Directions[direction]; !exist {
		return
	}
	if p.MoveDirection == "DOWN" && direction == "UP" {
		return
	}
	if p.MoveDirection == "UP" && direction == "DOWN" {
		return
	}
	if p.MoveDirection == "RIGHT" && direction == "LEFT" {
		return
	}
	if p.MoveDirection == "LEFT" && direction == "RIGHT" {
		return
	}
	p.MoveDirection = direction
}

func (p *Player) Move() {

	p.Position.X += Directions[p.MoveDirection].X
	p.Position.Y += Directions[p.MoveDirection].Y
	if p.Position.X >= FIELD_WIDTH-1 {
		p.Position.X = 0
	}
	if p.Position.Y >= FIELD_HEIGHT-1 {
		p.Position.Y = 0
	}
	if p.Position.X < 0 {
		p.Position.X = FIELD_WIDTH - 1
	}
	if p.Position.Y < 0 {
		p.Position.Y = FIELD_HEIGHT - 1
	}
	p.Tracer = append(p.Tracer, p.Position)

}
