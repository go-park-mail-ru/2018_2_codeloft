package game

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"unsafe"

	gamemodels "github.com/go-park-mail-ru/2018_2_codeloft/game/models"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

const (
	SIGNAL_CONNECT = "connect"
	SIGNAL_DEAD    = "dead"
	NO_SIGNAL      = "None"
)

const (
	MAXPLAYERS = 10
)

func NewRoom() *Room {
	id := uuid.NewV4().String()
	field := [gamemodels.FIELD_HEIGHT][gamemodels.FIELD_WIDTH]gamemodels.Cell{}
	for i := 0; i < gamemodels.FIELD_HEIGHT; i++ {
		for j := 0; j < gamemodels.FIELD_WIDTH; j++ {
			field[i][j] = gamemodels.Cell{Val: 0}
		}
	}
	//fmt.Println(field)
	return &Room{
		ID:          id,
		MaxPlayers:  MAXPLAYERS,
		Players:     make(map[int]*PlayerConn),
		Connections: make(chan *PlayerConn),
		Disconnects: make(chan *PlayerConn),
		Broadcast:   make(chan *OutMessage),
		Message:     make(chan *IncomingMessage),
		Ticker:      time.NewTicker(time.Millisecond * 100),
		Field:       field,
		LastId:      1,
	}
}

type Room struct {
	ID          string
	Ticker      *time.Ticker
	Players     map[int]*PlayerConn
	MaxPlayers  int
	Connections chan *PlayerConn
	Disconnects chan *PlayerConn
	Message     chan *IncomingMessage
	Broadcast   chan *OutMessage
	Field       [gamemodels.FIELD_HEIGHT][gamemodels.FIELD_WIDTH]gamemodels.Cell
	LastId      int
}

type PlayerConn struct {
	//ID   string
	ID     int
	Room   *Room
	Conn   *websocket.Conn
	Player *gamemodels.Player
	Signal string
}

//easyjson:json
type IncomingMessage struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	PlayerCon *PlayerConn     `json:"-"`
}

//easyjson:json
type OutMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

//easyjson:json
type State struct {
	Players []gamemodels.Player `json:"players"`
}

func (r *Room) ListenToPlayers() {
	for {
		select {
		case m := <-r.Message:
			log.Printf("message from player %s: %v", m.PlayerCon.Player.Username, string(m.Payload))

			switch m.Type {
			case "connect_player":
				m.PlayerCon.Signal = SIGNAL_CONNECT
				player := m.PlayerCon.Player
				player.IsDead = false
				player.Score = 0
				player.Position.RandomPos()
				m.PlayerCon.Player.Tracer = make([]gamemodels.Position, 0, 20)
				//r.Field[player.Position.Y][player.Position.X].Mu.Lock()
				for {
					if r.Field[player.Position.Y][player.Position.X].Val == 0 {
						r.Field[player.Position.Y][player.Position.X].Val = m.PlayerCon.Player.ID
						player.Tracer = append(player.Tracer, player.Position)
						break
					}
					player.Position.RandomPos()
				}
				//r.Field[player.Position.Y][player.Position.X].Val = m.PlayerCon.ID
				//r.Field[player.Position.Y][player.Position.X].Mu.Unlock()
				player.MoveDirection = "RIGHT"
			case "change_direction":
				direction := ""
				json.Unmarshal(m.Payload, &direction)
				m.PlayerCon.Player.ChangeDirection(direction)
				//fmt.Printf("Player %s, change direction to %v\n", m.PlayerCon.Player.Username, m.PlayerCon.Player.MoveDirection)
			}

		case p := <-r.Disconnects:
			for _, pos := range p.Player.Tracer {
				r.Field[pos.Y][pos.X].Val = 0
			}
			delete(r.Players, p.ID)
			if len(r.Players) == 0 {
				r.Ticker.Stop()
				game := GetGame()
				delete(game.Rooms, r.ID)
				log.Printf("Room %s was deleted", r.ID)
			}
			log.Printf("player was deleted from room %s", r.ID)
		}

	}
}

func (r *Room) Run() {

	go r.RunBroadcast()

	//players := make([]gamemodels.Player, 0, len(r.Players))
	//for _, p := range r.Players {
	//	players = append(players, *p.Player)
	//}
	//state := &State{
	//	Players: players,
	//}
	//r.Broadcast <- &OutMessage{Type: "SIGNAL_NEW_GAME_STATE", Payload: state}
	for {
		<-r.Ticker.C
		log.Printf("room %s tick with %d players", r.ID, len(r.Players))

		players := make([]gamemodels.Player, 0, len(r.Players))
		for _, p := range r.Players {
			players = append(players, *p.Player)
		}

		state := &State{
			Players: players,
		}

		r.Broadcast <- &OutMessage{Type: "IN_GAME", Payload: state}
		//fmt.Println(r.Field)
		r.MovePlayers()
	}
}

func (r *Room) MovePlayers() {
	for _, p := range r.Players {
		startpos := p.Player.Position
		if p.Player.IsDead == true {
			continue
		}
		p.Player.Move()
		for startpos.Y < p.Player.Position.Y && startpos.X < p.Player.Position.X {
			startpos.Y += gamemodels.Directions[p.Player.MoveDirection].Y
			startpos.X += gamemodels.Directions[p.Player.MoveDirection].X
			r.Field[startpos.Y][startpos.X].Val = p.ID
		}
		//r.Field[p.Player.Position.Y][p.Player.Position.X].Mu.Lock()
		if r.Field[p.Player.Position.Y][p.Player.Position.X].Val == 0 {
			r.Field[p.Player.Position.Y][p.Player.Position.X].Val = p.Player.ID
		} else {
			p.Player.IsDead = true
			for _, pos := range p.Player.Tracer {
				r.Field[pos.Y][pos.X].Val = 0
			}
			p.Player.Position = gamemodels.Position{-1, -1}
		}
		//r.Field[p.Player.Position.Y][p.Player.Position.X].Mu.Unlock()
	}
}

func (r *Room) RunBroadcast() {
	for {
		m := <-r.Broadcast
		for _, p := range r.Players {
			if p.Signal == SIGNAL_CONNECT {
				//log.Println(r.Field)
				p.Send(&OutMessage{"connected", r.Field})

			}
			if p.Player.IsDead != false {
				p.Send(&OutMessage{"DEAD", p.Player.Score})
			} else {
				p.Send(m)
			}
			p.Signal = NO_SIGNAL
		}
	}
}

func (p *PlayerConn) Send(s *OutMessage) {
	d, _ := json.Marshal(s)
	fmt.Println(unsafe.Sizeof(d))
	err := p.Conn.WriteJSON(s)
	if err != nil {
		log.Printf("cannot send state to client: %s", err)
	}
}

func (p *PlayerConn) Listen() {
	log.Printf("start listening messages from player %s", p.ID)

	initMessage := &IncomingMessage{"connect_player", json.RawMessage{}, p}
	p.Room.Message <- initMessage
	// p.Conn.WriteJSON(p.Room.Field) // send matrix
	for {

		m := &IncomingMessage{}

		err := p.Conn.ReadJSON(m)
		//_, b, err := p.Conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("player %s was disconnected", p.Player.Username)
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
