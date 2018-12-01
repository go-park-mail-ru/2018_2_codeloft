package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/authservice/auth"
	"github.com/go-park-mail-ru/2018_2_codeloft/game/game"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

type GameHandler struct {
	Sm auth.AuthCheckerClient
}

func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//cookie, err := r.Cookie("session_id")
	//if err != nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	log.Println("No cookie header with session_id name", err)
	//	return
	//}
	//userid, err := h.Sm.Check(context.Background(), &auth.SessionID{ID: cookie.Value})
	//if err != nil {
	//	fmt.Println("[ERROR] checkAuth:", err)
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

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

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Errorw("Recovered",
					"URL", r.URL.Path,
					"Method", r.Method,
					"Origin", r.Header.Get("Origin"),
					"Remote address", r.RemoteAddr,
					"Error", err,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var (
	authhost = "127.0.0.1"
)

func main() {
	if os.Getenv("ENV") == "production" {
		authhost = "auth"
	}
	log.Println("connect to authservice")
	grcpConn, err := grpc.Dial(
		fmt.Sprintf("%s:8081", authhost),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := auth.NewAuthCheckerClient(grcpConn)

	mux := http.NewServeMux()
	mux.Handle("/gamews", &GameHandler{sessManager})
	mux.Handle("/metrics", prometheus.Handler())
	pw := panicMiddleware(mux)

	log.Println("start gameserver :8082")
	http.ListenAndServe(":8082", pw)
}
