// package common declares all the interface and data types are being used
// in the codebase.

package common

// FoodTruckManager defines an interface to wrap a nosql database server
type FoodTruckDbManager interface {
	UpdateFoodTruck(oID int64, lat float64, lon float64, cl int) bool
	CloseFoodTruck(oID int64) bool
	FindFoodTruck(oID int64) *Location
	FoodTruckCluster(oID int64) int

	// Clusterdata returnes the all location data of a particular cluster.
	ClusterData(cl int) []*Location
}

// UserDbManager defines an interface to wrap a sql database.
type UserDbManager interface {
	AddUser(name, pw, cuisine string) int64
	ValidateUser(name, pw string) bool
	UserID(name string) int64
	CuisineType(uID int64) string
}

// DataContainer defines an interface which should provide a datastructure for
// querying co-ordinate based data. Ideally it should implement a quadtree.
type DataContainer interface {
	Insert(loc *Location)
	Remove(loc *Location)
	KNearestNeighbour(loc *Location, k int) []*Location 

	// Distance should implement a distance function of two locations based on 
	// usecase
	Distance(loc1 *Location, loc2 *Location) float64

	// Serialize serializes the whole structure in wire format, so that it can be 
	// store in any database. 
	Serialize() string

	// Deserialize builds the whole structure from the given serialized data.
	Deserialize(s string)

	// Generate builds the container from the given location data.
	Generate(locs []*Location)
}