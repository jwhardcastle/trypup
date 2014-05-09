// A sample front page to get started
package app

import (
	"html/template"
	"net/http"
)

// Show the front page
func root(w http.ResponseWriter, r *http.Request) {
	p := setup(w, r)
	
	p.Content = template.HTML(`<h1>TrypUp: travel, democratized</h1>`)
	t, err := template.ParseFiles("app/view/index.html")
	check(err, "Could not parse template.")
	
	t.Execute(w, p)
}