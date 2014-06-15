// Comments belong to items, or to other comments; all comments must reference the root Item regardless

package app

import (
	"html/template"
	"log"
	"time"

	"github.com/russross/blackfriday"

	"appengine"
	"appengine/datastore"
)

type Comment struct {
	Body            string `datastore:"Body,noindex"`
	DateCreated     time.Time
	Score           int
	Upvotes         int
	Downvotes       int
	owner           *User
	OwnerKey        *datastore.Key
	parent          Votable
	ParentKey       *datastore.Key
	item            *Item
	ItemKey         *datastore.Key
	SessionUserVote int8 `datastore:"-"`
	children        []*Comment
	commentKey      *datastore.Key
	CommentTree
}

func NewComment(c appengine.Context, body string, owner *User, parent Votable) *Comment {
	comment := new(Comment)
	comment.Body = body
	comment.owner = owner
	comment.OwnerKey = owner.userKey
	comment.parent = parent
	comment.ParentKey = parent.Key()

	// We want a direct link to the top Item for this tree
	parentItem, ok := parent.(*Item)
	if ok {
		comment.item = parentItem
		comment.ItemKey = parentItem.Key()
	} else {
		parentComment := parent.(*Comment)
		comment.item = parentComment.item
		comment.ItemKey = parentComment.ItemKey
	}

	comment.DateCreated = time.Now()
	comment.Score = 1
	comment.Upvotes = 1
	err := comment.Save(c)
	check(err, "Couldn't save your comment.")

	log.Print("New Comment:", comment.Key())

	// The submitter automatically upvotes the new comment
	NewVote(c, owner, comment, 1)

	return comment
}

func GetComment(c appengine.Context, intID int64) Comment {
	var comment Comment

	key := datastore.NewKey(c, "Comment", "", intID, nil)
	err := datastore.Get(c, key, &comment)
	check(err, "Could not find comment.")

	comment.commentKey = key

	return comment
}

func (comment *Comment) Save(c appengine.Context) error {
	if (*comment).commentKey == nil {
		(*comment).commentKey = datastore.NewIncompleteKey(c, "Comment", nil)
		log.Print("New comment key", (*comment).commentKey)
	}

	var err error
	(*comment).commentKey, err = datastore.Put(c, (*comment).commentKey, comment)

	return err
}

func (comment *Comment) AddComment(c appengine.Context, body string, owner *User) *Comment {
	return NewComment(c, body, owner, comment)
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

func (comment Comment) Upvoted() bool {
	if comment.SessionUserVote > 0 {
		return true
	}
	return false
}

func (comment Comment) Downvoted() bool {
	if comment.SessionUserVote < 0 {
		return true
	}
	return false
}

func (comment Comment) Format() template.HTML {
	output := blackfriday.MarkdownBasic([]byte(comment.Body))
	return template.HTML(string(output))
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

func (comment *Comment) Upvote(c appengine.Context) {
	comment.Score += 1
	comment.Upvotes += 1
}

func (comment *Comment) DeUpvote(c appengine.Context) {
	comment.Score -= 1
	comment.Upvotes -= 1
}

func (comment *Comment) Downvote(c appengine.Context) {
	comment.Score -= 1
	comment.Downvotes += 1
}

func (comment *Comment) DeDownvote(c appengine.Context) {
	comment.Score += 1
	comment.Downvotes -= 1
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

func (comment Comment) URL() string {
	return comment.item.URL() + "/" + comment.URLID()
}

// Take the IntID and convert it to base36 for use in URLs, etc.
func (comment Comment) URLID() string {
	return to36(comment.commentKey.IntID())
}
