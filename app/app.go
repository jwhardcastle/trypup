// Google AppEngine initialization
package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Get started with AppEngine
func init() {
	r := new(mux.Router)

	// TODO: put appHandler(...) or some equivalent back, so we safely handle error messages
	r.Handle("/share/{id}/{urltitle}", appHandler(ItemHandler))
	r.Handle("/user/{username}", appHandler(UserHandler))
	r.Handle("/dummy", appHandler(DummyHandler))
	r.Handle("/", appHandler(RootHandler))

	http.Handle("/", r)
}
