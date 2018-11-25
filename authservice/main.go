//client

package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/authservice/auth"
	"github.com/go-park-mail-ru/2018_2_codeloft/authservice/database"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	dbhost = "127.0.0.1"
	authhost = "127.0.0.1"
)


func main() {
	if os.Getenv("ENV") == "production" {
		dbhost = "db"
		authhost = "auth"
	}
	db := &database.DB{}
	if len(os.Args) < 3 {
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
	} else {
		db.DB_USERNAME = os.Args[1]
		db.DB_PASSWORD = os.Args[2]
	}
	db.DB_NAME = "codeloft"
	db.ConnectDataBase()
	defer db.DataBase.Close()
	val := ""
	err := db.DataBase.QueryRow("select value from sessions where id = 1").Scan(&val)
	if err != nil {
		fmt.Println("cant get val", err)
	}

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()

	auth.RegisterAuthCheckerServer(server, NewSessionManager(db.DataBase))

	fmt.Println("starting server at :8081")
	server.Serve(lis)

}
