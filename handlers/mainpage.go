package handlers

import (
	"net/http"
	"text/template"
)

var MainPage = func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, struct{}{})
}
