package main 

import (
	"flag"
	"time"

	"userdb/mockuserdb"
    "foodtruckdb/mockfoodtruckdb"
    "datacontainer/mockdatacontainer"
    "foodtruckdbserver/libs"
)

var (
	port = flag.Int("port", 1234, "Port number to start server")
	updInterval = flag.Duration("update_iterval", time.Hour, "Minimum duration for updating in-memory data of server.")
)

func main() {
	userdb := mockuserdb.New(1000)
    ftdb := mockfoodtruckdb.New(1000)
    container := mockdatacontainer.New(1000)

    libs.New(ftdb, userdb, container, *updInterval).Start(*port)
}