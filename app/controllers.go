package app

import (
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
	p := setup(c, r)

	var items []Item
	q := datastore.NewQuery("Item").Order("-Score")
	keys, err := q.GetAll(c, &items)
	check(err, "Could not load items.")

	var votes []Vote
	q = datastore.NewQuery("Vote").Filter("OwnerKey=", p.User.userKey).Filter("ParentType=", "Item")
	q.GetAll(c, &votes) // Eat the error, we don't care if we can't load votes

	var vote Vote

	for i, key := range keys {
		items[i].itemKey = key
		items[i].loadOwner(c)
		for i, vote = range votes {
			if items[i].itemKey.IntID() == vote.ParentKey.IntID() {
				items[i].SessionUserVote = vote.Value
			}
		}
	}

	p.Data["Items"] = items

	renderTemplate(w, "index.html", p)
}

// The detail page for a specific item
func ItemHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	p := setup(c, r)

	vars := mux.Vars(r)

	id := vars["id"]
	item := GetItem(c, decodeID(id))
	item.loadOwner(c)
	item.loadComments(c)

	// Since we've already loaded the comments, count 'em up and store for later reference
	item.CommentTree.Count()
	item.Save(c)

	p.Data["Item"] = item

	renderTemplate(w, "item.html", p)
}

// Allow a user to login to the site
func LoginHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	p := setup(c, r)

	err := r.ParseForm()
	check(err, "Could not process login information.")

	if len(r.PostForm) > 0 {
		log.Print(r.PostForm)
		if len(r.PostForm["Username"][0]) > 0 {
			user, err := GetUser(c, r.PostForm["Username"][0])
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

	renderTemplate(w, "login.html", p)
}

// Show the detail page for a user, with their submitted items, etc.
func UserHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	p := setup(c, r)

	vars := mux.Vars(r)
	username := vars["username"]
	user, err := GetUser(c, username) // TODO: do an actual lookup
	check(err, "A user with that name could not be found.")

	p.Data["User"] = user

	renderTemplate(w, "user.html", p)
}

func AddComment(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	p := setup(c, r)

	body := r.PostFormValue("b")
	parent := decodeID(r.PostFormValue("p"))
	parentType := r.PostFormValue("pt")

	var newComment *Comment

	if parentType == "i" {
		item := GetItem(c, parent)
		newComment = item.AddComment(c, body, &p.User)
	} else {
		comment := GetComment(c, parent)
		newComment = comment.AddComment(c, body, &p.User)
	}

	renderTemplate(w, "_comment", newComment)
}
