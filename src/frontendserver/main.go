package main

import (
	"log"
    "fmt"
    "net/http"
    "net/rpc"
)

var appSrvClient *rpc.Client


func main() {
	fmt.Println("Frontend server initializing ...")

	fmt.Println("Frontend server connecting with application server...")
    var err error
    appSrvClient, err = rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing application server:", err)
	}

    http.HandleFunc("/", handler)
    http.HandleFunc("/index", handler)
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/postregister", postRegisterHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/postlogin", postLoginHandler)
    http.HandleFunc("/home", homePageHandler)
    http.HandleFunc("/update", updateHandler)
    http.HandleFunc("/logout", logoutHandler)

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(":8080", nil)
}