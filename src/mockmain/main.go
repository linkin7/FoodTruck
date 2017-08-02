// package mockmain/main starts a frontend server, application server and 
// FoodTruck database server in a single machine using mock object.

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

var (
    SleepDuration = time.Minute
    DialTimeoutDuration = 10 * time.Minute
)

func main() {
    fmt.Println("Frontend server initializing ...")

    userdb := mockuserdb.New(1000)
    ftdb := mockfoodtruckdb.New(1000)
    container := mockdatacontainer.New(1000)

    go dslibs.New(ftdb, userdb, container, time.Second).Start(7777)
    time.Sleep(SleepDuration)

    fmt.Println("Connecting with Food truck database server...")
    conn, err := net.DialTimeout("tcp", "localhost:7777", DialTimeoutDuration)
    if err != nil {
        log.Fatal("dialing:", err)
    }    

    go aslibs.New(userdb, rpc.NewClient(conn)).Start(1234)
    time.Sleep(SleepDuration)

	fmt.Println("Connecting with application server...")
    conn, err = net.DialTimeout("tcp", "localhost:1234", DialTimeoutDuration)
    if err != nil {
        log.Fatal("dialing:", err)
    }

    handler.InitHandlers(rpc.NewClient(conn))

    fmt.Println("Frontend server successfully started ...")
    http.ListenAndServe(":8080", nil)
}