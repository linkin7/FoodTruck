package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"

    "frontendserver/handler"
)

var (
    appServerAddress = flag.String("application_server_address", "localhost:1234", "Network address and port of application server")
    port = flag.Int("port", 8080, "Port number to start server")
)

func main() {
	fmt.Println("Frontend server initializing ...")

	fmt.Println("Frontend server connecting with application server...")
    client, err := rpc.DialHTTP("tcp", *appServerAddress)
	if err != nil {
		log.Fatal("Dialing application server:", err)
	}

    handler.InitHandlers(client)

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(strconv.Itoa(*port), nil)
}