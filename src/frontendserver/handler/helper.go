package handler

import (
	"io/ioutil"
	"net/http"
)

type Page struct {
    Title string
    Body  string
}

func loadPage(title string) (*Page, error) {
    filename := "../src/frontendserver/html/" + title + ".html"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: string(body)}, nil
}

func isLoggedIn(r *http.Request) string {
	if cookie, err := r.Cookie("session"); err == nil {
 		return cookie.Value
 	}
	return ""
}

func makeLoggedIn(name string, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: name,
	})
}

func makeLoggedOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: "",
	})
}