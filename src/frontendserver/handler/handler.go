package handler 

import (
	"fmt"
	"net/http"
	"net/rpc"

	"common"
	"frontendserver/htmltemplate"
)

var appSrvClient *rpc.Client

func InitHandlers(client *rpc.Client) {
	appSrvClient = client

	http.HandleFunc("/", handler)
    http.HandleFunc("/index", handler)
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/postregister", postRegisterHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/postlogin", postLoginHandler)
    http.HandleFunc("/home", homePageHandler)
    http.HandleFunc("/update", updateHandler)
    http.HandleFunc("/logout", logoutHandler)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) > 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	p, _ := loadPage("register", htmltemplate.Register)
	fmt.Fprintf(w, "%v", p.Body)
}


func postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) > 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	ud := &common.UserData{
		Name: r.FormValue("name"),
		Pw: r.FormValue("password"),
		Cuisine: r.FormValue("cuisine"),
	}
	var reply int
	err := appSrvClient.Call("Server.Register", ud, &reply)
	if err != nil {
		fmt.Fprintf(w, "Registration error: %v", err)
		return
	}
	if reply == -1 {
		fmt.Fprintf(w, "Name %v already exists!", r.FormValue("name"))
		return
	}

	makeLoggedIn(r.FormValue("name"), w, r)
    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p, err := loadPage("home", htmltemplate.Home)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
        return
    }
	fmt.Fprintf(w, "%v", p.Body)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	fmt.Println("%v", r.Form)
	// need tocll rpc
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) > 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	p, err := loadPage("index", htmltemplate.Index)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
        return
    }
	fmt.Fprintf(w, "%v", p.Body)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) > 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	p, _ := loadPage("login", htmltemplate.Login)
	fmt.Fprintf(w, "%v", p.Body)
}


func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	if len(isLoggedIn(r)) > 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	ud := &common.UserData{
		Name: r.FormValue("name"),
		Pw: r.FormValue("password"),
	}
	var reply bool
	err := appSrvClient.Call("Server.Login", ud, &reply)
	if err != nil {
		fmt.Fprintf(w, "Login error: %v", err)
		return
	}
	if reply == false {
		fmt.Fprintln(w, "Name or password doesn't match!")
		return
	}

	makeLoggedIn(r.FormValue("name"), w, r)
    http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	makeLoggedOut(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}