package models

import (
	"github.com/satori/go.uuid"
	"log"
	"github.com/gorilla/websocket"
)

const MAXROOMS = 5

func NewGame() *Game {
	return &Game{
		Rooms:    make(map[string]*Room),
		MaxRooms: MAXROOMS,
		Register: make(chan *websocket.Conn),
	}
}

type Game struct {
	Rooms    map[string]*Room
	MaxRooms int
	Register chan *websocket.Conn
}

func (g *Game) Run() {
	for {
		conn := <-g.Register
		log.Printf("got new connection")

		g.ProcessConn(conn)
	}
}

func (g *Game) FindRoom() *Room {
	for _, r := range g.Rooms {
		if len(r.Players) < r.MaxPlayers {
			return r
		}
	}

	if len(g.Rooms) >= g.MaxRooms {
		return nil
	}

	r := NewRoom()
	go r.ListenToPlayers()
	go r.Run()
	g.Rooms[r.ID] = r
	log.Printf("room %s created", r.ID)

	return r
}

func (g *Game) ProcessConn(conn *websocket.Conn) {
	id := uuid.Must(uuid.NewV4()).String()
	p := &Player{
		Conn: conn,
		ID:   id,
	}

	r := g.FindRoom()
	if r == nil {
		return
	}
	r.Players[p.ID] = p
	p.Room = r
	log.Printf("player %s joined room %s", p.ID, r.ID)
	go p.Listen()

}