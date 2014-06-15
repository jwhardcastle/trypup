package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
		if len(votes) > 0 {
			for _, vote = range votes {
				if items[i].itemKey.IntID() == vote.ParentKey.IntID() {
					items[i].SessionUserVote = vote.Value
				}
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

	var votes []Vote
	q := datastore.NewQuery("Vote").Filter("OwnerKey=", p.User.userKey).Filter("ParentType=", "Comment")
	q.GetAll(c, &votes) // Eat the error, we don't care if we can't load votes

	var vote Vote

	for _, comment := range item.Comments() {
		if len(votes) > 0 {
			for _, vote = range votes {
				if comment.commentKey.IntID() == vote.ParentKey.IntID() {
					comment.SessionUserVote = vote.Value
				}
			}
		}
	}

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
		log.Print("Comment:", comment.Key())
		newComment = comment.AddComment(c, body, &p.User)
	}

	renderTemplate(w, "_comment", newComment)
}

func VoteHandler(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	p := setup(c, r)

	value64, err := strconv.ParseInt(r.PostFormValue("v"), 10, 8)
	value := int8(value64)
	check(err, "Invalid value.")
	parentID := decodeID(r.PostFormValue("p"))
	parentType := r.PostFormValue("pt")

	var parent Votable
	if parentType == "i" {
		parentItem := GetItem(c, parentID)
		parent = &parentItem
	} else {
		parentComment := GetComment(c, parentID)
		parent = &parentComment
	}

	var votes []Vote
	q := datastore.NewQuery("Vote").Filter("OwnerKey=", p.User.userKey).Filter("ParentKey=", parent.Key())
	keys, err := q.GetAll(c, &votes)
	check(err, "Couldn't load your upvotes.")

	var vote Vote
	if len(votes) > 0 {
		vote = votes[0]
		vote.voteKey = keys[0]
		vote.Parent = parent
		if (&vote).Value == value {
			// delete the current vote
			vote.Delete(c)
			vote.Value = 0
		} else {
			// change the value
			vote.Update(c, value)
			vote.Save(c)
		}
	} else {
		// create a new vote
		vote = *NewVote(c, &p.User, parent, value)
	}

	ret := struct {
		V  int8   `json:"v"`
		P  string `json:"p"`
		PT string `json:"pt"`
	}{vote.Value, r.PostFormValue("p"), parentType}

	var b []byte
	b, err = json.Marshal(ret)
	check(err, "Couldn't format return.")

	w.Write(b)
}
