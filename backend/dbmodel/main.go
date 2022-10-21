package dbmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbInfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var dbErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var DbConnectionString string
var Db *gorm.DB
var err error

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
func (a *Articles) deleteArticles(Db *gorm.DB) error {
	err := a.getArticles(Db)
	if err != nil {
		dbErrorLog.Println("Error retrieving articles from database.")
		return err
	}
	result := Db.Unscoped().Delete(&a)
	if result.Error != nil {
		dbErrorLog.Println("Error deleting articles from database.")
		return err
	} else {
		dbInfoLog.Printf("Deleted %v records from db.", result.RowsAffected)
		return nil
	}
}

func (a *Articles) DeleteArticles() error {
	err := a.deleteArticles(Db)
	if err != nil {
		dbErrorLog.Println(err.Error())
		return errors.New("aaaa")
	} else {
		return nil
	}
}

// create article
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
