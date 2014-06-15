package app

import (
	"log"
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
	log.Print("Vote:", vote)

	vote.id(c)

	err := vote.Save(c)
	check(err, "Couldn't save vote.")

	return vote
}

func (vote *Vote) id(c appengine.Context) {
	//log.Print(vote.Owner)
	log.Print("Parent: ", vote.ParentKey)
	//return
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

func (vote *Vote) Update(c appengine.Context, newValue int8) error {
	if vote.Value > 0 && newValue < 0 {
		log.Print((*vote).Parent)
		(*vote).Parent.DeUpvote(c)
	} else if vote.Value < 0 && newValue > 0 {
		(*vote).Parent.DeDownvote(c)
	}
	if newValue > 0 {
		(*vote).Parent.Upvote(c)
	} else if newValue < 0 {
		(*vote).Parent.Downvote(c)
	}

	(*vote).Parent.Save(c)

	vote.Value = newValue
	return vote.Save(c)
}

func (vote *Vote) Delete(c appengine.Context) error {
	if (*vote).Value == 1 {
		(*vote).Parent.DeUpvote(c)
		(*vote).Parent.Save(c)
	} else {
		(*vote).Parent.DeDownvote(c)
		(*vote).Parent.Save(c)
	}
	return datastore.Delete(c, (*vote).voteKey)
}
