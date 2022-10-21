package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andreistefanciprian/go_web_api_demo/backend/dbmodel"
)

// home page
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Home Page\n")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Book Library"))
}

func ViewArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: articles")
	var a dbmodel.Articles
	w.Header().Set("content-type", "application/json")
	err := a.JSONViewArticles(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: create")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	article := &dbmodel.Article{}
	article.FromJSON(r.Body)

	err := article.AddArticle()
	if err != nil {
		http.Error(w, "Coudn't add article to db.", http.StatusInternalServerError)
	}
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: update article")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	article := &dbmodel.Article{}
	// reqBody, _ := ioutil.ReadAll(r.Body)
	article.FromJSON(r.Body)

	err = article.UpdateArticle(id)
	if err != nil {
		http.Error(w, "Coudn't add article to db.", http.StatusInternalServerError)
	}
}

// start server
func StartServer() {
	log.Print("Listening on port 8080 ...")
	http.Handle("/", dbmodel.JwtAuthentication(homePage))
	http.Handle("/articles", dbmodel.JwtAuthentication(ViewArticles))
	http.Handle("/article/create", dbmodel.JwtAuthentication(CreateArticle))
	http.Handle("/article/delete", dbmodel.JwtAuthentication(dbmodel.DeleteArticle))
	http.Handle("/article/view", dbmodel.JwtAuthentication(dbmodel.ViewArticle))
	http.Handle("/article/update", dbmodel.JwtAuthentication(UpdateArticle))
	http.Handle("/articles/delete_all", dbmodel.JwtAuthentication(dbmodel.DeleteArticles))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
