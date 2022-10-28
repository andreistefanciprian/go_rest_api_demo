package dbmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"gorm.io/gorm"
)

var dbInfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var dbErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

type UserModel struct {
	DB *gorm.DB
}

// Article struct holds the data table in the db
type Article struct {
	gorm.Model
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []*Article

// InitialMigration creates the table if it doesn't exist
func (ug *UserModel) InitialMigration() {
	ug.DB.AutoMigrate(&Article{})
}

func (ug *UserModel) getArticles(a *Articles) error {
	result := ug.DB.Find(&a)
	if result.Error != nil {
		return result.Error
	}
	dbInfoLog.Printf("Retrieved %v records from db.", result.RowsAffected)
	return nil
}

func (ug *UserModel) JSONViewArticles(w io.Writer, a *Articles) error {
	err := ug.getArticles(a)
	if err != nil {
		dbErrorLog.Println("Error retrieving articles from database.")
		return err
	}
	e := json.NewEncoder(w)
	e.Encode(a)
	return nil
}

// delete all articles
func (ug *UserModel) deleteArticles(a *Articles) error {
	err := ug.getArticles(a)
	if err != nil {
		dbErrorLog.Println("Could not retrieve articles from database.")
		return err
	}
	result := ug.DB.Unscoped().Delete(&a)
	if result.Error != nil {
		return result.Error
	} else {
		dbInfoLog.Printf("Deleted %v records from db.", result.RowsAffected)
		return nil
	}
}

func (ug *UserModel) DeleteArticles(a *Articles) error {
	err := ug.deleteArticles(a)
	if err != nil {
		return errors.New("could not delete articles")
	} else {
		return nil
	}
}

// create article
func (ug *UserModel) addArticle(a *Article) error {
	result := ug.DB.Create(&a)
	if result.Error != nil {
		return errors.New("could not create article")
	}
	dbInfoLog.Printf("Added '%s' article in database.", a.Title)
	return nil
}

func (a *Article) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func (ug *UserModel) AddArticle(a *Article) error {
	err := ug.addArticle(a)
	if err != nil {
		return err
	}
	return nil
}

// delete article
func (ug *UserModel) deleteArticle(id int) error {
	var article Article
	result := ug.DB.Unscoped().Delete(&article, id) // hard delete
	msg := fmt.Sprintf("Deleted %v record from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		dbErrorLog.Println(msg)
		return errors.New(msg)
	} else {
		dbInfoLog.Println(msg)
		return nil
	}
}

func (ug *UserModel) DeleteArticle(id int) error {
	err := ug.deleteArticle(id)
	if err != nil {
		dbErrorLog.Println("Couldn't find article with ID", id)
		return err
	}
	dbInfoLog.Println("Deleted article with ID", id)
	return nil
}

// view article
func (ug *UserModel) getArticle(id int) (*Article, error) {
	var a Article
	result := ug.DB.First(&a, id)
	msg := fmt.Sprintf("Retrieved %v records from db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		dbErrorLog.Println(msg)
		return &a, errors.New(msg)
	} else {
		dbInfoLog.Println(msg)
		return &a, nil
	}
}

func (ug *UserModel) GetArticle(id int) (*Article, error) {
	a, err := ug.getArticle(id)
	if err != nil {
		msg := fmt.Sprintf("Coudn't find article with id %d", id)
		dbErrorLog.Println(msg)
		return a, errors.New(msg)
	}
	dbInfoLog.Printf("Retrieved article '%s' with id %d.", a.Title, id)
	return a, nil
}

func (a *Article) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	e.Encode(a)
	return nil
}

// update article
func (ug *UserModel) updateArticle(a *Article, id int) error {
	result := ug.DB.Model(&a).Where("id = ?", id).Updates(a)
	msg := fmt.Sprintf("Updated %v record/s in db.", result.RowsAffected)
	if result.RowsAffected == 0 {
		// dbErrorLog.Println(msg)
		return errors.New("an article with this id doesn't exist")
	} else {
		dbInfoLog.Println(msg)
		return nil
	}
}

func (ug *UserModel) UpdateArticle(a *Article, id int) error {
	err := ug.updateArticle(a, id)
	if err != nil {
		// dbErrorLog.Printf("Coudn't update article with id %d: %s", id, err)
		return err
	}
	dbInfoLog.Printf("Updated article with id %d. New book title: '%s'", id, a.Title)
	return nil
}
