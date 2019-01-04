package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/go-park-mail-ru/2018_2_codeloft/auth"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

const uploadPath = "/var/www/avatars/"

type UserAvatar struct {
	Db *sql.DB
	Sm auth.AuthCheckerClient
}

func userUpdateAvatar(w http.ResponseWriter, r *http.Request, db *sql.DB, sm auth.AuthCheckerClient) {
	//var session *models.Session
	//if session = services.GetCookie(r, db); session == nil {
	//	w.WriteHeader(http.StatusUnauthorized)
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
		fmt.Println("[ERROR] checkAuth in userUpdateLang:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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

	file, handle, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	name := uuid.NewV4().String()
	switch mimeType {
	case "image/jpeg":
		name = name + ".jpeg"
		saveFile(w, file, name)
	case "image/png":
		name = name + ".png"
		saveFile(w, file, name)
	default:
		jsonResponse(w, http.StatusBadRequest, "The format file is not valid.")
		return
	}
	user.Avatar = name
	err = user.UpdateAvatar(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		zap.S().Infow("Error in lang update", "err", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserAvatar) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodPost:
		userUpdateAvatar(w, r, h.Db, h.Sm)
	}
}

func saveFile(w http.ResponseWriter, file multipart.File, filename string) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	err = ioutil.WriteFile(uploadPath+filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
