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
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	appServerAddress = flag.String("application_server_address", "localhost:1234", "Network address and port of application server")
	port             = flag.Int("port", 8080, "Port number to start server")
	mapAPIKey        = flag.String("map_api_key", "AIzaSyDZymgaZni09lCVqoeGtL7c87cGlMkWw-s", "Developer key to invoke Google map API.")
)

func main() {
	fmt.Println("Frontend server initializing ...")

	ctx := context.Background()

	fmt.Println("Frontend server connecting with application server...")
	conn, err := net.DialTimeout("tcp", *appServerAddress, 15*time.Minute)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	asCl := rpc.NewClient(conn)
	mapCl, err := maps.NewClient(maps.WithAPIKey(*mapAPIKey))
	if err != nil {
		log.Fatal("fatal error: ", err)
	}

	handler.InitHandlers(ctx, asCl, mapCl)

	fmt.Println("Frontend server successfully started ...")
	http.ListenAndServe(strconv.Itoa(*port), nil)
}
