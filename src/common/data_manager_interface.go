package common

type FoodTruckDbManager interface {
	UpdateFoodTruck(oID int64, lat float64, lon float64) bool
	CloseFoodTruck(oID int64) bool
	FindFoodTruck(oID int64) *Location
	FoodTruckCluster(oID int64) int
	ClusterData(cl int) []*Location
}

type UserDbManager interface {
	AddUser(name, pw, cuisine string) int64
	ValidateUser(name, pw string) bool
	UserID(name string) int64
}

type DataContainer interface {
	Insert(loc *Location)
	Remove(loc *Location)
	KNearestNeighbour(loc *Location, k int) []*Location 
	Distance(loc1 *Location, loc2 *Location) float64
	Serialize() string
	Deserialize(s string)
	Generate(locs []*Location)
}