package handlers

import (
	"database/sql"
	"github.com/go-park-mail-ru/2018_2_codeloft/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type GameHandler struct {
	Db *sql.DB
}

func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
	game.Connect(conn)
}
