// A sample front page to get started
package app

import (
	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
	"net/http"
	"bytes"
)

// Load dummy data
func DummyHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	dummyData(r, c)
}

// Show the front page
func RootHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates := setup(w, r)

	var items []Item
	q := datastore.NewQuery("Item").Order("-Score")
	keys, err := q.GetAll(c, &items)
	check(err, "Could not load items.")

	for i, key := range keys {
		items[i].itemKey = key
		items[i].loadOwner(c)
	}

	var b bytes.Buffer

	err = templates.ExecuteTemplate(&b, "index.html", items)
	check(err, "Could not process template.")
	b.WriteTo(w)
}

func ItemHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates := setup(w, r)

	vars := mux.Vars(r)

	id := vars["id"]
	intID:=decodeID(id)
	key := datastore.NewKey(c,"Item","",intID,nil)

	var item Item
	err := datastore.Get(c, key, &item)
	//var items []Item
	//keys, err := q.GetAll(c, &items)
	check(err, "Could not load item.")

	/*
	items[0].itemKey = keys[0]
	items[0].loadOwner(c)
	items[0].loadComments(c)
	*/
	item.itemKey = key
	item.loadOwner(c)
	item.loadComments(c)
	
	var b bytes.Buffer

	err = templates.ExecuteTemplate(&b, "item.html", item)
	check(err, "Could not process template.")
	b.WriteTo(w)
}

func UserHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates := setup(w, r)

	vars := mux.Vars(r)
	username := vars["username"]
	user := getUser(c, username) // TODO: do an actual lookup

	var b bytes.Buffer

	err := templates.ExecuteTemplate(&b, "user.html", user)
	check(err, "Could not process template.")
	b.WriteTo(w)
}
