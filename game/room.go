package game

import (
	"container/list"
	"encoding/json"
	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const MAXPLAYERS = 10

type Command struct {
	Player  *gamemodels.Player
	Command string
}

func NewRoom() *Room {
	id := uuid.Must(uuid.NewV4()).String()
	return &Room{
		ID:          id,
		MaxPlayers:  MAXPLAYERS,
		Players:     make(map[string]*PlayerConn),
		Connections: make(chan *PlayerConn),
		Disconnects: make(chan *PlayerConn),
		Broadcast:   make(chan *OutMessage),
		Message:     make(chan *IncomingMessage),
		Ticker:      time.NewTicker(time.Millisecond * 200),
	}
}

type Room struct {
	ID          string
	Ticker      *time.Ticker
	Players     map[string]*PlayerConn
	MaxPlayers  int
	Connections chan *PlayerConn
	Disconnects chan *PlayerConn
	Message     chan *IncomingMessage
	Broadcast   chan *OutMessage
	Commands    []*Command
}

type PlayerConn struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Player gamemodels.Player
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	PlayerCon  *PlayerConn     `json:"-"`
}

type OutMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type State struct {
	Players []gamemodels.Player `json:"players"`
}

func (r *Room) ListenToPlayers() {
	for {
		select {
		case m := <-r.Message:
			log.Printf("message from player %s: %v", m.PlayerCon.ID, string(m.Payload))

			switch m.Type {
			case "connect_player":
				m.PlayerCon.Player.Position.RandomPos()
				m.PlayerCon.Player.Tracer = list.New()
				m.PlayerCon.Player.Tracer.PushBack(m.PlayerCon.Player.Position)
			}

		case p := <-r.Disconnects:
			delete(r.Players, p.ID)
			if len(r.Players) == 0 {
				r.Ticker.Stop()
				game := GetGame()
				delete(game.Rooms, r.ID)
			}
			log.Printf("player was deleted from room %s", r.ID)
		}

	}
}


func (r *Room) Run() {

	go r.RunBroadcast()

	players := []gamemodels.Player{}
	for _, p := range r.Players {
		players = append(players, p.Player)
	}
	state := &State{
		Players: players,
	}
	r.Broadcast <- &OutMessage{Type: "SIGNAL_NEW_GAME_STATE", Payload: state}
	for {
		<-r.Ticker.C
		log.Printf("room %s tick with %d players", r.ID, len(r.Players))

		players := []gamemodels.Player{}
		for _, p := range r.Players {
			players = append(players, p.Player)
		}

		state := &State{
			Players: players,
		}

		r.Broadcast <- &OutMessage{Type: "IN_GAME", Payload: state}
	}
}


func (r *Room) RunBroadcast() {
	for {
		m := <-r.Broadcast
		for _, p := range r.Players {
			p.Send(m)
		}
	}
}

func (p *PlayerConn) Send(s *OutMessage) {
	err := p.Conn.WriteJSON(s)
	if err != nil {
		log.Printf("cannot send state to client: %s", err)
	}
}

func (p *PlayerConn) Listen() {
	log.Printf("start listening messages from player %s", p.ID)

	for {

		m := &IncomingMessage{}

		err := p.Conn.ReadJSON(m)
		//_, b, err := p.Conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %s was disconnected", p.ID)
			p.Room.Disconnects <- p
			return
		} else if err != nil {
			log.Println("Error READJSON in Listen")
			continue
		}
		//fmt.Println(string(b))
		m.PlayerCon = p
		p.Room.Message <- m

	}
}