package chat

import (
	"log"
	"sync"

	"github.com/go-park-mail-ru/2018_2_codeloft/chat/models"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type Chat struct {
	//Connections chan *connectInfo
	Users       map[string]*UserConn
	Connections chan *websocket.Conn
}

var globalChat *Chat
var once sync.Once

func GetChat() *Chat {
	once.Do(func() {
		globalChat = &Chat{
			//Connections: make(chan *connectInfo),
			Users:       make(map[string]*UserConn),
			Connections: make(chan *websocket.Conn),
		}
	})
	return globalChat
}

func Connect(conn *websocket.Conn) {
	globalChat.Connections <- conn
}

func (g *Chat) Run() {
	for {
		conn := <-g.Connections
		log.Printf("got new connection")
		//g.ProcessConn(conn.ws, conn.nickname)
		g.ProcessConn(conn)
	}
}

func SendMessage(m *models.Message) {
	for _, conn := range globalChat.Users {
		switch conn.UserLogin {
		case "", m.SenderLogin, m.ReceiverLogin:
			conn.Send(m)
		}
	}
}

func init() {
	chat := GetChat()
	go chat.Run()
}

func (g *Chat) ProcessConn(conn *websocket.Conn) {
	id := uuid.NewV4().String()
	var nickname string
	err := conn.ReadJSON(&nickname)
	if err != nil {
		log.Println(err)
	}
	p := &UserConn{
		Conn:      conn,
		ID:        id,
		UserLogin: nickname,
		//Player: &gamemodels.Player{Speed: gamemodels.DEFAULT_SPEED},
	}
	globalChat.Users[p.ID] = p

	log.Printf("player %s joined", p.UserLogin)
	go p.Listen()
}
