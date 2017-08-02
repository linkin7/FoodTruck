// package mockuserdb implements a in-memory mock version of FoodTruckDbManager
// interface (common/data_manager_interface.go). Currently it pushes all the data in
// default 0 cluster.

package mysql

import (
    "database/sql"
    "fmt"
    "log"

    "common"
)

type FoodTrucks struct {
    db *sql.DB
}

func New(dataSourceName string) *FoodTrucks {
    log.Printf("Opening mysql food truck database: %s\n", dataSourceName)
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
    	log.Fatalf("Could not open db: %v", err)
    	return nil
    }
    return &FoodTrucks{
    	db: db,
    }
}

func (fts *FoodTrucks) UpdateFoodTruck(oID int64, lat float64, lon float64, cl int) bool {
	fts.CloseFoodTruck(oID)
	q := fmt.Sprintf("INSERT INTO foodtrucks (oID, latitude, longitude, cluster) VALUES (%v, %v, %v, %v)", oID, lat, lon, cl)
	rows, err := fts.db.Query(q)
    if err != nil {
        log.Printf("%v", err)
        return false
    }
    defer rows.Close()
    return false
}

func (fts *FoodTrucks) CloseFoodTruck(oID int64) bool {
	q := fmt.Sprintf("DELETE FROM foodtrucks WHERE oID = %v", oID)
	rows, err := fts.db.Query(q)
    if err != nil {
        log.Printf("%v", err)
        return false
    }
    defer rows.Close()
    return false
}

func (fts *FoodTrucks) FindFoodTruck(oID int64) *common.Location {
	q := fmt.Sprintf("SELECT oID, latitude, longitude FROM foodtrucks WHERE oID = %v", oID)
	rows, err := fts.db.Query(q)
    if err != nil {
        log.Printf("%v", err)
        return nil
    }
    defer rows.Close()

    loc := common.Location{}
    for rows.Next() {
        if err := rows.Scan(&loc.ID, &loc.Lat, &loc.Lon); err != nil {
            log.Printf("%v", err)
            return nil
        }
        return &loc
    }
    return &loc
}

func (fts *FoodTrucks) ClusterData(cl int) []*common.Location {
	q := fmt.Sprintf("SELECT oID, latitude, longitude FROM foodtrucks WHERE cluster = %v", cl)
	rows, err := fts.db.Query(q)
    if err != nil {
        log.Printf("%v", err)
        return nil
    }
    defer rows.Close()

    locs := []*common.Location{}
    for rows.Next() {
    	loc := common.Location{}
        if err := rows.Scan(&loc.ID, &loc.Lat, &loc.Lon); err != nil {
            log.Printf("%v", err)
            return nil
        }
        locs = append(locs, &loc)
    }
    return locs
}

func (fts *FoodTrucks) FoodTruckCluster(oID int64) int {
	q := fmt.Sprintf("SELECT cluster FROM foodtrucks WHERE oID = %v", oID)
	rows, err := fts.db.Query(q)
    if err != nil {
        log.Printf("%v", err)
        return 0
    }
    defer rows.Close()

    for rows.Next() {
    	var cl int
        if err := rows.Scan(&cl); err != nil {
            log.Printf("%v", err)
            return 0
        }
        return cl
    }
    return 0
}