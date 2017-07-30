package main

import (
	"log"
    "fmt"
    "net/http"
    "net/rpc"

    "frontendserver/handler"
    "applicationserver/libs"
    "userdb/mockuserdb"
)

func main() {
	fmt.Println("Frontend server initializing ...")

    go libs.New(mockuserdb.New(1000), nil).Start(1234)

	fmt.Println("Frontend server connecting with application server...")
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing application server:", err)
	}

    handler.InitHandlers(client)

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(":8080", nil)
}