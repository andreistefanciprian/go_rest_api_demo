package dbmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbInfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var dbErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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

type Articles []*Article

// Connect establishes connection to mysql
func Connect(DbConnectionString string) {
	Db, err = gorm.Open(mysql.Open(DbConnectionString), &gorm.Config{})
	if err != nil {
		dbErrorLog.Fatal("Failed to connect database", err)
	}
}

// InitialMigration creates the table if it doesn't exist
func InitialMigration(Db *gorm.DB) {
	Db.AutoMigrate(&Article{})
}

func (a *Articles) getArticles(Db *gorm.DB) error {
	result := Db.Find(&a)
	if result.Error != nil {
		return result.Error
	}
	dbInfoLog.Printf("Retrieved %v records from db.", result.RowsAffected)
	return nil
}

func (a *Articles) JSONViewArticles(w io.Writer) error {
	err := a.getArticles(Db)
	if err != nil {
		dbErrorLog.Println("Error retrieving articles from database.")
		return err
	}
	e := json.NewEncoder(w)
	e.Encode(a)
	return nil
}

// delete all articles
func DbDeleteArticles(Db *gorm.DB) {
	var allArticles []Article
	resultFind := Db.Find(&allArticles)
	dbInfoLog.Printf("Retrieved %v records from db.", resultFind.RowsAffected)
	result := Db.Unscoped().Delete(&allArticles) // hard delete
	dbInfoLog.Printf("Deleted %v records from db.", result.RowsAffected)
}

func DeleteArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nEndpoint Hit: delete all articles")
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	DbDeleteArticles(Db)
	dbInfoLog.Println("Deleted all articles from database.")
}

func addArticle(Db *gorm.DB, article Article) error {
	result := Db.Create(&article)
	if result.Error != nil {
		return result.Error
	}
	dbInfoLog.Printf("Added '%s' article in database.", article.Title)
	return nil
}

func (a *Article) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func (a *Article) AddArticle() error {
	err := addArticle(Db, *a)
	if err != nil {
		dbErrorLog.Println("Coudn't add article to db:", err)
	}
	return nil
}

// delete article
func deleteArticle(Db *gorm.DB, id int) error {
	var article Article
	result := Db.Unscoped().Delete(&article, id) // hard delete
	msg := fmt.Sprintf("Deleted %v record from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		dbErrorLog.Println(msg)
		return errors.New(msg)
	} else {
		dbInfoLog.Println(msg)
		return nil
	}
}

func DeleteArticle(id int) error {
	err := deleteArticle(Db, id)
	if err != nil {
		dbErrorLog.Println("Couldn't find article with ID", id)
		return err
	}
	dbInfoLog.Println("Deleted article with ID", id)
	return nil
}

// view article
func (a *Article) getArticle(Db *gorm.DB, id int) error {
	result := Db.First(&a, id)
	msg := fmt.Sprintf("Retrieved %v records from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		dbErrorLog.Println(msg)
		return errors.New(msg)
	} else {
		dbInfoLog.Println(msg)
		return nil
	}
}

func (a *Article) GetArticle(id int) error {
	err := a.getArticle(Db, id)
	if err != nil {
		msg := fmt.Sprintf("Coudn't find article with id %d", id)
		dbErrorLog.Println(msg)
		return errors.New(msg)
	}
	dbInfoLog.Printf("Retrieved article '%s' with id %d.", a.Title, id)
	return nil
}

func (a *Article) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	e.Encode(a)
	return nil
}

// update article
func updateArticle(Db *gorm.DB, article Article, id int) error {
	result := Db.Model(&article).Where("id = ?", id).Updates(article)
	msg := fmt.Sprintf("Updated %v record/s in db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		// dbErrorLog.Println(msg)
		return errors.New("an article with this id doesn't exist")
	} else {
		dbInfoLog.Println(msg)
		return nil
	}
}

func (a *Article) UpdateArticle(id int) error {
	err := updateArticle(Db, *a, id)
	if err != nil {
		dbErrorLog.Printf("Coudn't update article with id %d\n%s", id, err)
		return err
	}
	dbInfoLog.Printf("Updated article with id %d. New book title: '%s'", id, a.Title)
	return nil
}

// middleware for parsing HTTP Token Header from incoming requests
func JwtAuthentication(endpoint http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbInfoLog.Println("Executing JWT middleware")
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
				dbErrorLog.Println(err)
			}
			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	fmt.Println(claims)
			if token.Valid {
				dbInfoLog.Println("JWT Auth is successful!")
				endpoint.ServeHTTP(w, r)
			} else {
				dbErrorLog.Println("JWT Auth Token is NOT valid!")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not Authorised!\nJWT Auth Token is NOT valid!"))
			}
		} else {
			dbErrorLog.Println("JWT Auth Token HTTP Header is NOT Present!")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorised!\nJWT Auth Token HTTP Header is NOT Present!"))
		}
	})
}
