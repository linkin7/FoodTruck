package mockfoodtruckdb

import "common"

type foodTruck struct {
	oID int64
	lat float64
	lon float64
	cluster int
}

type Collections struct {
	capacity int64
	fts []foodTruck
}

func New(cap int64) *Collections {
	return &Collections{
		capacity: cap,
	}
}


func (c *Collections) UpdateFoodTruck(oID int64, lat float64, lon float64) bool {
	c.CloseFoodTruck(oID)
	c.fts = append(c.fts, foodTruck{
		oID: oID,
		lat: lat,
		lon: lon,
		cluster: 0,
		})
	return true
}

func (c *Collections) CloseFoodTruck(oID int64) bool {
	for i, ft := range c.fts {
		if (ft.oID == oID) {
			c.fts = append(c.fts[:i], c.fts[i+1:]...)
			return true
		}
	}
	return true
}

func (c *Collections) FindFoodTruck(oID int64) *common.Location {
	for _, ft := range c.fts {
		if (ft.oID == oID) {
			return &common.Location{
				ID: oID,
				Lat: ft.lat,
				Lon: ft.lon, 
			}
		}
	}
	return nil
}


