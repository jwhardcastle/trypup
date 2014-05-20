package app

import (
	"appengine"
	"appengine/datastore"
	"net/http"
)

func dummyData(r *http.Request) {
	c := appengine.NewContext(r)

	// Start by wiping out the data we're about to populate
	// Delete all users
	dq := datastore.NewQuery("User").KeysOnly()
	d, err := dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	// Delete all items
	dq = datastore.NewQuery("Item").KeysOnly()
	d, err = dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	// Delete all comments
	dq = datastore.NewQuery("Comment").KeysOnly()
	d, err = dq.GetAll(c, nil)
	datastore.DeleteMulti(c, d)

	u1 := User{Username: "jwhardcastle", userKey: datastore.NewKey(c, "User", "jwhardcastle", 0, nil)}
	u2 := User{Username: "jhutton", userKey: datastore.NewKey(c, "User", "jhutton", 0, nil)}
	u3 := User{Username: "rkavalsky", userKey: datastore.NewKey(c, "User", "rkavalsky", 0, nil)}
	u4 := User{Username: "teej", userKey: datastore.NewKey(c, "User", "teej", 0, nil)}

	i1 := Item{
		Title:        "Baltimore Museum of Industry, learn how a linotype works, among the city's industrial history",
		owner:        &u1,
		URLTitle:     "baltimore-museum-of-industry-learn-how-a",
		Score:        36,
		Upvotes:      40,
		Downvotes:    4,
		Lat:          39.273556,
		Long:         -76.601806,
		CommentCount: 2,
		Color:        "cadetblue",
		Icon:         "truck",
	}

	i2 := Item{
		Title:        "OPACY: Oriole Park at Camden Yards, Home of the Baltimore Orioles",
		owner:        &u2,
		URLTitle:     "opacy-oriole-park-at-camden-yards-home-o",
		Score:        129,
		Upvotes:      150,
		Downvotes:    11,
		Lat:          39.283501,
		Long:         -76.6219798,
		CommentCount: 1,
		Color:        "orange",
		Icon:         "sun-o",
	}

	c1 := Comment{
		owner:     &u3,
		Body:      "We love going here!",
		Score:     3,
		Upvotes:   3,
		Downvotes: 0,
	}

	c2 := Comment{
		owner:     &u4,
		Body:      "typography geek heaven",
		Score:     5,
		Upvotes:   5,
		Downvotes: 0,
	}
	c3 := Comment{
		owner:         &u1,
		Body:          "Agreed! Among other things.",
		parentComment: &c2,
		Score:         0,
		Upvotes:       1,
		Downvotes:     1,
	}

	//c2.Children = append(c2.Children, &c3)

	// We are generating comments as if we had these associations
	//i1.comments = []*Comment{&c2}
	//i2.comments = []*Comment{&c1}

	userKeys, err := datastore.PutMulti(c, []*datastore.Key{u1.userKey, u2.userKey, u3.userKey, u4.userKey}, []interface{}{&u1, &u2, &u3, &u4})
	check(err, "Could not store users in datastore.")

	itemKeys := []*datastore.Key{
		datastore.NewIncompleteKey(c, "Item", nil),
		datastore.NewIncompleteKey(c, "Item", nil),
	}

	i1.OwnerKey = userKeys[0]
	i2.OwnerKey = userKeys[1]

	itemKeys, err = datastore.PutMulti(c, itemKeys, []interface{}{&i1, &i2})
	check(err, "Could not store items in datastore.")

	commentKeys := []*datastore.Key{
		datastore.NewIncompleteKey(c, "Comment", nil),
		datastore.NewIncompleteKey(c, "Comment", nil),
	}

	c1.ParentKey = itemKeys[1]
	c1.OwnerKey = userKeys[0]

	c2.ParentKey = itemKeys[0]
	c2.OwnerKey = userKeys[3]

	commentKeys, err = datastore.PutMulti(c, commentKeys, []interface{}{&c1, &c2})
	check(err, "Could not store comments in datastore.")

	c3.ParentKey = commentKeys[1]
	c3.OwnerKey = userKeys[0]
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Comment", nil), &c3) // The third comment is a child on the second
	check(err, "Could not store comments in datastore.")

	/*q := datastore.NewQuery("User").Order("Id")

	var users []User
	q.GetAll(c, &users)

	var items []Item
	q = datastore.NewQuery("Item").Order("-Score")
	q.GetAll(c, &items)*/

	//return items, users
}
