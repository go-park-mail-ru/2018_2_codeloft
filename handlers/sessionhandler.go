package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2018_2_codeloft/auth"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"github.com/go-park-mail-ru/2018_2_codeloft/validator"
	"github.com/mailru/easyjson"

	"go.uber.org/zap"
)

func checkAuth(w http.ResponseWriter, r *http.Request, db *sql.DB, sm auth.AuthCheckerClient) {
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
	//var s *models.Session
	//if s = services.GetCookie(r, db); s == nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

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
	if !user.GetUserByID(db, userid.UserID) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(generateError(models.MyError{URL: r.URL.Path, What: "User Does Not Exist in Users table, but exist in session", Err: fmt.Errorf("")}))
		zap.L().Info("User Does Not Exist in Users table, but exist in session",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.String("Session value", cookie.Value),
			zap.Int64("User id", userid.UserID),
		)
		return
	}
	res, err := easyjson.Marshal(&user)
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

func signIn(w http.ResponseWriter, r *http.Request, db *sql.DB, sm auth.AuthCheckerClient) {
	//var s *models.Session
	//// Если уже залогинен
	//if s = services.GetCookie(r, db); s != nil {
	//	w.WriteHeader(http.StatusConflict)
	//	return
	//}
	//cooka, err := r.Cookie("session_id")
	//if err == nil {
	//	w.WriteHeader(http.StatusConflict)
	//	log.Println("[ERROR] signIn Cookie exist.AlreadyAuth:", cooka.Value)
	//	return
	//}
	cooka, err := r.Cookie("session_id")
	if cooka != nil {
		userid, err := sm.Check(context.Background(), &auth.SessionID{ID: cooka.Value})
		if err == nil {
			fmt.Println("[ERROR] signIn: Already auth. UserID:", userid.UserID)
			w.WriteHeader(http.StatusConflict)
			return
		}
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

	u := models.HelpUser{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{URL: r.URL.Path, What: "wrong requst format", Err: err}))
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
		w.Write(generateError(models.MyError{URL: r.URL.Path, What: "bad login", Err: err}))
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
		err := models.MyError{URL: r.URL.Path, What: "User does not exist", Err: models.UserDoesNotExist(u.Login)}
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

		err := models.MyError{URL: r.URL.Path, What: "wrong password", Err: fmt.Errorf("wrong password")}
		w.Write(generateError(err))
		zap.L().Info("Wrong password",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(&err),
		)

		//w.Write(generateError(models.MyError{r.URL.Path, "wrong password", fmt.Errorf("wrong password")}))

		return
	}

	cookieVal, err := sm.Create(context.Background(), &auth.Session{UserID: dbUser.Id})
	if err != nil {
		fmt.Println("[ERROR] signIn CantCreateCookie:", cookieVal.ID, "\n", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    cookieVal.ID,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	//cookie := services.GenerateCookie(dbUser.Login)
	if os.Getenv("ENV") == "production" {
		cookie.Secure = true
	}
	//s = &models.Session{cookie.Value, dbUser.Id}
	//err = s.AddCookie(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		myErr := models.MyError{URL: r.URL.Path, What: "Cant AddCookie", Err: err}
		w.Write(generateError(myErr))
		zap.L().Info("Cant AddCookie",
			zap.String("URL", r.URL.Path),
			zap.String("Method", r.Method),
			zap.String("Origin", r.Header.Get("Origin")),
			zap.String("Remote addres", r.RemoteAddr),
			zap.Error(&myErr),
		)

		//w.Write(generateError(models.MyError{r.URL.Path, "Cant AddCookie", err}))

		return
	}
	http.SetCookie(w, cookie)
	res, err := easyjson.Marshal(&dbUser)
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

func logout(w http.ResponseWriter, r *http.Request, db *sql.DB, sm auth.AuthCheckerClient) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = sm.Delete(context.Background(), &auth.SessionID{ID: cookie.Value})
	if err != nil {
		fmt.Println("[ERROR] logOut Cant Delete cookie:", cookie.Value, "\n", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	// var s models.Session
	// if !s.CheckCookie(db, cookie.Value) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	//var s *models.Session
	//if s = services.GetCookie(r, db); s == nil {
	//	zap.L().Info("StatusConflist",
	//		zap.String("URL", r.URL.Path),
	//		zap.String("Method", r.Method),
	//		zap.String("Origin", r.Header.Get("Origin")),
	//		zap.String("Remote addres", r.RemoteAddr),
	//		zap.Int("Code", http.StatusConflict),
	//	)
	//	w.WriteHeader(http.StatusConflict)
	//	return
	//}
	//cookie, _ := r.Cookie("session_id")
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	//err := s.DeleteCookie(db)
	//if err != nil {
	//	zap.L().Warn("Can not delete cookie",
	//		zap.String("URL", r.URL.Path),
	//		zap.String("Method", r.Method),
	//		zap.String("Origin", r.Header.Get("Origin")),
	//		zap.String("Remote addres", r.RemoteAddr),
	//		zap.Error(err),
	//	)
	//}
	w.WriteHeader(http.StatusOK)
}

type SessionHandler struct {
	Db *sql.DB
	Sm auth.AuthCheckerClient
}

func (h *SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {

	case http.MethodGet:
		checkAuth(w, r, h.Db, h.Sm)
	case http.MethodPost:
		signIn(w, r, h.Db, h.Sm)
	case http.MethodDelete:
		logout(w, r, h.Db, h.Sm)
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
