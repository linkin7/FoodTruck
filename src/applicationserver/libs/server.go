package libs

import (
	"fmt"
	"log"
	"strconv"
	"net"
	"net/rpc"

	"common"
)

type AppServer struct {
	uDbMgr common.UserDbManager
    ftDbClient *rpc.Client
}

func New(uDbM common.UserDbManager, client *rpc.Client) *AppServer {
	return &AppServer{
		uDbMgr: uDbM,
		ftDbClient: client,
	}
}

func (srv *AppServer) Start(port int) {
	fmt.Println("Application server starting ...")
	rpc.Register(srv)
	
	fmt.Println("Application server opening tcp port ...")
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("listen error:", err)
	}
	fmt.Println("Application server successfully started ...")
	rpc.Accept(l)
}

func (srv *AppServer) Register(ud *common.UserData, ID *int64) error {
	*ID = srv.uDbMgr.AddUser(ud.Name, ud.Pw, ud.Cuisine)
	return nil
}

func (srv *AppServer) Login(ud *common.UserData, ok *bool) error {
	*ok = srv.uDbMgr.ValidateUser(ud.Name, ud.Pw)
	return nil
}

func (srv *AppServer) UserID(ud *common.UserData, ID *int64) error {
	*ID = srv.uDbMgr.UserID(ud.Name)
	return nil
}

func (srv *AppServer) UpdateFoodTruck(td *common.TruckData, ok *bool) error {
	err := srv.ftDbClient.Call("FTServer.UpdateFoodTruck", td, ok)
	if err != nil {
		return err
	}
	return nil
}

func (srv *AppServer) CloseFoodTruck(td *common.TruckData, ok *bool) error {
	err := srv.ftDbClient.Call("FTServer.CloseFoodTruck", td, ok)
	if err != nil {
		return err
	}
	return nil
}