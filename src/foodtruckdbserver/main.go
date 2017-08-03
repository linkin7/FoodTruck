package main

import (
	"flag"
	"time"

	"datacontainer/mockdatacontainer"
	"foodtruckdb/mockfoodtruckdb"
	"foodtruckdbserver/libs"
	"userdb/mockuserdb"
)

var (
	port        = flag.Int("port", 1234, "Port number to start server")
	updInterval = flag.Duration("update_iterval", time.Hour, "Minimum duration for updating in-memory data of server.")
)

func main() {
	userdb := mockuserdb.New()
	ftdb := mockfoodtruckdb.New()
	container := mockdatacontainer.New()

	libs.New(ftdb, userdb, container, *updInterval).Start(*port)
}
