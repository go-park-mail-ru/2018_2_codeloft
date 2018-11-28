package game

import (
	"log"
	"sync"

	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
)

const MAXROOMS = 5

var globalGame *Game
var once sync.Once
var RoomsCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "game_rooms_count",
	Help: "Count of rooms in game",
})

func GetGame() *Game {
	once.Do(func() {
		globalGame = &Game{
			Rooms:    make(map[string]*Room),
			MaxRooms: MAXROOMS,
			//Connections: make(chan *connectInfo),
			Connections: make(chan *websocket.Conn),
		}
	})
	return globalGame
}

type connectInfo struct {
	ws       *websocket.Conn
	nickname string
}

type Game struct {
	Rooms    map[string]*Room
	MaxRooms int
	//Connections chan *connectInfo
	Connections chan *websocket.Conn
}

//func Connect(conn *websocket.Conn, nickname string) {
//	globalGame.Connections <- &connectInfo{conn, nickname}
//}
func Connect(conn *websocket.Conn) {
	globalGame.Connections <- conn
}

func (g *Game) Run() {
	for {
		conn := <-g.Connections
		log.Printf("got new connection")
		//g.ProcessConn(conn.ws, conn.nickname)
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
	RoomsCount.Inc()
	go r.ListenToPlayers()
	go r.Run()
	g.Rooms[r.ID] = r
	log.Printf("room %s created", r.ID)

	return r
}

//func (g *Game) ProcessConn(conn *websocket.Conn, nickname string) {
func (g *Game) ProcessConn(conn *websocket.Conn) {
	//id := uuid.NewV4().String()
	var nickname string
	err := conn.ReadJSON(&nickname)
	if err != nil {
		log.Println(err)
	}
	p := &PlayerConn{
		Conn: conn,
		//ID:   id,
		Player: &gamemodels.Player{Username: nickname, Speed: gamemodels.DEFAULT_SPEED},
		//Player: &gamemodels.Player{Speed: gamemodels.DEFAULT_SPEED},
	}
	r := g.FindRoom()
	if r == nil {
		return
	}
	p.ID = r.LastId
	p.Player.ID = r.LastId
	r.LastId++
	r.Players[p.ID] = p
	p.Room = r
	log.Printf("player %s joined room %s", p.Player.Username, r.ID)
	go p.Listen()
}
