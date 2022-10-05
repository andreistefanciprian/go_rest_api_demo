package dbmodel

import (
	"encoding/json"
	"errors"
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
func DbViewArticles(Db *gorm.DB) []Article {
	var allArticles []Article
	result := Db.Find(&allArticles)
	fmt.Printf("Retrieved %v records from db.", result.RowsAffected)
	return allArticles
}

func ViewArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: articles")
	books := DbViewArticles(Db)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// delete all articles
func DbDeleteArticles(Db *gorm.DB) {
	var allArticles []Article
	resultFind := Db.Find(&allArticles)
	fmt.Printf("Retrieved %v records from db.", resultFind.RowsAffected)
	result := Db.Unscoped().Delete(&allArticles) // hard delete
	msg := fmt.Sprintf("Deleted %v records from db.", result.RowsAffected)
	fmt.Println(msg)
}

func DeleteArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: delete all articles")
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	DbDeleteArticles(Db)
	fmt.Fprintf(w, "Deleted all articles.")
}

// add article
func DbCreateArticle(Db *gorm.DB, article Article) {
	Db.Create(&article) // pass pointer of data to Create
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
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

	DbCreateArticle(Db, article)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// delete article
func DbDeleteArticle(Db *gorm.DB, id int) error {
	var article Article
	result := Db.Unscoped().Delete(&article, id) // hard delete
	msg := fmt.Sprintf("Deleted %v records from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		fmt.Println(msg)
		return errors.New(msg)
	} else {
		fmt.Println(msg)
		return nil
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
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
	result := DbDeleteArticle(Db, id)
	if result == nil {
		fmt.Fprintf(w, "Deleted article with ID %d", id)
	} else {
		fmt.Fprintf(w, "%s\nCouldn't find article with ID %d", result, id)
	}
}

// view article
func DbViewArticle(Db *gorm.DB, id int) (Article, error) {
	var article Article
	result := Db.First(&article, id)
	msg := fmt.Sprintf("Retrieved %v records from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		fmt.Println(msg)
		return Article{}, errors.New(msg)
	} else {
		fmt.Println(msg)
		return article, nil
	}
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

	article, err := DbViewArticle(Db, id)
	if err != nil {
		fmt.Fprintf(w, "%s\nCouldn't find article with ID %d", err, id)
	} else {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

// update article
func dbUpdateArticle(Db *gorm.DB, article Article, id int) error {
	result := Db.Model(&article).Where("id = ?", id).Updates(article)
	msg := fmt.Sprintf("Updated %v records from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		fmt.Println(msg)
		return errors.New(msg)
	} else {
		fmt.Println(msg)
		return nil
	}
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
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

	result := dbUpdateArticle(Db, article, id)
	if result != nil {
		fmt.Fprintf(w, "%s\nCouldn't find article with ID %d", result, id)
	} else {
		fmt.Fprintf(w, "Updated article with ID %d", id)
	}
}
