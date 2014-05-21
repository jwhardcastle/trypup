// A sample front page to get started
package app

import (
	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
	"net/http"
)

// Load dummy data
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	dummyData(r)
}

// Show the front page
func RootHandler(w http.ResponseWriter, r *http.Request) {
	templates, _ := setup(w, r)

	c := appengine.NewContext(r)

	var items []Item
	q := datastore.NewQuery("Item").Order("-Score")
	keys, err := q.GetAll(c, &items)
	check(err, "Could not load items.")

	for i, key := range keys {
		items[i].itemKey = key
		items[i].loadOwner(c)
	}

	err = templates.ExecuteTemplate(w, "index.html", items)
	check(err, "Could not process template.")
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	templates, c := setup(w, r)

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
	
	err = templates.ExecuteTemplate(w, "item.html", item)
	check(err, "Could not process template.")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	templates, c := setup(w, r)

	vars := mux.Vars(r)
	username := vars["username"]
	user := getUser(c, username) // TODO: do an actual lookup

	err := templates.ExecuteTemplate(w, "user.html", user)
	check(err, "Could not process template.")
}
