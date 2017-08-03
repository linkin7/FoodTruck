package handler

import (
	"errors"
	"fmt"
	"net/http"

	"common"
	"googlemaps.github.io/maps"
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

// Map related API

func findLatLon(adr string) (float64, float64, error) {
	if len(adr) == 0 {
		return 0, 0, nil
	}
	r := &maps.GeocodingRequest{
		Address: adr,
	}
	res, err := mapClient.Geocode(ctx, r)
	if err != nil {
		return 0, 0, err
	}
	if len(res) != 1 {
		return 0, 0, fmt.Errorf("Address is not correct!")
	}

	return res[0].Geometry.Location.Lat, res[0].Geometry.Location.Lng, nil
}

func findAddress(lat, lon float64) (string, error) {
	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{lat, lon},
	}
	res, err := mapClient.ReverseGeocode(ctx, r)
	if err != nil {
		return "Unknown address", err
	}
	if len(res) == 0 {
		return "Unknown address", fmt.Errorf("Unknown address")
	}

	return res[0].FormattedAddress, nil
}
