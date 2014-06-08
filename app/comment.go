// Comments belong to items, or to other comments; all comments must reference the root Item regardless

package app

import (
	"time"

	"appengine"
	"appengine/datastore"
)

type Comment struct {
	Body        string
	DateCreated time.Time
	Score       int
	Upvotes     int
	Downvotes   int
	owner       *User
	OwnerKey    *datastore.Key
	parent      Votable
	ParentKey   *datastore.Key
	children    []*Comment
	commentKey  *datastore.Key
	CommentTree
}

func NewComment(c appengine.Context, body string, owner *User, parent Votable) *Comment {
	comment := new(Comment)
	comment.Body = body
	comment.owner = owner
	comment.OwnerKey = owner.userKey
	comment.parent = parent
	comment.ParentKey = parent.Key()
	comment.Save(c)

	return comment
}

func (comment *Comment) Save(c appengine.Context) error {
	if (*comment).commentKey == nil {
		(*comment).commentKey = datastore.NewIncompleteKey(c, "Comment", nil)
	}

	var err error
	(*comment).commentKey, err = datastore.Put(c, (*comment).commentKey, comment)

	return err
}

func (comment *Comment) CountVotes(c appengine.Context) {
	q := datastore.NewQuery("Vote").Filter("ParentKey=", (*comment).commentKey)
	var votes []Vote
	_, err := q.GetAll(c, &votes)
	check(err, "Couldn't load votes.")

	var s, u, d int
	for _, vote := range votes {
		s += int(vote.Value)

		if vote.Value > 0 {
			u++
		} else {
			d++
		}
	}
	(*comment).Score = s
	(*comment).Upvotes = u
	(*comment).Downvotes = d
	comment.Save(c)
}

func loadComment(c appengine.Context, key *datastore.Key, recursive bool) Comment {
	var comment Comment
	err := datastore.Get(c, key, &comment)
	check(err, "Could not load comment.")

	comment.commentKey = key
	comment.loadOwner(c)

	if recursive {
		comment.children = comment.CommentTree.loadComments(c, key, recursive)
	}

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

		if recursive {
			children[i].loadChildren(c, recursive)
		}

		childs = append(childs, &(children[i]))
	}

	(*comment).children = childs
}

func (comment Comment) Key() *datastore.Key {
	return comment.commentKey
}

func (comment Comment) Owner() User {
	return *(comment.owner)
}

func (comment Comment) Children() []*Comment {
	return comment.children
}
