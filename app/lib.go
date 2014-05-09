// Methods and datastructures that support the application
package app

import (
	"appengine"
	"html/template"
	"net/http"
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

// Make sure we're ready to go, with Content-Type and more
func setup(w http.ResponseWriter, r *http.Request) *Page {
	w.Header().Set("Content-Type", "text/html")
	return &Page{Title: "TrypUp: travel, democratized", Content: ""}
}
