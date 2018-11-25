package handlers

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_codeloft/chat"
	"github.com/go-park-mail-ru/2018_2_codeloft/database"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	Db *database.MongoDB
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("cannot upgrade connection: %s", err)
		return
	}

	//defer conn.Close()
	ctx := r.Context()
	login := ctx.Value("login")
	log.Println("login from context:", login)
	//conn.WriteJSON(login)
	//game.Connect(conn, login.(string))
	chat.Connect(conn)
}
