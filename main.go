package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2018_2_codeloft/database"
	"github.com/go-park-mail-ru/2018_2_codeloft/handlers"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"log"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2018_2_codeloft/logger"
	"github.com/rs/cors"
	"github.com/go-park-mail-ru/2018_2_codeloft/auth"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)


var (
	dbhost = "127.0.0.1"
	authhost = "127.0.0.1"
	mongohost = "127.0.0.1"
	databasename = "codeloft"
)



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

//TO DO
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("REQUEST",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleWare(next http.Handler, db *sql.DB, sm auth.AuthCheckerClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//var s *models.Session
		//if s = services.GetCookie(r, db); s == nil {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}
		//var user models.User
		//if !user.GetUserByID(db, s.User_id) {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	log.Println("User Does Not Exist in Users table, but exist in session", s.Value, s.User_id)
		//	return
		//}
		cookie, err := r.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("No cookie header with session_id name", err)
			return
		}
		userid, err := sm.Check(context.Background(), &auth.SessionID{ID: cookie.Value})
		if err != nil {
			fmt.Println("[ERROR] checkAuth:", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var user models.User
		if !user.GetUserByID(db, userid.UserID) {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("User Does Not Exist in Users table, but exist in session", cookie.Value, userid.UserID)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "login", user.Login)
		next.ServeHTTP(w, r.WithContext(ctx))
		//next.ServeHTTP(w,r)
	})
}

func main() {
	if os.Getenv("ENV") == "production" {
		dbhost = "db"
		authhost = "auth"
		mongohost = "mongodb"
	}
	zapLogger, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Can not initialize zap logger Error %v", err)
	}
	defer zapLogger.Sync()

	dbUserName := ""
	dbPassword := ""
	db := &database.DB{}
	if len(os.Args) < 3 {
		fmt.Println("Usage ./2018_2_codeloft <username> <password>")
		fmt.Println("Getting USERNAME and PASSWORD from env")
		var exist bool
		db.DB_USERNAME, exist = os.LookupEnv("USERNAME")
		if !exist {
			zap.L().Info("USERNAME don't set")
		}
		db.DB_PASSWORD, exist = os.LookupEnv("PASSWORD")
		if !exist {
			zap.L().Info("PASSWORD don't set")
		}
	} else {
		db.DB_USERNAME = os.Args[1]
		db.DB_PASSWORD = os.Args[2]
		dbUserName = os.Args[1]
		dbPassword = os.Args[2]
	}
	db.DB_NAME = databasename
	db.DB_URL = os.Getenv("DATABASE_URL") // for heroku
	db.ConnectDataBase()
	defer db.DataBase.Close()
	var filepath string = "resources/initdb.sql"
	if _, err := os.Stat(filepath); err == nil {
		db.Init(filepath)
	} else {
		zap.S().Warn("file does not exist\n", filepath)
	}

	//gameMux := http.NewServeMux()
	//gameMux.Handle("/gamews", &handlers.GameHandler{db.DataBase})
	//authHandler := AuthMiddleWare(gameMux, db.DataBase)

	log.Println("Connecting to MongoDB:")
	mongoDb := &database.MongoDB{}
	mongoDb.DB_USERNAME = dbUserName
	mongoDb.DB_PASSWORD = dbPassword
	mongoDb.DB_NAME = databasename
	mongoDb.DB_URL = fmt.Sprintf("mongodb://%s:%s@%s:27017/%s", mongoDb.DB_USERNAME,
		mongoDb.DB_PASSWORD,
		mongohost,
		mongoDb.DB_NAME,
	)
	err = mongoDb.Connect()
	if err != nil {
		log.Println("[ERROR] MognoConnection:",err)
	}

	grcpConn, err := grpc.Dial(
		fmt.Sprintf("%s:8081",authhost),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := auth.NewAuthCheckerClient(grcpConn)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.MainPage)
	mux.Handle("/user", &handlers.UserHandler{db.DataBase, sessManager})
	mux.Handle("/session", &handlers.SessionHandler{db.DataBase, sessManager})
	mux.Handle("/user/", &handlers.UserById{db.DataBase, sessManager})
	//mux.Handle("/gamews", authHandler)
	mux.Handle("/gamews", &handlers.GameHandler{db.DataBase})
	mux.Handle("/chatws", &handlers.ChatHandler{mongoDb})
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
		zap.S().Infow("get port from env: ", port)
	} else {
		port = "8080"
	}

	if len(os.Args) > 3 {
		port = os.Args[3]
		fmt.Println(port)
	}
	addr := fmt.Sprintf(":%s", port)

	fmt.Println("starting server on http://127.0.0.1:8080")
	http.ListenAndServe(addr, panicMW)

}
