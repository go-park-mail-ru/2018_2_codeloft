package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/authservice/auth"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"net/http"
	"time"
)

func GenerateCookie() *http.Cookie {
	val := uuid.NewV4().String()
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    val,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	return &cookie
}

type SessionManager struct {
	DB *sql.DB
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		DB: db,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *auth.Session) (*auth.SessionID, error) {
	fmt.Println("call Create", in)
	cookie := GenerateCookie()
	val := &auth.SessionID{ID: cookie.Value}
	_, err := sm.DB.Exec("insert into sessions(value, id) values ($1, $2) on CONFLICT do nothing", val.ID, in.UserID)
	if err != nil {
		log.Println("[ERROR] Create:", err)
		return nil, grpc.Errorf(codes.AlreadyExists, "User already Auth")
	}
	return val, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *auth.SessionID) (*auth.Session, error) {
	fmt.Println("call Check", in)
	row := sm.DB.QueryRow("select id from sessions where value = $1", in.ID)
	var userid int64
	err := row.Scan(&userid)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "auth not found")
	}
	return &auth.Session{UserID: userid}, nil
}

func (sm *SessionManager) Delete(ctx context.Context, in *auth.SessionID) (*auth.Nothing, error) {
	fmt.Println("call Delete", in)
	_, err := sm.DB.Exec("delete from sessions where value = $1", in.ID)
	if err != nil {

		log.Printf("cant DelCookie: %v\n", err)

		return nil, grpc.Errorf(codes.NotFound, "not found cookie")
	}
	return &auth.Nothing{Dummy: true}, nil
}
