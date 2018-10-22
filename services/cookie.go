package services

import (
	"golang.org/x/crypto/sha3"
	"net/http"
	"fmt"
	"time"
	"database/sql"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"os"
)

func GetCookie(s *models.Session,r *http.Request, db *sql.DB) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	if !s.CheckCookie(db, cookie.Value) {
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
		Value:    fmt.Sprintf("%x", h),
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	return &cookie
}
