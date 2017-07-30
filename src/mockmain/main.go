package main

import (
    "log"
    "fmt"
    "net"
    "net/http"
    "net/rpc"
    "time"

    "frontendserver/handler"
    "applicationserver/libs"
    "userdb/mockuserdb"
)

func main() {
    fmt.Println("Frontend server initializing ...")

    go libs.New(mockuserdb.New(1000), nil).Start(1234)
    time.Sleep(time.Minute)

	fmt.Println("Frontend server connecting with application server...")
    conn, err := net.DialTimeout("tcp", "localhost:1234", 10 * time.Minute)
      if err != nil {
        log.Fatal("dialing:", err)
      }

    client := rpc.NewClient(conn)

    handler.InitHandlers(client)

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(":8080", nil)
}