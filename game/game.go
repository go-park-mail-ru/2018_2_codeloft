package game

import (
	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
	"github.com/gorilla/websocket"
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
			Connections: make(chan *connectInfo),
		}
	})
	return globalGame
}

type connectInfo struct {
	ws       *websocket.Conn
	nickname string
}

type Game struct {
	Rooms       map[string]*Room
	MaxRooms    int
	Connections chan *connectInfo
}

func Connect(conn *websocket.Conn, nickname string) {
	globalGame.Connections <- &connectInfo{conn, nickname}
}

func (g *Game) Run() {
	for {
		conn := <-g.Connections
		log.Printf("got new connection")
		g.ProcessConn(conn.ws, conn.nickname)
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

func (g *Game) ProcessConn(conn *websocket.Conn, nickname string) {
	//id := uuid.NewV4().String()
	p := &PlayerConn{
		Conn: conn,
		//ID:   id,
		Player: &gamemodels.Player{Username: nickname, Speed: gamemodels.DEFAULT_SPEED},
	}
	r := g.FindRoom()
	if r == nil {
		return
	}
	p.ID = len(r.Players) + 1
	r.Players[p.ID] = p
	p.Room = r
	log.Printf("player %s joined room %s", p.Player.Username, r.ID)
	go p.Listen()
}
