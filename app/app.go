// Google AppEngine initialization
package app

import (
	"net/http"
)

// Get started with AppEngine
func init() {
	http.Handle("/", appHandler(root))
}

