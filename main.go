package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/database"
	"github.com/go-park-mail-ru/2018_2_codeloft/handlers"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
)


func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("in URL: %v With method %v\n", r.URL.Path, r.Method)
				log.Println("recovered", err)

			}
		}()
		next.ServeHTTP(w, r)
	})
}

//TO DO
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("URL: %v; Method: %v; Origin: %v\n", r.URL.Path, r.Method, r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
	})
}


func main() {
	db := &database.DB{}
	if (len(os.Args) < 3){
		fmt.Println("Usage ./2018_2_codeloft <username> <password>")
		fmt.Println("Getting USERNAME and PASSWORD from env")
		var exist bool
		db.DB_USERNAME, exist = os.LookupEnv("USERNAME")
		if !exist {
			log.Println("USERNAME don't set")
		}
		db.DB_PASSWORD, exist = os.LookupEnv("PASSWORD")
		if !exist {
			log.Println("PASSWORD don't set")
		}
	} else
	{
		db.DB_USERNAME = os.Args[1]
		db.DB_PASSWORD = os.Args[2]
	}
	db.DB_NAME = "codeloft"
	db.DB_URL = os.Getenv("DATABASE_URL") // for heroku
	db.ConnectDataBase()
	defer db.DataBase.Close()
	var filepath string = "resources/initdb.sql"
	if _,err := os.Stat(filepath); err == nil {
		db.Init(filepath)
	} else {
		log.Printf("file %s does not exist\n", filepath)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.MainPage)
	mux.Handle("/user", &handlers.UserHandler{db.DataBase})
	mux.Handle("/session", &handlers.SessionHandler{db.DataBase})
	mux.Handle("/user/", &handlers.UserById{db.DataBase})

	fmt.Println("starting server on http://127.0.0.1:8080")
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "codeloft") ||
				strings.Contains(origin, "localhost") ||
				strings.Contains(origin, "127.0.0.1")
		},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
	})
	logHandler := logMiddleware(mux)
	corsMW := c.Handler(logHandler)
	panicMW := panicMiddleware(corsMW)
	port := os.Getenv("PORT") // for heroku
	if port != "" {
		log.Println("get port from env: ", port)
	} else {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, panicMW)
}
