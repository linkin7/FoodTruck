package mockfoodtruckdb

import (
	"testing"
	"reflect"

	"common"
)

func TestUpdateFoodTruck(t *testing.T) {
	db := New(10)
	db.UpdateFoodTruck(1, 1.5, 2.5, 0)
	db.UpdateFoodTruck(2, 3.5, 4.5, 0)

	want := &common.Location{
		ID: 1,
		Lat: 1.5,
		Lon: 2.5,
	}
	got := db.FindFoodTruck(1);
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateFoodTruck #1: Got %v, want %v", got, want)
	}

	db.UpdateFoodTruck(1, 3, 3.5)
	want = &common.Location{
		ID: 1,
		Lat: 3,
		Lon: 3.5,
	}
	got = db.FindFoodTruck(1);
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateFoodTruck #2: Got %v, want %v", got, want)
	}
}