// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v2
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> http://john.doe.com
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//     oauth2:
//         type: oauth2
//         authorizationUrl: /oauth2/auth
//         tokenUrl: /oauth2/token
//         in: header
//         scopes:
//           bar: foo
//         flow: accessCode
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
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
)

// swagger:operation GET /session checkAuth
//
// Checks auth with cookie
// ---
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - name: name
//   in: path
//   description: Name to be returned.
//   required: true
//   type: string
// responses:
//   '200':
//     description: The hello message
//     type: string
func checkAuth(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	s := &models.Session{}
	if !services.GetCookie(s, r, db) {
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
		log.Println("User Does Not Exist in Users table, but exist in session", s.Value, s.User_id)
		return
	}
	res, err := json.Marshal(&user)
	if err != nil {
		log.Println("error while Marshaling in /user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// swagger:operation POST /session signIn
//
// Logins
// ---
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - name: name
//   in: path
//   description: Name to be returned.
//   required: true
//   type: string
// responses:
//   '200':
//     description: The hello message
//     type: string
func signIn(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	s := &models.Session{}
	// Если уже залогинен
	if services.GetCookie(s, r, db) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading body in /session", err)
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
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path, "bad login", err}))
		return
	}
	var dbUser models.User
	if !dbUser.GetUserByLogin(db, u.Login) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path, "User does not exist", models.UserDoesNotExist(u.Login)}))
		return
	}
	if dbUser.Password != u.Password {
		w.WriteHeader(http.StatusBadRequest)
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
		w.Write(generateError(models.MyError{r.URL.Path, "Cant AddCookie", err}))
		return
	}
	http.SetCookie(w, cookie)
	res, err := json.Marshal(&dbUser)
	if err != nil {
		log.Println("error while Marshaling in /session POST")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// swagger:operation DELETE /session LogOut
//
// LogOut and delete cookie
// ---
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - name: name
//   in: path
//   description: Name to be returned.
//   required: true
//   type: string
// responses:
//   '200':
//     description: The hello message
//     type: string
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
	s := &models.Session{}
	if !services.GetCookie(s, r, db) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	cookie, _ := r.Cookie("session_id")
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	s.DeleteCookie(db)
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
