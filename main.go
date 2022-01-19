package main

import (
	"database/sql"
	"fmt"

	"github.com/gorilla/mux"

	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                             uint16
	Title, Announcement, Full_text string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_test")

	if err != nil {
		panic(err)
	}
	defer db.Close()
	//выборка данных
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Announcement, &post.Full_text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
		//fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		panic(err)
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	announcement := r.FormValue("announcement")
	full_text := r.FormValue("full_text")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_test")

	if err != nil {
		panic(err)
	}
	defer db.Close()
	//установка данных
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`,`announcement`,`full_text`) VALUES ('%s','%s','%s')", title, announcement, full_text))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang_test")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = %s", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Announcement, &post.Full_text)
		if err != nil {
			panic(err)
		}
		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)

}

func contacts(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		panic(err)
	}

	t.ExecuteTemplate(w, "contacts", nil)
}

func handleFunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":3306", nil)
}

func main() {
	handleFunc()
}
