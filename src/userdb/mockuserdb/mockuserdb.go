package mockuserdb

type user struct {
	id int64
	name string
	password string
	cuisine string
}

type Collections struct {
	capacity int64
	users []user
}

func New(cap int64) *Collections {
	return &Collections{
		capacity: cap,
	}
}


func (c *Collections) AddUser(name, pw, cuisine string) int64 {
	if c.UserID(name) != -1 {
		return -1
	}
	c.users = append(c.users, user{
		id: 1,
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
