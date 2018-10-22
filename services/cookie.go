package services

import (
	"golang.org/x/crypto/sha3"
	"fmt"
	"net/http"
	"time"
	"database/sql"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"os"
)

func CheckCookie(r *http.Request, db *sql.DB) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	var s models.Session
	if !s.CheckCookie(db, cookie.Value) {
		return false
	}
	if cookie.Expires < time.Now() {
		return false
	}
	return true
}

func GenerateCookie(val string) *http.Cookie {
	buf := []byte(val+os.Getenv("USERNAME"))
	h := make([]byte, 64)
	sha3.ShakeSum256(h, buf)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    string(h),
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	return *http.Cookie
}

