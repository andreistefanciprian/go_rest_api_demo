package webserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/andreistefanciprian/go_web_api_demo/backend/dbmodel"
	jwt "github.com/golang-jwt/jwt/v4"
)

var mySigningKey = []byte(os.Getenv("JWT_SECRET_KEY"))
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

// middleware for parsing HTTP Token Header from incoming requests
func JwtAuthentication(endpoint http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InfoLog.Println("***** Executing JWT middleware *****")
		// verify if Token header exists
		headers := r.Header
		_, exists := headers["Token"]
		if exists {
			tokenString := r.Header.Get("Token")
			// validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return mySigningKey, nil
			})
			if err != nil {
				ErrorLog.Println(err)
			}
			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	fmt.Println(claims)
			if token.Valid {
				InfoLog.Println("JWT Auth is successful!")
				endpoint.ServeHTTP(w, r)
			} else {
				ErrorLog.Println("JWT Auth Token is NOT valid!")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not Authorised! JWT Auth Token is NOT valid!"))
			}
		} else {
			ErrorLog.Println("JWT Auth Token HTTP Header is NOT Present!")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorised! JWT Auth Token HTTP Header is NOT Present!"))
		}
	})
}

// start server
func StartServer() {

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.Handle("/", JwtAuthentication(homePage))
	mux.Handle("/articles", JwtAuthentication(ViewArticles))
	mux.Handle("/article/create", JwtAuthentication(CreateArticle))
	mux.Handle("/article/delete", JwtAuthentication(DeleteArticle))
	mux.Handle("/article/view", JwtAuthentication(ViewArticle))
	mux.Handle("/article/update", JwtAuthentication(UpdateArticle))
	mux.Handle("/articles/delete_all", JwtAuthentication(DeleteArticles))

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
