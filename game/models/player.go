package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerData struct {
	Username string `json:"username"`
	HP       string
	Position Position `json:"position"`
}

type Player struct {
	ID   string
	Room *Room
	Conn *websocket.Conn
	Data PlayerData
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Player  *Player         `json:"-"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}