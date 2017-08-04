// Package handler implements all the URL handle processor for the frontend server.
// Before using this package InitHandlers() function should be called with a client
// of Application server.
package handler

import (
	"fmt"
	"net/http"
	"net/rpc"

	"common"
	"frontendserver/htmltemplate"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	ctx context.Context

	// appSrvClient holds the TCP connection with Appllication server.
	appSrvClient *rpc.Client

	// mapClient holds client connection with Google map API server.
	mapClient *maps.Client
)

// InitHandles registers all the handler function with corresponding url path.
func InitHandlers(c context.Context, asCl *rpc.Client, mapCl *maps.Client) {
	ctx = c
	appSrvClient = asCl
	mapClient = mapCl

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

	lat, lon, err := findLatLon(r.FormValue("address"))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	uID, err := currentUserID(r)
	if err != nil {
		fmt.Fprintf(w, "FoodTruck update error: %v", err)
		return
	}

	td := &common.TruckData{uID, lat, lon, ""}
	method := ""
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

	lat, lon, err := findLatLon(r.FormValue("address"))
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	reply := &[]*common.TruckData{}
	err = appSrvClient.Call("AppServer.FindNearest", &common.Location{
		Lat:     lat,
		Lon:     lon,
		Payload: 10,
	}, &reply)
	if err != nil {
		log.Println("Finding Nearest Truck error: ", err)
		fmt.Fprintln(w, "Can't find nearest Truck")
		return
	}

	if len(*reply) == 0 {
		fmt.Fprintln(w, "No available truck around")
		return
	}

	fmt.Fprintln(w, "<table border='1'><tr><th>Cuisine Type</th><th>Address</th></tr>")
	for _, td := range *reply {
		adr, err := findAddress(td.Lat, td.Lon)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td></tr>", td.Cuisine, adr)
	}
	fmt.Fprintln(w, "</table>")
}
