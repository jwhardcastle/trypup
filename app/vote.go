package app

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

type Vote struct {
	Owner       *User `datastore:"-"`
	OwnerKey    *datastore.Key
	Parent      Votable `datastore:"-"`
	ParentKey   *datastore.Key
	ParentType  string
	Value       int8           // 1 for upvote, -1 for downvote
	voteKey     *datastore.Key `datastore:"-"`
	DateCreated time.Time
}

func NewVote(c appengine.Context, owner *User, parent Votable, value int8) *Vote {
	vote := new(Vote)
	vote.Owner = owner
	vote.OwnerKey = (*owner).userKey
	vote.Parent = parent
	vote.ParentKey = parent.Key()
	vote.Value = value
	vote.DateCreated = time.Now()
	vote.ParentType = strings.Split(reflect.TypeOf(parent).String(), ".")[1] // e.g. Item, Comment, not *app.Comment

	vote.id(c)

	err := vote.Save(c)
	check(err, "Couldn't save vote.")

	return vote
}

func (vote *Vote) id(c appengine.Context) {
	// Construct an ID that ensures each user can only vote once on each Model
	id := (*vote).Owner.Username + "$$" + strconv.FormatInt(vote.ParentKey.IntID(), 10)
	(*vote).voteKey = datastore.NewKey(c, "Vote", id, 0, nil)
}

func (vote *Vote) Save(c appengine.Context) error {
	if (*vote).voteKey == nil {
		vote.id(c)
	}

	var err error
	(*vote).voteKey, err = datastore.Put(c, (*vote).voteKey, vote)

	return err
}
