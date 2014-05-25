// Methods and datastructures that support the application
package app

import (
	"appengine"
	"html/template"
	"net/http"
	"math"
)

// Dignified error handling
type appError struct {
	Error   error
	Message string
	Code    int
}

// http.Handle doesn't expect you to return an error, but we want to surface them
type appHandler func(appengine.Context, http.ResponseWriter, *http.Request)

// Serve HTTP, and handle error messages with dignity by displaying a nice error page
func (fn appHandler) ServeHTTP(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recv := recover(); recv != nil {
			e := recv.(*appError)

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
	}()
	fn(c, w, r)
}

// Check errors, panic if we have an error
func check(err error, message string) {
	if err != nil {
		panic(&appError{err, message, 500})
	}
}

// Make sure we're ready to go, with Content-Type and more
func setup(w http.ResponseWriter, r *http.Request) *template.Template {
	w.Header().Set("Content-Type", "text/html")

	templates, err := template.ParseFiles(
		"app/view/header.html",
		"app/view/footer.html",
		"app/view/index.html",
		"app/view/item.html",
		"app/view/login.html",
		"app/view/_comment.html",
		"app/view/_item.html",
		"app/view/user.html",
	)
	check(err, "Could not process templates.")

	return templates
}

// Shamelessly stolen from Reddit, ported to Go
func to_base(q int64, alphabet string) string {
	l := len(alphabet) // The base
	maxdigits := int(math.Ceil(math.Log(float64(q))/math.Log(float64(l))))
	var buffer [64]byte
	var r int // remainder
	var i int
	for i=0 ; q != 0; i++ {
		r = int(math.Mod(float64(q),float64(l)))
		buffer[maxdigits-i-1]=alphabet[r]
		q = q/int64(l)
	}
	return string(buffer[:i])
}

func to36(q int64) string  {
	return to_base(q, "0123456789abcdefghijklmnopqrstuvwxyz")
}
