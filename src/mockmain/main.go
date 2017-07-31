package main

import (
    "log"
    "fmt"
    "net"
    "net/http"
    "net/rpc"
    "time"

    "frontendserver/handler"
    "userdb/mockuserdb"
    "foodtruckdb/mockfoodtruckdb"
    "datacontainer/mockdatacontainer"
    aslibs "applicationserver/libs"
    dslibs "foodtruckdbserver/libs"
)

func main() {
    fmt.Println("Frontend server initializing ...")

    userdb := mockuserdb.New(1000)
    ftdb := mockfoodtruckdb.New(1000)
    container := mockdatacontainer.New(1000)

    go dslibs.New(ftdb, container).Start(7777)
    time.Sleep(time.Minute)

    fmt.Println("Connecting with Food truck database server...")
    conn, err := net.DialTimeout("tcp", "localhost:7777", 10 * time.Minute)
    if err != nil {
        log.Fatal("dialing:", err)
    }    

    go aslibs.New(userdb, rpc.NewClient(conn)).Start(1234)
    time.Sleep(time.Minute)

	fmt.Println("Connecting with application server...")
    conn, err = net.DialTimeout("tcp", "localhost:1234", 10 * time.Minute)
    if err != nil {
        log.Fatal("dialing:", err)
    }

    handler.InitHandlers(rpc.NewClient(conn))

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(":8080", nil)
}