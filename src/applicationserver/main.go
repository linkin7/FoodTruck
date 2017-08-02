// Package applicationserver/main starts an application server.

package main 

import (
	"fmt"
	"flag"
	"log"
	"net/rpc"

	"applicationserver/libs"
	"userdb/mockuserdb"
)

var (
	dbServerAddress = flag.String("database_server_address", "localhost:7777", "Network address and port of database server")
	port = flag.Int("port", 1234, "Port number to start server")
)

func main() {
	fmt.Println("Application server initializing ...")

	fmt.Println("Application server connecting with database server...")
    client, err := rpc.DialHTTP("tcp", *dbServerAddress)
	if err != nil {
		log.Fatal("Dialing database server:", err)
	}

	// TODO: change the mockuserdb by actual database class.
	libs.New(mockuserdb.New(1000), client).Start(*port)
}
