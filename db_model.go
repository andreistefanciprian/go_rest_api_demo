package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConnectionString string

type Article struct {
	gorm.Model
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func initialMigration() {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	db.AutoMigrate(&Article{})
}

func getArticles() []Article {
	var allArticles []Article
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	result := db.Find(&allArticles)
	// fmt.Printf("Retrieved %v records from db.", result.RowsAffected)
	fmt.Printf("Retrieved %v records from db.", result.RowsAffected)
	return allArticles
}

func createArticle(article Article) {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	db.Create(&article) // pass pointer of data to Create
}

func deleteArticle(id int) {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	db.Unscoped().Delete(&Article{}, id) // hard delete
}

func getArticle(id int) Article {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	var article Article
	db.First(&article, id)
	return article
	// fmt.Println(article)
}

func updateArticle(article Article, id int) {
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	db.Model(&article).Where("id = ?", id).Updates(article)
}
