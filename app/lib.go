// Methods and datastructures that support the application
package app

import (
	"appengine"
	"html/template"
	"net/http"
	"time"
)

// Dignified error handling
type appError struct {
	Error   error
	Message string
	Code    int
}

// An item is an activity, place, restaurant, or point of interest that has been shared
type Item struct {
	Title		string
	Description	string
	Id			uint
	Lat			float32
	Long		float32
	URLTitle	string
	Icon		string
	Color		string		
	Comments	[]*Comment
	DateCreated	time.Time
	Owner		User
	Score		int
	Upvotes		int
	Downvotes	int
	CommentCount int
}

// Comments belong to items, or to other comments; all comments must reference the root Item regardless
type Comment struct {
	Owner		User
	Body		string
	Children	[]*Comment
	DateCreated	time.Time
	Parent		*Comment
	Item		*Item
	Score		int
	Upvotes		int
	Downvotes	int
}

// Users log in to vote, share, and leave comments
type User struct {
	Username		string
	PasswordHash	string
	Id				uint
	DateCreated		time.Time
}

// http.Handle doesn't expect you to return an error, but we want to surface them
type appHandler func(http.ResponseWriter, *http.Request)

// Serve HTTP, and handle error messages with dignity by displaying a nice error page
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recv := recover(); recv != nil {
			e := recv.(*appError)
			c := appengine.NewContext(r)
			c.Errorf("%v", e.Error)

			http.StatusText(e.Code)
			t, err := template.ParseFiles("errors/500.html")
			if err != nil {
				http.Error(w, "A very serious error has occurred.", 500)
			}
			err = t.Execute(w, e)
			if err != nil {
				http.Error(w, "A very serious error has occurred.", 500)
			}
			
		}  
	} ()
	fn(w,r)
}

func (fn Item) string(i Item) string {
	return i.Title
}

// Check errors, panic if we have an error
func check(err error, message string) { if err != nil { panic(&appError{err, message, 500} ) } }

// Make sure we're ready to go, with Content-Type and more
func setup(w http.ResponseWriter, r *http.Request) *template.Template {
	w.Header().Set("Content-Type", "text/html")

	templates, err := template.ParseFiles(
		"app/view/header.html",
		"app/view/footer.html",
		"app/view/index.html",
		"app/view/item.html",
		"app/view/_comment.html",
		"app/view/_item.html",
		"app/view/user.html",
	)
	check(err, "Could not process templates.")

	return templates
}
