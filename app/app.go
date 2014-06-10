// Google AppEngine initialization
package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oxtoacart/bpool"
)

// Get started with AppEngine
func init() {

	// Buffers are used to hold page content while we build, and then release all at once
	bufpool = bpool.NewBufferPool(32)

	// Load the various .html templates into memory
	initTemplates()

	r := new(mux.Router)

	r.Handle("/login/", NewAppHandler(LoginHandler))
	r.Handle("/share/{id}/{urltitle}", NewAppHandler(ItemHandler))
	r.Handle("/user/{username}", NewAppHandler(UserHandler))
	r.Handle("/dummy", NewAppHandler(DummyHandler))
	r.Handle("/", NewAppHandler(RootHandler))

	http.Handle("/", r)
}
