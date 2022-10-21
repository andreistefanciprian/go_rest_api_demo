package webserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/andreistefanciprian/go_web_api_demo/backend/dbmodel"
)

var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// home page
func homePage(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Book Library"))
}

func ViewArticles(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /articles")
	var a dbmodel.Articles
	w.Header().Set("content-type", "application/json")
	err := a.JSONViewArticles(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /article/create")
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
	InfoLog.Println("Endpoint Hit: /article/update")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		ErrorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	article := &dbmodel.Article{}
	// reqBody, _ := ioutil.ReadAll(r.Body)
	article.FromJSON(r.Body)

	err = article.UpdateArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Article was updated successfully."))
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /article/delete")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		ErrorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	err = dbmodel.DeleteArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Article was deleted successfully."))
	}

}

func ViewArticle(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /article/view")
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		ErrorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	article := &dbmodel.Article{}
	err = article.GetArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Header().Set("content-type", "application/json")
		article.ToJSON(w)
	}
}

func DeleteArticles(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Endpoint Hit: /articles/delete_all")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var articles dbmodel.Articles
	err := articles.DeleteArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		msg := fmt.Sprintf("Deleted %d articles.", len(articles))
		w.Write([]byte(msg))
	}

}

// start server
func StartServer() {

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.Handle("/", dbmodel.JwtAuthentication(homePage))
	mux.Handle("/articles", dbmodel.JwtAuthentication(ViewArticles))
	mux.Handle("/article/create", dbmodel.JwtAuthentication(CreateArticle))
	mux.Handle("/article/delete", dbmodel.JwtAuthentication(DeleteArticle))
	mux.Handle("/article/view", dbmodel.JwtAuthentication(ViewArticle))
	mux.Handle("/article/update", dbmodel.JwtAuthentication(UpdateArticle))
	mux.Handle("/articles/delete_all", dbmodel.JwtAuthentication(DeleteArticles))

	// create a new server
	var httpPort = ":8080"
	srv := http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// start the server
	InfoLog.Println("Listening on port", httpPort)
	err := srv.ListenAndServe()
	if err != nil {
		ErrorLog.Fatal(err)
	}
}
