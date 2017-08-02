package mockuserdb

import "testing"

func TestValidateUser(t *testing.T) {
	db := New()
	db.AddUser("user1", "pw1", "italian")
	db.AddUser("user2", "pw2", "swiss")

	var tests = []struct {
		name string
		pw string
		want bool
	} {
		{"user2", "pw2", true},
		{"user1", "pw2", false},
		{"user4", "pw2", false},
		{"user1", "pw6", false},
	}

	for i, ts := range tests {
		got := db.ValidateUser(ts.name, ts.pw)
		if got != ts.want {
			t.Errorf("Validate User %v: Got %v, want %v", i, got, ts.want)
		}
	}
}