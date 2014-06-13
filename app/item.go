// An item is an activity, place, restaurant, or point of interest that has been shared

package app

import (
	"math"
	"regexp"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

type Item struct {
	Title       string
	Description string
	Lat         float32
	Long        float32
	URLTitle    string
	Icon        string
	Color       string
	DateCreated time.Time
	Score       int
	Upvotes     int
	Downvotes   int
	//CommentCount int
	owner    *User          `datastore:"-"`
	OwnerKey *datastore.Key `datastore:"owner"`
	comments []*Comment     `datastore:"-"`
	itemKey  *datastore.Key `datastore:"-"`
	CommentTree
	SessionUserVote int8 `datastore:"-"`
}

// TODO: shard counters for votes https://developers.google.com/appengine/articles/sharding_counters

func (fn Item) string(i Item) string {
	return i.Title
}

func NewItem(c appengine.Context, title string, description string, icon string, color string, lat float32, long float32, owner *User) *Item {
	item := new(Item)
	item.Title = title
	item.Description = description
	item.Icon = icon
	item.Color = color
	item.Lat = lat
	item.Long = long
	item.owner = owner
	item.OwnerKey = (*owner).userKey

	titleRegEx, _ := regexp.Compile("[^a-z0-9]")
	item.URLTitle = string(titleRegEx.ReplaceAll([]byte(strings.ToLower(title)), []byte("-")))

	item.DateCreated = time.Now()

	// The submitter automatically upvotes the new item
	item.Score = 1
	item.Upvotes = 1

	item.Save(c)

	NewVote(c, owner, item, 1)

	return item
}

func GetItem(c appengine.Context, intID int64) Item {
	var item Item

	key := datastore.NewKey(c, "Item", "", intID, nil)
	err := datastore.Get(c, key, &item)
	check(err, "Could not find item.")

	item.itemKey = key

	return item
}

func (item *Item) Save(c appengine.Context) error {
	if (*item).itemKey == nil {
		item.itemKey = datastore.NewIncompleteKey(c, "Item", nil)
	}
	var err error
	(*item).itemKey, err = datastore.Put(c, (*item).itemKey, item)
	return err
}

func (item *Item) AddComment(c appengine.Context, body string, owner *User) *Comment {
	return NewComment(c, body, owner, item)
}

func (item *Item) loadComments(c appengine.Context) {
	var comments []*Comment

	comments = item.CommentTree.loadComments(c, item.itemKey, true)

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

func (item *Item) CountVotes(c appengine.Context) {
	q := datastore.NewQuery("Vote").Filter("ParentKey=", (*item).itemKey)
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
	(*item).Score = s
	(*item).Upvotes = u
	(*item).Downvotes = d
	item.Save(c)
}

func (item Item) Upvoted() bool {
	if item.SessionUserVote > 0 {
		return true
	}
	return false
}

func (item Item) Downvoted() bool {
	if item.SessionUserVote < 0 {
		return true
	}
	return false
}

func (item Item) Owner() User {
	// TODO: figure out how to make this lazy load
	//if item.owner == nil {
	//err:=datastore.Get(c, item.parent, &item.owner)
	//check(err, "Could not load comments.")
	//}
	return *(item.owner)
}

func (item Item) Comments() []*Comment {
	// TODO: figure out how to make this lazy load
	//if (len(item.comments) == 0) {
	//item.loadComments() // TODO, no context here, can't load
	//}
	return item.comments
}

func (item Item) URL() string {
	return "/share/" + item.URLID() + "/" + item.URLTitle
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
	for i := 0; i < len(id); i++ {
		m := strings.IndexByte(alphabet, id[i])
		output = output + int64(float64(m)*(math.Pow(float64(36), float64(len(id)-1-i))))
	}
	return output
}
