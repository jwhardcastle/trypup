// A sample front page to get started
package root

import (
	"appengine"
	"fmt"
	"html/template"
	"net/http"
)

// Information we feed to the render template
type Page struct {
	Title   string
	Content string
}

// Dignified error handling
type appError struct {
	Error   error
	Message string
	Code    int
}

// Get started with AppEngine
func init() {
	http.Handle("/", appHandler(root))
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

// Check errors, panic if we have an error
func check(err error, message string) { if err != nil { panic(&appError{err, message, 500} ) } }

// Show the front page
func root(w http.ResponseWriter, r *http.Request) {
	p := setup(w, r)
	
	p.Content = fmt.Sprintf("%s", "<h1>TrypUp: travel, democratized</h1>")
	t, err := template.ParseFiles("root/view/index.html")
	check(err, "Could not parse template.")
	
	t.Execute(w, p)
}

// Make sure we're ready to go, with Content-Type and more
func setup(w http.ResponseWriter, r *http.Request) *Page {
	w.Header().Set("Content-Type", "text/html")
	return &Page{Title: "TrypUp: travel, democratized", Content: ""}
}
