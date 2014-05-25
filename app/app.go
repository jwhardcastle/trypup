// Google AppEngine initialization
package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/mjibson/appstats"
)

// Get started with AppEngine
func init() {
	r := new(mux.Router)

	// TODO: put appHandler(...) or some equivalent back, so we safely handle error messages
	r.Handle("/login/", appstats.NewHandler(appHandler(LoginHandler)))
	r.Handle("/share/{id}/{urltitle}", appstats.NewHandler(appHandler(ItemHandler)))
	r.Handle("/user/{username}", appstats.NewHandler(appHandler(UserHandler)))
	r.Handle("/dummy", appstats.NewHandler(appHandler(DummyHandler)))
	r.Handle("/", appstats.NewHandler(appHandler(RootHandler)))

	http.Handle("/", r)
}
