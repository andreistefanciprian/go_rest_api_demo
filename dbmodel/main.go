package dbmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbConnectionString string
var Db *gorm.DB
var err error
var mySigningKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Article struct holds the data table in the db
type Article struct {
	gorm.Model
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Connect establishes connection to mysql
func Connect(DbConnectionString string) {
	Db, err = gorm.Open(mysql.Open(DbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
}

// InitialMigration creates the table if it doesn't exist
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
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
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

// middleware for parsing HTTP Token Header from incoming requests
func JwtAuthentication(endpoint http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middleware")
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
				fmt.Println(err)
			}
			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	fmt.Println(claims)
			if token.Valid {
				log.Print("JWT Auth is successful!")
				endpoint.ServeHTTP(w, r)
			} else {
				log.Print("JWT Auth Token is NOT valid!")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not Authorised!\nJWT Auth Token is NOT valid!"))
			}
		} else {
			log.Print("JWT Auth Token HTTP Header is NOT Present!")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorised!\nJWT Auth Token HTTP Header is NOT Present!"))
		}
	})
}
