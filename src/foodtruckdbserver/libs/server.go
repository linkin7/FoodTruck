package libs

import (
	"fmt"
	"log"
	"strconv"
	"net"
	"net/rpc"
	"sync"
	"time"

	"common"
)

type FTServer struct {
	ftDbMgr common.FoodTruckDbManager
	uDbMgr common.UserDbManager
    container common.DataContainer

    curCluster int

    lastUpdTime time.Time
    updInterval time.Duration

    mu sync.Mutex
}

func New(ftmgr common.FoodTruckDbManager, umgr common.UserDbManager, c common.DataContainer, updInterval time.Duration) *FTServer {
	srv := &FTServer{
		ftDbMgr: ftmgr,
		uDbMgr: umgr,
		container: c,
		updInterval: updInterval,
	}
	srv.updateContainer()
	return srv
}

func (srv *FTServer) Start(port int) {
	fmt.Println("Food Truck Database server starting ...")
	rpc.Register(srv)
	
	fmt.Println("Food Truck Database server opening tcp port ...")
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Food Truck Database server successfully started ...")
	rpc.Accept(l)
}

func (srv *FTServer) updateContainer() {
	if curTime := time.Now(); curTime.Sub(srv.lastUpdTime) < srv.updInterval {
		return
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	if curTime := time.Now(); curTime.Sub(srv.lastUpdTime) < srv.updInterval {
		return
	}
	locs := srv.ftDbMgr.ClusterData(srv.curCluster)
	srv.container.Generate(locs)
	srv.lastUpdTime = time.Now()
}

func (srv *FTServer) UpdateFoodTruck(td *common.TruckData, ok *bool) error {
	if err := srv.CloseFoodTruck(td, ok); err != nil {
		return err
	}
	*ok = srv.ftDbMgr.UpdateFoodTruck(td.UID, td.Lat, td.Lon, 0)
	return nil
}

func (srv *FTServer) CloseFoodTruck(td *common.TruckData, ok *bool) error {
	*ok = srv.ftDbMgr.CloseFoodTruck(td.UID)
	return nil
}

func (srv *FTServer) FindNearestFoodTruck(loc *common.Location, list *[]*common.TruckData) error {
	srv.updateContainer()
	locs := srv.container.KNearestNeighbour(loc, loc.Payload)
	*list = srv.processFoodTruck(locs)
	return nil
}

func (srv *FTServer) processFoodTruck(list []*common.Location) []*common.TruckData {
	tds := []*common.TruckData{}
	for _, loc := range list {
		tds = append(tds, &common.TruckData{
			UID: loc.ID,
			Lat: loc.Lat,
			Lon: loc.Lon,
			Cuisine: srv.uDbMgr.CuisineType(loc.ID),
			})
	}
	return tds
}