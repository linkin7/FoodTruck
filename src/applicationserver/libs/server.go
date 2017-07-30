package libs

import (
	//"errors"
	"net/rpc"

	"common"
)

type Server struct {
	uDbMgr common.UserDbManager
    ftDbClient *rpc.Client
}

func New(uDbM common.UserDbManager, client *rpc.Client) *Server {
	return &Server{
		uDbMgr: uDbM,
		ftDbClient: client,
	}
}

func (srv *Server) Register(ud *common.UserData, ID *int64) error {
	*ID = srv.uDbMgr.AddUser(ud.Name, ud.Pw, ud.Cuisine)
	return nil
}

func (srv *Server) Login(ud *common.UserData, ok *bool) error {
	*ok = srv.uDbMgr.ValidateUser(ud.Name, ud.Pw)
	return nil
}

func (srv *Server) UserID(ud *common.UserData, ID *int64) error {
	*ID = srv.uDbMgr.UserID(ud.Name)
	return nil
}

func (srv *Server) UpdateFoodTruck(td *common.TruckData, ok *bool) error {
	err := srv.ftDbClient.Call("FTServer.UpdateFoodTruck", td, ok)
	if err != nil {
		return err
	}
	return nil
}

func (srv *Server) CloseFoodTruck(td *common.TruckData, ok *bool) error {
	err := srv.ftDbClient.Call("FTServer.CloseFoodTruck", td, ok)
	if err != nil {
		return err
	}
	return nil
}