// Users log in to vote, share, and leave comments

package app

import (
	"appengine"
	"appengine/datastore"
	"time"
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
)

// cost=15 takes 5 seconds
// cost=14 takes 2.5 seconds
const PASSWORDCOST int =  14

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
func New(c appengine.Context, username string, password string) *User {
	user := new(User)
	user.Username = username
	user.setPassword(password)
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

	if(len(username) == 0) {
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
	if (err != nil) {
		return err
	}
	
	(*user).PasswordHash = ph
	
	return nil
}

func (user User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}

