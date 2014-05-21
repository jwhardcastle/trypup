// Comments belong to items, or to other comments; all comments must reference the root Item regardless

package app

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Comment struct {
	Body          string
	DateCreated   time.Time
	Score         int
	Upvotes       int
	Downvotes     int
	owner         *User
	OwnerKey      *datastore.Key
	parentComment *Comment
	parentitem    *Item
	ParentKey     *datastore.Key
	children      []*Comment 
	commentKey    *datastore.Key
}

func loadComment(c appengine.Context, key *datastore.Key) Comment {
	var comment Comment
	err := datastore.Get(c, key, &comment)
	check(err, "Could not load comment.")

	comment.loadOwner(c)

	return comment
}

func (comment *Comment) loadOwner(c appengine.Context) {
	var u User
	err := datastore.Get(c, (*comment).OwnerKey, &u)
	check(err, "Could not load comment owner.")

	(*comment).owner = &u
	
}

func (comment *Comment) loadChildren(c appengine.Context, recursive bool) {
	var children []Comment
	q := datastore.NewQuery("Comment").Filter("ParentKey=", (*comment).commentKey).Order("-Score")
	keys, err := q.GetAll(c, &children)
	check(err, "Could not load child comments.")


	var childs []*Comment
	for i, key := range keys {
		children[i].commentKey = key
		children[i].loadOwner(c)

		if(recursive) {
			children[i].loadChildren(c, recursive)
		}

		childs = append(childs, &(children[i]))
	}

	(*comment).children = childs
}

func (comment Comment) Owner() User {
	return *(comment.owner)
}

func (comment Comment) Children() []*Comment {
	return comment.children
}