// An item is an activity, place, restaurant, or point of interest that has been shared

package app

import (
	"appengine"
	"appengine/datastore"
	"time"
	"strings"
	"math"
)

type Item struct {
	Title        string
	Description  string
	Lat          float32
	Long         float32
	URLTitle     string
	Icon         string
	Color        string
	DateCreated  time.Time
	Score        int
	Upvotes      int
	Downvotes    int
	CommentCount int
	owner        *User `datastore:"-"`
	OwnerKey     *datastore.Key `datastore:"owner"`
	comments     []*Comment `datastore:"-"`
	itemKey      *datastore.Key `datastore:"-"`
}

// TODO: shard counters for votes https://developers.google.com/appengine/articles/sharding_counters

func (fn Item) string(i Item) string {
	return i.Title
}

func (item *Item) loadComments(c appengine.Context) {
	var comments []*Comment
	q := datastore.NewQuery("Comment").Filter("ParentKey=", item.itemKey).Order("-Score")
	keys, err := q.GetAll(c, &comments)

	for i, key := range keys {
		comments[i].commentKey = key
		comments[i].loadOwner(c)
		comments[i].loadChildren(c, true)
	}

	check(err, "Could not load comments.")
	(*item).comments = comments
}

func (item *Item) loadOwner(c appengine.Context) {
	var u User
	err := datastore.Get(c, (*item).OwnerKey, &u)
	(*item).owner = &u
	check(err, "Could not load owner.")
}

func (item Item) Key() *datastore.Key {
	return item.itemKey
}

func (item Item) Owner() User {
	// TODO: figure out how to make this lazy load, maybe from itemKey.parent()
	//if item.owner == nil {
		//err:=datastore.Get(c, item.parent, &item.owner)
		//check(err, "Could not load comments.")
	//}
	return *(item.owner) 
}

// Lazy-load comments
func (item Item) Comments() []*Comment {
	if (len(item.comments) == 0) {
		//item.loadComments() // TODO, no context here, can't load
	}
	return item.comments	
}

// Take the IntID and convert it to base36 for use in URLs, etc.
func (item Item) URLID() string {
	return to36(item.itemKey.IntID())
}

// Take a base36 string and put it back to int
func decodeID(id string) int64 {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	var output int64
	output = 0
	for i:=0; i < len(id) ; i++ {
		m := strings.IndexByte(alphabet, id[i])
		output = output + int64(float64(m) * (math.Pow(float64(36), float64(len(id)-1-i))) )
	}
	return output
}
