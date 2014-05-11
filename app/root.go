// A sample front page to get started
package app

import (
	"html/template"
	"net/http"
)


// Show the front page
func root(w http.ResponseWriter, r *http.Request) {
	setup(w, r)

	u1 := User{Username: "jwhardcastle"}
	u2 := User{Username: "jhutton"}
	u3 := User{Username: "rkavalsky"}
	u4 := User{Username: "teej"}

	c1 := Comment{Owner: u3, Body: "We love going here!"}
	c2 := Comment{Owner: u4, Body: "typography geek heaven", Children: []*Comment{}}
	//c3 := Comment{Owner: u1, Body: "Agreed! Among other things.", Parent: &c2}

	//c2.Children[0] = &c3

	i1 := Item{
		Title: "Baltimore Museum of Industry, learn how a linotype works, among the city's industrial history",
		Comments: []*Comment{&c2},
		Owner: u1,
		URLTitle: "baltimore-museum-of-industry-learn-how-a",
		Id: 1,
		Score: 36,
		Upvotes: 40,
		Downvotes: 4,
	}
	i2 := Item{
		Title: "OPACY: Oriole Park at Camden Yards, Home of the Baltimore Orioles",
		Comments: []*Comment{&c1},
		Owner: u2,
		URLTitle: "opacy-oriole-park-at-camden-yards-home-o",
		Id: 2,
		Score: 129,
		Upvotes: 150,
		Downvotes: 11,
	}

	items := []Item{i2,i1}

	//p.Content = template.HTML(`<h1>TrypUp: travel, democratized</h1>`)
	t, err := template.ParseFiles("app/view/index.html")
	check(err, "Could not parse template.")
	
	t.Execute(w, items)
}