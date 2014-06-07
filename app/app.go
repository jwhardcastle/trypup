// Google AppEngine initialization
package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Get started with AppEngine
func init() {
	r := new(mux.Router)

	r.Handle("/login/", NewAppHandler(LoginHandler))
	r.Handle("/share/{id}/{urltitle}", NewAppHandler(ItemHandler))
	r.Handle("/user/{username}", NewAppHandler(UserHandler))
	r.Handle("/dummy", NewAppHandler(DummyHandler))
	r.Handle("/", NewAppHandler(RootHandler))

	http.Handle("/", r)
}
