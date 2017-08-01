package common

type UserData struct {
	UID int64
	Name string
	Pw string
	Cuisine string
}

type TruckData struct {
	UID int64
	Lat float64
	Lon float64
	Cuisine string
}

type Location struct {
	ID int64
	Lat float64
	Lon float64
	Payload int
}