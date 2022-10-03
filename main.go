package main

import (
	"fmt"
	"os"

	"github.com/andreistefanciprian/go_web_api_demo/dbmodel"
	"github.com/andreistefanciprian/go_web_api_demo/webserver"
)

func main() {
	// connect to db
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbmodel.DbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// migrate db
	dbmodel.Connect(dbmodel.DbConnectionString)
	dbmodel.InitialMigration(dbmodel.Db)

	// insert articles
	for i := 0; i < 5; i++ {
		dbmodel.CreateArticle(dbmodel.Db, dbmodel.Article{Title: "Book Title", Desc: "Book Description", Content: "Book Content"})
	}

	webserver.StartServer()
}
