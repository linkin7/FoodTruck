package handler

import (
	"errors"
	"fmt"
	"net/http"

	"common"
)

type Page struct {
	Title string
	Body  string
}

func loadPage(title, body string) (*Page, error) {
	return &Page{Title: title, Body: string(body)}, nil
}

func isLoggedIn(r *http.Request) bool {
	return len(currentUserName(r)) > 0
}

func currentUserName(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil || len(cookie.Value) == 0 {
		return ""
	}
	return cookie.Value
}

func currentUserID(r *http.Request) (int64, error) {
	name := currentUserName(r)
	if len(name) == 0 {
		return -1, errors.New("Username is empty")
	}
	var reply int64
	err := appSrvClient.Call("AppServer.UserID", &common.UserData{Name: name}, &reply)
	if err != nil {
		fmt.Printf("Fetching userID error: %v", err)
		return -1, err
	}
	if reply == -1 {
		return -1, fmt.Errorf("UserID for username %v can't be found", name)
	}
	return reply, nil
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
