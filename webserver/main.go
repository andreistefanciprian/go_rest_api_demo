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
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", dbmodel.GetAllArticles)
	http.HandleFunc("/article/create", dbmodel.AddArticle)
	http.HandleFunc("/article/delete", dbmodel.RemoveArticle)
	http.HandleFunc("/article/view", dbmodel.ViewArticle)
	http.HandleFunc("/article/update", dbmodel.ChangeArticle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
