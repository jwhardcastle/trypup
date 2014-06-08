// Methods and datastructures that support the application
package app

import (
	"html/template"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mjibson/appstats"

	"appengine"
	"appengine/datastore"
)

// Dignified error handling
type appError struct {
	Error   error
	Message string
	Code    int
}

type page struct {
	LoggedIn bool
	User     User
	Title    string
	Session  *sessions.Session
	Data     []interface{}
	Flashes  []interface{}
}

type Commentable interface {
	loadComments() CommentTree
}

type Votable interface {
	Key() *datastore.Key
	CountVotes(appengine.Context)
}

// http.Handle doesn't expect you to return an error, but we want to surface them
type AppHandler struct {
	f appstats.Handler
}

func NewAppHandler(f func(appengine.Context, http.ResponseWriter, *http.Request)) AppHandler {
	return AppHandler{
		f: appstats.NewHandler(f),
	}
}

// Serve HTTP, and handle error messages with dignity by displaying a nice error page
func (a AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recv := recover(); recv != nil {

			e := recv.(*appError)

			log.Print(e.Error)

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
	a.f.ServeHTTP(w, r)
}

// Check errors, panic if we have an error
func check(err error, message string) {
	if err != nil {
		panic(&appError{err, message, 500})
	}
}

// Make sure we're ready to go, with Content-Type and more
func setup(c appengine.Context, w http.ResponseWriter, r *http.Request) (*template.Template, page) {
	w.Header().Set("Content-Type", "text/html")

	defer func() {
		if recv := recover(); recv != nil {
			err := recv.(error)
			check(err, "Could not process templates.")
		}
	}()

	templates := template.Must(template.New("").Funcs(template.FuncMap{
	//"printFlashes": printFlashes,
	}).ParseGlob("app/view/*.html"))

	var p page
	var err error

	p.Session, err = store.Get(r, "trypup")
	check(err, "Couldn't load session.")

	if p.Session.Values["Username"] != nil {
		p.User, err = getUser(c, p.Session.Values["Username"].(string))
		p.LoggedIn = true
		check(err, "Could not load your user profile.")
	} else {
		p.LoggedIn = false
	}
	p.Flashes = p.Session.Flashes()

	return templates, p
}

// Shamelessly stolen from Reddit, ported to Go
func to_base(q int64, alphabet string) string {
	l := len(alphabet) // The base
	maxdigits := int(math.Ceil(math.Log(float64(q)) / math.Log(float64(l))))
	var buffer [64]byte
	var r int // remainder
	var i int
	for i = 0; q != 0; i++ {
		r = int(math.Mod(float64(q), float64(l)))
		buffer[maxdigits-i-1] = alphabet[r]
		q = q / int64(l)
	}
	return string(buffer[:i])
}

func to36(q int64) string {
	return to_base(q, "0123456789abcdefghijklmnopqrstuvwxyz")
}
