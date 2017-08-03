// package mockmain/main starts a frontend server, application server and
// FoodTruck database server in a single machine using mock object.

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	aslibs "applicationserver/libs"
	"datacontainer/mockdatacontainer"
	ftdb "foodtruckdb/mysql"
	dslibs "foodtruckdbserver/libs"
	"frontendserver/handler"
	udb "userdb/mysql"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	_ "github.com/go-sql-driver/mysql"
)

var (
	SleepDuration       = time.Minute
	DialTimeoutDuration = 10 * time.Minute
)

func mustGetenv(k string) string {
	val := os.Getenv(k)
	if len(val) == 0 {
		log.Panicf("%s environment variable not set.", k)
	}
	return val
}

func main() {
	fmt.Println("Frontend server initializing ...")

	ctx := context.Background()

	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	dbuser := mustGetenv("CLOUDSQL_USER")
	password := mustGetenv("CLOUDSQL_PASSWORD")
	mapAPIKey := mustGetenv("MAP_API_KEY")
	dbInstanceAddress := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/datastore", dbuser, password, connectionName)

	userDB := udb.New(dbInstanceAddress)
	ftDB := ftdb.New(dbInstanceAddress)
	container := mockdatacontainer.New()

	go dslibs.New(ftDB, userDB, container, time.Second).Start(7777)
	time.Sleep(SleepDuration)

	fmt.Println("Connecting with Food truck database server...")
	conn, err := net.DialTimeout("tcp", "localhost:7777", DialTimeoutDuration)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	go aslibs.New(userDB, rpc.NewClient(conn)).Start(1234)
	time.Sleep(SleepDuration)

	fmt.Println("Connecting with application server...")
	conn, err = net.DialTimeout("tcp", "localhost:1234", DialTimeoutDuration)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	mapCl, err := maps.NewClient(maps.WithAPIKey(mapAPIKey))
	if err != nil {
		log.Fatal("fatal error: ", err)
	}

	handler.InitHandlers(ctx, rpc.NewClient(conn), mapCl)

	fmt.Println("Frontend server successfully started ...")
	http.ListenAndServe(":8080", nil)
}
