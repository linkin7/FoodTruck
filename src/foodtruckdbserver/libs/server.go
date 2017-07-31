package libs

import (
	"fmt"
	"log"
	"strconv"
	"net"
	"net/rpc"

	"common"
)

type FTServer struct {
	ftDbMgr common.FoodTruckDbManager
    container common.DataContainer
    cur_cluster int
}

func New(mgr common.FoodTruckDbManager, c common.DataContainer) *FTServer {
	return &FTServer{
		ftDbMgr: mgr,
		container: c,
	}
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
	locs := srv.ftDbMgr.ClusterData(srv.cur_cluster)
	srv.container.Generate(locs)
}

func (srv *FTServer) UpdateFoodTruck(td *common.TruckData, ok *bool) error {
	if err := srv.CloseFoodTruck(td, ok); err != nil {
		return err
	}
	*ok = srv.ftDbMgr.UpdateFoodTruck(td.UID, td.Lat, td.Lon)
	return nil
}

func (srv *FTServer) CloseFoodTruck(td *common.TruckData, ok *bool) error {
	*ok = srv.ftDbMgr.CloseFoodTruck(td.UID)
	return nil
}

func (srv *FTServer) FindNearestFoodTruck(loc *common.Location, list *[]*common.Location) error {
	*list = srv.container.KNearestNeighbour(loc, loc.Payload)
	return nil
}