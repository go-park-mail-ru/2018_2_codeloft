package models

import (
	"github.com/satori/go.uuid"
	"time"
)

const MAXPLAYERS = 10

type Command struct {
	Player  *Player
	Command string
}

type NewPlayer struct {
	Username string `json:"username"`
}

func NewRoom() *Room {
	id := uuid.Must(uuid.NewV4()).String()
	return &Room{
		ID:         id,
		MaxPlayers: MAXPLAYERS,
		Players:    make(map[string]*Player),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Broadcast:  make(chan *Message),
		Message:    make(chan *IncomingMessage),
	}
}

type Room struct {
	ID         string
	Ticker     *time.Ticker
	Players    map[string]*Player
	MaxPlayers int
	Register   chan *Player
	Unregister chan *Player
	Message    chan *IncomingMessage
	Broadcast  chan *Message
	Commands   []*Command
}