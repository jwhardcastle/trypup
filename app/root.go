// A sample front page to get started
package app

import (
	//"html/template"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

// Show the front page
func RootHandler(w http.ResponseWriter, r *http.Request) {
	templates := setup(w, r)

	items, _ := dummyData()

	//t, err := templates.ParseFiles("app/view/index.html")
	//check(err, "Could not parse template.")
	
	templates.ExecuteTemplate(w, "index.html", items)
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	templates := setup(w,r)

	items, _ := dummyData()
	vars := mux.Vars(r)

	//t := template.New("Item handler")
	//t = t.Funcs(template.FuncMap{"CommentHandler": CommentHandler})

	//t, err := template.ParseFiles("app/view/item.html", "app/view/header.html", "app/view/footer.html", "app/view/comment.html")
	//check(err, "Could not parse template.")

	id, err := strconv.Atoi(vars["id"])
	check(err, "Invalid identifier.")

	templates.ExecuteTemplate(w, "item.html", items[id-1])
	//t.Execute(w, items[id-1])
}

func CommentHandler(comment Comment) string {
	//t, err := template.ParseFiles("app/view/_comment.html")
	//check(err, "Could not parse comments.")

	//t.Execute(, comment)
	return "!"
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	templates := setup(w,r)

	_, users := dummyData()
	//vars := mux.Vars(r)
	//username := vars["username"]
	user := users[0] // TODO: do an actual lookup

	templates.ExecuteTemplate(w, "user.html", user)
}

func dummyData() ([]Item, []User) {
	u1 := User{Username: "jwhardcastle"}
	u2 := User{Username: "jhutton"}
	u3 := User{Username: "rkavalsky"}
	u4 := User{Username: "teej"}

	c1 := Comment{Owner: u3, Body: "We love going here!", Score: 3, Upvotes: 3, Downvotes: 0}
	c2 := Comment{Owner: u4, Body: "typography geek heaven", Children: []*Comment{}, Score: 5, Upvotes: 5, Downvotes: 0}
	c3 := Comment{Owner: u1, Body: "Agreed! Among other things.", Parent: &c2, Score: 0, Upvotes: 1, Downvotes: 1}

	c2.Children = append(c2.Children, &c3)

	i1 := Item{
		Title: "Baltimore Museum of Industry, learn how a linotype works, among the city's industrial history",
		Comments: []*Comment{&c2},
		Owner: u1,
		URLTitle: "baltimore-museum-of-industry-learn-how-a",
		Id: 1,
		Score: 36,
		Upvotes: 40,
		Downvotes: 4,
		Lat: 39.273556,
		Long: -76.601806,
		CommentCount: 2,
		Color: "cadetblue",
		Icon: "truck",
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
		Lat: 39.283501,
		Long: -76.6219798,
		CommentCount: 1,
		Color: "orange",
		Icon: "sun-o",
	}

	items := []Item{i1,i2}
	users := []User{u1,u2,u3,u4}
	return items, users
}
