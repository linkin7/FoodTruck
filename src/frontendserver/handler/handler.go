// Package handler implements all the URL handle processor for the frontend server.
// Before using this package InitHandlers() function should be called with a client
// of Application server.
package handler

import (
	"fmt"
	"net/http"
	"net/rpc"
	"strconv"

	"common"
	"frontendserver/htmltemplate"
)

// appSrvClient holds the TCP connection with Appllication server.
var appSrvClient *rpc.Client

// InitHandles registers all the handler function with corresponding url path.
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
	http.HandleFunc("/updateconfirm", updateConfirmHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/findnearest", findNearestHandler)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	p, _ := loadPage("register", htmltemplate.Register)
	fmt.Fprintf(w, "%v", p.Body)
}

func postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	ud := &common.UserData{
		Name:    r.FormValue("name"),
		Pw:      r.FormValue("password"),
		Cuisine: r.FormValue("cuisine"),
	}

	var reply int
	err := appSrvClient.Call("AppServer.Register", ud, &reply)
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
	if !isLoggedIn(r) {
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
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	method := ""
	// TODO: check error after parsing float.
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	uID, err := currentUserID(r)
	if err != nil {
		fmt.Fprintf(w, "FoodTruck update error: %v", err)
		return
	}
	td := &common.TruckData{uID, lat, lon, ""}

	if len(r.FormValue("start")) > 0 {
		method = "AppServer.UpdateFoodTruck"
	} else {
		method = "AppServer.CloseFoodTruck"
	}

	var reply bool
	err = appSrvClient.Call(method, td, &reply)
	if err != nil {
		fmt.Fprintf(w, "FoodTruck update error: %v", err)
		return
	}
	if reply == false {
		fmt.Fprintln(w, "FoodTruck update failed! Try again!")
		return
	}

	http.Redirect(w, r, "/updateconfirm", http.StatusSeeOther)
}

func updateConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Fprintln(w, "<b>Food Truck status successfully updated</b>")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
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
	if isLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	p, _ := loadPage("login", htmltemplate.Login)
	fmt.Fprintf(w, "%v", p.Body)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	ud := &common.UserData{
		Name: r.FormValue("name"),
		Pw:   r.FormValue("password"),
	}
	var reply bool
	err := appSrvClient.Call("AppServer.Login", ud, &reply)
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

func findNearestHandler(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// TODO: check error after parsing float.
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)

	reply := &[]*common.TruckData{}
	err := appSrvClient.Call("AppServer.FindNearest", &common.Location{
		Lat:     lat,
		Lon:     lon,
		Payload: 3,
	}, &reply)
	if err != nil {
		fmt.Fprintf(w, "Finding Nearest Truck error: %v", err)
		return
	}

	if len(*reply) == 0 {
		fmt.Fprintln(w, "No available truck around")
		return
	}

	fmt.Fprintln(w, "<table border='1'><tr><th>Cuisine Type</th><th>Latitude</th><th>Longitude</th></tr>")
	for _, td := range *reply {
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td></tr>", td.Cuisine, td.Lat, td.Lon)
	}
	fmt.Fprintln(w, "</table>")
}
