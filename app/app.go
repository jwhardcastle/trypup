// Google AppEngine initialization
package app

import (
	"net/http"
	"github.com/gorilla/mux"
)

// Get started with AppEngine
func init() {
	r := new(mux.Router)

	// TODO: put appHandler(...) or some equivalent back, so we safely handle error messages
	r.Handle("/share/{id:[0-9]+}/{urltitle}", appHandler(ItemHandler))
	r.Handle("/user/{username}", appHandler(UserHandler))
	r.Handle("/", appHandler(RootHandler))

	http.Handle("/", r)
}

