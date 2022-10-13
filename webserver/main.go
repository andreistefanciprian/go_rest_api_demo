package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andreistefanciprian/go_web_api_demo/dbmodel"
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

// start server
func StartServer() {
	log.Print("Listening on port 8080 ...")
	http.Handle("/", dbmodel.JwtAuthentication(homePage))
	http.Handle("/articles", dbmodel.JwtAuthentication(dbmodel.ViewArticles))
	http.Handle("/article/create", dbmodel.JwtAuthentication(dbmodel.CreateArticle))
	http.Handle("/article/delete", dbmodel.JwtAuthentication(dbmodel.DeleteArticle))
	http.Handle("/article/view", dbmodel.JwtAuthentication(dbmodel.ViewArticle))
	http.Handle("/article/update", dbmodel.JwtAuthentication(dbmodel.UpdateArticle))
	http.Handle("/articles/delete_all", dbmodel.JwtAuthentication(dbmodel.DeleteArticles))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
