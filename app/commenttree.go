package app

import (
	"appengine"
	"appengine/datastore"
)

type CommentTree struct {
	comments     []*Comment
	Parent       *datastore.Key
	CommentCount int
	ScoreSum     int
}

func (ct *CommentTree) loadComments(c appengine.Context, key *datastore.Key, recursive bool) []*Comment {
	ct.Parent = key

	q := datastore.NewQuery("Comment").Filter("ParentKey=", key).KeysOnly().Order("-Score")
	keys, err := q.GetAll(c, nil)
	check(err, "Could not load child comments.")

	for _, key := range keys {
		var comment = loadComment(c, key, recursive)
		ct.comments = append(ct.comments, &comment)
	}

	return ct.comments
}

func (ct *CommentTree) Count() int {
	if ct.CommentCount > 0 {
		return ct.CommentCount
	}

	count := len(ct.comments)

	if count > 0 {
		for _, ctchild := range ct.comments {
			count += ctchild.CommentTree.Count()
		}
	}

	ct.CommentCount = count

	return count
}

func (ct CommentTree) Comments() []*Comment {
	return ct.comments
}
