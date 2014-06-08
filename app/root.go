// A sample front page to get started
package app

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"
)

// Load dummy data
func DummyHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	dummyData(r, c)
}

// Show the front page
func RootHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates, p := setup(c, w, r)

	var items []Item
	q := datastore.NewQuery("Item").Order("-Score")
	keys, err := q.GetAll(c, &items)
	check(err, "Could not load items.")

	for i, key := range keys {
		items[i].itemKey = key
		items[i].loadOwner(c)
	}

	p.Data = append(p.Data, items)

	var b bytes.Buffer

	err = templates.ExecuteTemplate(&b, "index.html", p)
	check(err, "Could not process template.")
	b.WriteTo(w)
}

func ItemHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates, p := setup(c, w, r)

	vars := mux.Vars(r)

	id := vars["id"]
	intID := decodeID(id)
	key := datastore.NewKey(c, "Item", "", intID, nil)

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

	p.Data = append(p.Data, item)

	var b bytes.Buffer

	err = templates.ExecuteTemplate(&b, "item.html", p)
	check(err, "Could not process template.")
	b.WriteTo(w)
}

func LoginHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates, p := setup(c, w, r)

	err := r.ParseForm()
	check(err, "Could not process login information.")

	if len(r.PostForm) > 0 {
		log.Print(r.PostForm)
		if len(r.PostForm["Username"][0]) > 0 {
			user, err := getUser(c, r.PostForm["Username"][0])
			if err != nil {
				p.Session.AddFlash("Could not find that username.")

			} else {
				if len(r.PostForm["Password"][0]) > 0 {
					err = user.checkPassword(r.PostForm["Password"][0])

					if err != nil {
						p.Session.AddFlash("Your password is incorrect.")

					} else {
						p.Session.Values["Username"] = user.Username
						p.User = user
						p.Session.Save(r, w)
						http.Redirect(w, r, "/", 302)

					}

				} else {
					p.Session.AddFlash("Password may not be blank.")

				}

			}
		} else {
			p.Session.AddFlash("Username may not be blank.")

		}

		p.Session.Save(r, w)
	}

	if flashes := p.Session.Flashes(); len(flashes) > 0 {
		for _, f := range flashes {
			p.Flashes = append(p.Flashes, f)
		}
	}

	var b bytes.Buffer
	err = templates.ExecuteTemplate(&b, "login.html", p)
	b.WriteTo(w)
}

func UserHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	templates, p := setup(c, w, r)

	vars := mux.Vars(r)
	username := vars["username"]
	user, err := getUser(c, username) // TODO: do an actual lookup

	check(err, "A user with that name could not be found.")

	p.Data = append(p.Data, user)

	var b bytes.Buffer

	err = templates.ExecuteTemplate(&b, "user.html", p)
	check(err, "Could not process template.")
	b.WriteTo(w)
}
