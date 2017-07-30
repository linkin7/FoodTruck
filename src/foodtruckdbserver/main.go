package main 

import (
	"fmt"
	"flag"
	"log"

	"net"
	"net/http"
	"net/rpc"

	"applicationserver/libs"
	"userdb/mockuserdb"
)

var port = flag.Int("port", 1234, "Port number to start server")

func main() {
	fmt.Println("Application server initializing ...")

	fmt.Println("Application server connecting with database server...")
    client, err := rpc.DialHTTP("tcp", "localhost:7777")
	if err != nil {
		log.Fatal("Dialing database server:", err)
	}

	srv := libs.New(mockuserdb.New(1000), client)
	rpc.Register(srv)
	rpc.HandleHTTP()
	
	fmt.Println("Application server opening tcp port ...")
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Application server successfully started ...")
	http.Serve(l, nil)
}
