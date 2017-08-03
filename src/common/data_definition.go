// package common declares all the interface and data types are being used
// in the codebase.

package common

// UserData contains the registered user related data to pass in a rpc call.
type UserData struct {
	UID     int64
	Name    string
	Pw      string
	Cuisine string
}

// TruckData contains the locaton and the info of a registered user to pass
// in a rpc call. UID field denotes the user id of truck owner.
type TruckData struct {
	UID     int64
	Lat     float64
	Lon     float64
	Cuisine string
}

// TruckData contains the coordinates of a particular locaton, Payload field
// denotes different value based on context. For example for nearestNeighbour
// it contains the information about number of nearest locations should be returned.
type Location struct {
	ID      int64
	Lat     float64
	Lon     float64
	Payload int
}
