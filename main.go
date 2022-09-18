package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	// read db credentials from.env
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	dbUser := viper.Get("MYSQL_USER")
	dbPassword := viper.Get("MYSQL_PASSWORD")
	dbHost := viper.Get("MYSQL_HOST")
	dbPort := viper.Get("MYSQL_PORT")
	dbName := viper.Get("MYSQL_DB_NAME")
	dbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// migrate db
	initialMigration()

	// // delete articles
	// for i := 1; i < 1000; i++ {
	// 	deleteArticle(i)
	// }

	// insert articles
	for i := 0; i < 5; i++ {
		createArticle(Article{Title: "Book Title", Desc: "Book Description", Content: "Book Content"})
	}

	// // update article
	// newArticle := Article{Title: "Updated Book Title 1", Desc: "Book Description 1", Content: "Book Content 1"}
	// updateArticle(newArticle, 36)

	// // delete article
	// deleteArticle(38)

	// getArticles()
	getArticle(168)
	startServer()

}
