package app

import (
	"math/rand"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
)

func dummyData(r *http.Request, c appengine.Context) {

	// Start by wiping out the data we're about to populate
	// Delete all users
	dq := datastore.NewQuery("User").KeysOnly()
	d, _ := dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	// Delete all items
	dq = datastore.NewQuery("Item").KeysOnly()
	d, _ = dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	// Delete all comments
	dq = datastore.NewQuery("Comment").KeysOnly()
	d, _ = dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	// Delete all comments
	dq = datastore.NewQuery("Vote").KeysOnly()
	d, _ = dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	u1 := NewUser(c, "jwhardcastle", "")
	u2 := NewUser(c, "jhutton", "")
	u3 := NewUser(c, "rkavalsky", "")
	u4 := NewUser(c, "teej", "")

	i1 := NewItem(c, "Baltimore Museum of Industry, learn how a linotype works, among the city's industrial hiCstory", "This is a really cool museum that has lots of interesting displays.", "truck", "cadetblue", 39.273556, -76.601806, u1)

	i2 := NewItem(c, "OPACY: Oriole Park at Camden Yards, Home of the Baltimore Orioles", "Camden Yards is the first of the modern \"retro\" stadiums that harkens back to an earlier age of Baseball", "sun-o", "orange", 39.283501, -76.6219798, u2)

	c1 := NewComment(c, "We love going here!", u3, i2)
	c2 := NewComment(c, "typography geek heaven", u4, i1)
	c3 := NewComment(c, "Agreed! Among other things.", u1, c2)

	votables := []Votable{i1, i2, c1, c2, c3}

	for i := 0; i < 100; i++ {
		u := NewUser(c, "user"+strconv.FormatInt(int64(i), 10), "")
		for _, v := range votables {
			var value = rand.Intn(5) - 1
			if value == 0 {
				continue
			} else if value > 0 {
				value = 1
			}
			NewVote(c, u, v, int8(value))
		}
	}

	for _, v := range votables {
		v.CountVotes(c)
	}

}
