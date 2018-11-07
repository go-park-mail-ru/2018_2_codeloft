package game

import (
	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"sync"
)

const MAXROOMS = 5

var globalGame *Game
var once sync.Once

func GetGame() *Game {
	once.Do(func() {
		globalGame = &Game{
			Rooms:       make(map[string]*Room),
			MaxRooms:    MAXROOMS,
			Connections: make(chan *websocket.Conn),
		}
	})
	return globalGame
}

type Game struct {
	Rooms       map[string]*Room
	MaxRooms    int
	Connections chan *websocket.Conn
}

func Connect(conn *websocket.Conn) {
	globalGame.Connections <- conn
}

func (g *Game) Run() {
	for {
		conn := <-g.Connections
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
	var username string
	err := conn.ReadJSON(&username)
	if err != nil {
		log.Println("error while reading json in processConn")
		username = ""
	}
	p := &PlayerConn{
		Conn: conn,
		ID:   id,
		Player: gamemodels.Player{Username: username},
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
