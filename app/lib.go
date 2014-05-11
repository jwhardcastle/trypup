// Methods and datastructures that support the application
package app

import (
	"appengine"
	"html/template"
	"net/http"
	"time"
)

// Information we feed to the render template
type Page struct {
	Title   string
	Content template.HTML
}

// Dignified error handling
type appError struct {
	Error   error
	Message string
	Code    int
}

type Item struct {
	Title		string
	Description	string
	Id			uint
	Lat			float32
	Long		float32
	URLTitle	string
	Comments	[]*Comment
	DateCreated	time.Time
	Owner		User
	Score		int
	Upvotes		int
	Downvotes	int
}

type Comment struct {
	Owner		User
	Body		string
	Children	[]*Comment
	DateCreated	time.Time
	Parent		*Comment
}

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
func setup(w http.ResponseWriter, r *http.Request) *Page {
	w.Header().Set("Content-Type", "text/html")
	return &Page{Title: "TrypUp: travel, democratized", Content: ""}
}
