package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// get all articles
func allArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: articles")
	books := getArticles()
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// add article
func addArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: create")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	createArticle(article)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// delete article
func removeArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: delete article")
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	deleteArticle(id)
	fmt.Fprintf(w, "Deleted article with ID %d", id)
}

// view article
func viewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: view article")
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	article := getArticle(id)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// update article
func changeArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: update article")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	updateArticle(article, id)

	fmt.Fprintf(w, "Updated article with ID %d", id)
}

// start server
func startServer() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", allArticles)
	http.HandleFunc("/article/create", addArticle)
	http.HandleFunc("/article/delete", removeArticle)
	http.HandleFunc("/article/view", viewArticle)
	http.HandleFunc("/article/update", changeArticle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
