// Integration test of FTServer using mock DataContainer and FoodTruckDbManager.
package libs

import (
	"reflect"
	"testing"
	"time"

	"common"
	"datacontainer/mockdatacontainer"
	"foodtruckdb/mockfoodtruckdb"
	"userdb/mockuserdb"
)

func initServer(t *testing.T) *FTServer {
	userdb := mockuserdb.New(1000)
	// User #1
	userdb.AddUser("a", "b", "mexican")
	// User #2
	userdb.AddUser("m", "n", "italian")
	// User #3
	userdb.AddUser("x", "v", "french")

	ftdb := mockfoodtruckdb.New(1000)
	ftdb.UpdateFoodTruck(1, 2, 3, 0)
	ftdb.UpdateFoodTruck(2, 4, 3, 0)
	ftdb.UpdateFoodTruck(3, 6, 4, 0)

	container := mockdatacontainer.New(1000)

	return New(ftdb, userdb, container, time.Nanosecond)
}

func TestFindNearestFoodTruck(t *testing.T) {
	srv := initServer(t)
	var b bool

	srv.CloseFoodTruck(&common.TruckData{
		UID: 1,
	}, &b)
	time.Sleep(time.Second)

	// 2 nearest FoodTruclk of location co-ordinate (0, 0)
	loc = &common.Location{
		Lat:     0,
		Lon:     0,
		Payload: 2,
	}
	got := &[]*common.TruckData{}

	srv.FindNearestFoodTruck(loc, got)
	want := &[]*common.TruckData{
		{2, 4, 3, "italian"},
		{3, 6, 4, "french"},
	}
	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("FindNearestFoodTruck #1: Got %v, want %v", *got, *want)
	}

	srv.UpdateFoodTruck(&common.TruckData{
		UID: 1,
		Lat: 1,
		Lon: 1,
	}, &b)
	time.Sleep(time.Second)

	// 2 nearest FoodTruck of location co-ordinate (0, 0)
	loc = &common.Location{
		Lat:     0,
		Lon:     0,
		Payload: 2,
	}
	got = &[]*common.TruckData{}

	srv.FindNearestFoodTruck(loc, got)
	want = &[]*common.TruckData{
		{1, 1, 1, "mexican"},
		{2, 4, 3, "italian"},
	}
	if !reflect.DeepEqual(*got, *want) {
		t.Errorf("FindNearestFoodTruck #2: Got %v, want %v", *got, *want)
	}
}
