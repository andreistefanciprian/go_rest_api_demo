package dbmodel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var DbConnectionString string
var Db *gorm.DB
var err error

func Connect(DbConnectionString string) {
	Db, err = gorm.Open(mysql.Open(DbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
}

func InitialMigration(Db *gorm.DB) {
	Db.AutoMigrate(&Article{})
}

// get all articles
func getArticles(Db *gorm.DB) []Article {
	var allArticles []Article
	result := Db.Find(&allArticles)
	fmt.Printf("Retrieved %v records from db.", result.RowsAffected)
	return allArticles
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: articles")
	books := getArticles(Db)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// add article
func CreateArticle(Db *gorm.DB, article Article) {
	Db.Create(&article) // pass pointer of data to Create
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
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

	CreateArticle(Db, article)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// delete article
func deleteArticle(Db *gorm.DB, id int) {
	Db.Unscoped().Delete(&Article{}, id) // hard delete
}

func RemoveArticle(w http.ResponseWriter, r *http.Request) {
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

	deleteArticle(Db, id)
	fmt.Fprintf(w, "Deleted article with ID %d", id)
}

// view article
func getArticle(Db *gorm.DB, id int) Article {
	var article Article
	Db.First(&article, id)
	return article
}

func ViewArticle(w http.ResponseWriter, r *http.Request) {
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

	article := getArticle(Db, id)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// update article
func updateArticle(Db *gorm.DB, article Article, id int) {
	Db.Model(&article).Where("id = ?", id).Updates(article)
}

func ChangeArticle(w http.ResponseWriter, r *http.Request) {
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
	updateArticle(Db, article, id)

	fmt.Fprintf(w, "Updated article with ID %d", id)
}
