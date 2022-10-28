package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreistefanciprian/go_rest_api_demo/backend/dbmodel"
	jwt "github.com/golang-jwt/jwt/v4"
)

// home page
func (app *application) homePage(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Book Library"))
}

func (app *application) ViewArticles(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /articles")
	var a dbmodel.Articles
	w.Header().Set("content-type", "application/json")
	err := app.articles.JSONViewArticles(w, &a)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (app *application) CreateArticle(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /article/create")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	article := &dbmodel.Article{}
	article.FromJSON(r.Body)

	err := app.articles.AddArticle(article)
	if err != nil {
		http.Error(w, "Coudn't add article to db.", http.StatusInternalServerError)
	}
}

func (app *application) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /article/update")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		app.errorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	article := &dbmodel.Article{}
	// reqBody, _ := ioutil.ReadAll(r.Body)
	article.FromJSON(r.Body)

	err = app.articles.UpdateArticle(article, id)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Article was updated successfully."))
	}
}

func (app *application) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /article/delete")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		app.errorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	err = app.articles.DeleteArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Article was deleted successfully."))
	}

}

func (app *application) ViewArticle(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /article/view")
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	rawId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id < 1 {
		msg := fmt.Sprintf("Id '%s' is not  a valid id number!", rawId)
		app.errorLog.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	article := &dbmodel.Article{}
	article, err = app.articles.GetArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		w.Header().Set("content-type", "application/json")
		article.ToJSON(w)
	}
}

func (app *application) DeleteArticles(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Endpoint Hit: /articles/delete_all")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var articles dbmodel.Articles
	err := app.articles.DeleteArticles(&articles)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		msg := fmt.Sprintf("Deleted %d articles.", len(articles))
		w.Write([]byte(msg))
	}

}

// middleware for parsing HTTP Token Header from incoming requests
func (app *application) JwtAuthentication(endpoint http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				return app.mySigningKey, nil
			})
			if err != nil {
				app.errorLog.Println(err)
			}
			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	fmt.Println(claims)
			if token.Valid {
				endpoint.ServeHTTP(w, r)
			} else {
				app.errorLog.Println("JWT Auth Token is NOT valid!")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not Authorised! JWT Auth Token is NOT valid!"))
			}
		} else {
			app.errorLog.Println("JWT Auth Token HTTP Header is NOT Present!")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorised! JWT Auth Token HTTP Header is NOT Present!"))
		}
	})
}
