// package mockuserdb implements a in-memory mock version of UserDbManager
// interface (common/data_manager_interface.go). During inserting a new
// user data it creates a unique user id.

package mockuserdb

import "sync"

type user struct {
	id int64
	name string
	password string
	cuisine string
}

type Collections struct {
	capacity int64
	users []user

	mu sync.Mutex
	id int64
}

func New(cap int64) *Collections {
	return &Collections{
		capacity: cap,
	}
}

func (c *Collections) generateID() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.id++
	return c.id
}


func (c *Collections) AddUser(name, pw, cuisine string) int64 {
	if c.UserID(name) != -1 {
		return -1
	}
	// If no cuisine is set by a user, then by default it should be None.
	if len(cuisine) == 0 {
		cuisine = "None"
	}
	c.users = append(c.users, user{
		id: c.generateID(),
		name: name,
		password: pw,
		cuisine: cuisine,
		})
	return 1
}

func (c *Collections) ValidateUser(name, pw string) bool {
	for _, u := range c.users {
		if u.name == name {
			return u.password == pw
		}
	}
	return false
}

func (c *Collections) UserID(name string) int64 {
	for _, u := range c.users {
		if u.name == name {
			return u.id
		}
	}
	return -1
}

func (c *Collections) CuisineType(id int64) string {
	for _, u := range c.users {
		if u.id == id {
			return u.cuisine
		}
	}
	return "None"
}