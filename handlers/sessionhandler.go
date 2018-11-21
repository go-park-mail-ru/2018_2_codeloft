package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"github.com/go-park-mail-ru/2018_2_codeloft/services"
	"github.com/go-park-mail-ru/2018_2_codeloft/validator"
	"go.uber.org/zap"
)

func checkAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var s *models.Session
	if s = services.GetCookie(r, db); s == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// cookie, err := r.Cookie("session_id")
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// var s models.Session
	// if !s.CheckCookie(db, cookie.Value) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	var user models.User
	if !user.GetUserByID(db, s.User_id) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(generateError(models.MyError{r.URL.Path, "User Does Not Exist in Users table, but exist in session", fmt.Errorf("")}))
		zap.L().Info("User Does Not Exist in Users table, but exist in session",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.String("Session value", s.Value),
			zap.Int64("User id", s.User_id),
		)
		return
	}
	res, err := json.Marshal(&user)
	if err != nil {
		zap.L().Info("error while Marshaling in /user",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func signIn(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var s *models.Session
	// Если уже залогинен
	if s = services.GetCookie(r, db); s != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error while reading body in /session", err)

		zap.L().Info("error while reading body in /session",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var u struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path, "wrong requst format", err}))
		zap.L().Info("error while Marshaling in /session",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)

		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path, "bad login", err}))
		zap.L().Info("error while validating in /session",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)
		return
	}

	var dbUser models.User
	if !dbUser.GetUserByLogin(db, u.Login) {
		w.WriteHeader(http.StatusBadRequest)
		err := models.MyError{r.URL.Path, "User does not exist", models.UserDoesNotExist(u.Login)}
		w.Write(generateError(err))
		zap.L().Info("User does not exist",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(&err),
		)
		return
	}
	if dbUser.Password != u.Password {
		w.WriteHeader(http.StatusBadRequest)

		err := models.MyError{r.URL.Path, "wrong password", fmt.Errorf("wrong password")}
		w.Write(generateError(err))
		zap.L().Info("Wrong password",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(&err),
		)

		w.Write(generateError(models.MyError{r.URL.Path, "wrong password", fmt.Errorf("wrong password")}))

		return
	}
	// cookie := http.Cookie{
	// 	Name:     "session_id",
	// 	Value:    "testCookie"+dbUser.Login,
	// 	Expires:  time.Now().Add(30 * 24 * time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   false,
	// }
	cookie := services.GenerateCookie(dbUser.Login)
	if os.Getenv("ENV") == "production" {
		cookie.Secure = true
	}
	s = &models.Session{cookie.Value, dbUser.Id}
	err = s.AddCookie(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		myErr := models.MyError{r.URL.Path, "Cant AddCookie", err}
		w.Write(generateError(myErr))
		zap.L().Info("Cant AddCookie",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(&myErr),
		)

		w.Write(generateError(models.MyError{r.URL.Path, "Cant AddCookie", err}))

		return
	}
	http.SetCookie(w, cookie)
	res, err := json.Marshal(&dbUser)
	if err != nil {
		zap.L().Info("error while Marshaling",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// cookie, err := r.Cookie("session_id")
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// var s models.Session
	// if !s.CheckCookie(db, cookie.Value) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	var s *models.Session
	if s = services.GetCookie(r, db); s == nil {
		zap.L().Info("StatusConflist",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Int("Code", http.StatusConflict),
		)
		w.WriteHeader(http.StatusConflict)
		return
	}
	cookie, _ := r.Cookie("session_id")
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	err := s.DeleteCookie(db)
	if err != nil {
		zap.L().Warn("Can not delete cookie",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(err),
		)
	}
	w.WriteHeader(http.StatusOK)
}

type SessionHandler struct {
	Db *sql.DB
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {

	case http.MethodGet:
		checkAuth(w, r, h.Db)
	case http.MethodPost:
		signIn(w, r, h.Db)
	case http.MethodDelete:
		logout(w, r, h.Db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//var SessionHandler = func(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	switch r.Method {
//
//	case http.MethodGet:
//		checkAuth(w, r)
//	case http.MethodPost:
//		signIn(w, r)
//	case http.MethodDelete:
//		logout(w, r)
//
//	}
//
//}
