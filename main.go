package main

import (
	"fmt"
	"os"
)

func main() {
	// connect to db
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// migrate db
	initialMigration()

	// insert articles
	for i := 0; i < 5; i++ {
		createArticle(Article{Title: "Book Title", Desc: "Book Description", Content: "Book Content"})
	}

	startServer()
}
