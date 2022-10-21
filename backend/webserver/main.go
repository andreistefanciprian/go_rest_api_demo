package webserver

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/andreistefanciprian/go_web_api_demo/backend/dbmodel"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// home page
func homePage(w http.ResponseWriter, r *http.Request) {
	infoLog.Println("Endpoint Hit: /")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Book Library"))
}

func ViewArticles(w http.ResponseWriter, r *http.Request) {
	infoLog.Println("Endpoint Hit: /articles")
	var a dbmodel.Articles
	w.Header().Set("content-type", "application/json")
	err := a.JSONViewArticles(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	infoLog.Println("Endpoint Hit: /article/create")
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
	infoLog.Println("Endpoint Hit: /article/update")
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
		http.Error(w, "Coudn't add article to database.", http.StatusInternalServerError)
	}
}

// start server
func StartServer() {

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.Handle("/", dbmodel.JwtAuthentication(homePage))
	mux.Handle("/articles", dbmodel.JwtAuthentication(ViewArticles))
	mux.Handle("/article/create", dbmodel.JwtAuthentication(CreateArticle))
	mux.Handle("/article/delete", dbmodel.JwtAuthentication(dbmodel.DeleteArticle))
	mux.Handle("/article/view", dbmodel.JwtAuthentication(dbmodel.ViewArticle))
	mux.Handle("/article/update", dbmodel.JwtAuthentication(UpdateArticle))
	mux.Handle("/articles/delete_all", dbmodel.JwtAuthentication(dbmodel.DeleteArticles))

	// create a new server
	var httpPort = ":8080"
	srv := http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// start the server
	infoLog.Println("Listening on port", httpPort)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
