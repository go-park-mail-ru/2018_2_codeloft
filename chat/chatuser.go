package chat

import (
	"log"

	"github.com/go-park-mail-ru/2018_2_codeloft/chat/models"
	"github.com/gorilla/websocket"
)

type UserConn struct {
	ID        string
	Conn      *websocket.Conn
	UserLogin string
}

func (u *UserConn) Send(m *models.Message) {
	err := u.Conn.WriteJSON(m)
	if err != nil {
		log.Printf("cannot send state to client: %s", err)
	}
}

func (p *UserConn) Listen() {
	log.Printf("start listening messages from player %s", p.UserLogin)

	// p.Conn.WriteJSON(p.Room.Field) // send matrix
	for {

		m := &models.Message{}

		err := p.Conn.ReadJSON(m)
		//_, b, err := p.Conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("user %s was disconnected", p.UserLogin)
			delete(globalChat.Users, p.ID)
			return
		} else if err != nil {
			log.Println("Error READJSON in Listen")
			log.Println(err)
			continue
		}
		m.Type = "user_message"
		log.Println(m)
		go SendMessage(m)
	}
}
