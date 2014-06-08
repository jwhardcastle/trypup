// Users log in to vote, share, and leave comments

package app

import (
	"errors"
	"time"

	"code.google.com/p/go.crypto/bcrypt"

	"appengine"
	"appengine/datastore"
)

// cost=15 takes 5 seconds
// cost=14 takes 2.5 seconds
const PASSWORDCOST int = 14

type User struct {
	Username     string
	PasswordHash []byte
	DateCreated  time.Time
	userKey      *datastore.Key `datastore:"-"`
}

func (fn User) string(u User) string {
	return u.Username
}

// Create a new user
func NewUser(c appengine.Context, username string, password string) *User {
	user := new(User)
	user.Username = username
	// to speed things up while creating dummy accounts, if the password is blank, don't set it
	// this means these accounts can't login
	if len(password) > 0 {
		user.setPassword(password)
	}
	user.userKey = datastore.NewKey(c, "User", username, 0, nil)
	user.Save(c)
	return user
}

// Store the user in the datastore
func (user *User) Save(c appengine.Context) error {
	var err error
	(*user).userKey, err = datastore.Put(c, (*user).userKey, user)
	return err
}

// Return the specified user from the datastore
func getUser(c appengine.Context, username string) (User, error) {
	user := new(User)

	if len(username) == 0 {
		return *user, errors.New("Username can't be blank.")
	}

	k := datastore.NewKey(c, "User", username, 0, nil)

	err := datastore.Get(c, k, user)
	//check(err, "Could not load user profile for " + username + ".")

	return *user, err
}

// Set the user's password
func (user *User) setPassword(newPassword string) error {
	var ph []byte
	ph, err := bcrypt.GenerateFromPassword([]byte(newPassword), PASSWORDCOST)
	if err != nil {
		return err
	}

	(*user).PasswordHash = ph

	return nil
}

func (user User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}
