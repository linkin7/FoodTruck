// package frontendserver/main starts a frontend server.
package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/rpc"
    "strconv"
    "time"

    "frontendserver/handler"
)

var (
    appServerAddress = flag.String("application_server_address", "localhost:1234", "Network address and port of application server")
    port = flag.Int("port", 8080, "Port number to start server")
)

func main() {
	fmt.Println("Frontend server initializing ...")

	fmt.Println("Frontend server connecting with application server...")
    conn, err := net.DialTimeout("tcp", *appServerAddress, 15 * time.Minute)
      if err != nil {
        log.Fatal("dialing:", err)
      }

    client := rpc.NewClient(conn)

    handler.InitHandlers(client)

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(strconv.Itoa(*port), nil)
}