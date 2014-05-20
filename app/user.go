// Users log in to vote, share, and leave comments

package app

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type User struct {
	Username     string
	PasswordHash string
	DateCreated  time.Time
	userKey      *datastore.Key `datastore:"-"`
}

func (fn User) string(u User) string {
	return u.Username
}

func getUser(c appengine.Context, username string) User {
	k := datastore.NewKey(c, "User", username, 0, nil)
	user := new(User)

	err := datastore.Get(c, k, user)
	check(err, "Could not load user profile.")

	return *user
}