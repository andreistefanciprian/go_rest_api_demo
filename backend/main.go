package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andreistefanciprian/go_rest_api_demo/backend/dbmodel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	mySigningKey []byte
	articles     *dbmodel.UserModel
}

// establishes connection to mysql
func connectDB(dbConnectionString string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		// errorLog.Fatal("Failed to connect database", err)
		return nil, err
	}
	// infoLog.Println("Successfully connected to db.")
	return db, nil
}

func main() {
	// define vars
	httpPort := ":8080"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	// connect to db + migrate db
	db, err := connectDB(dbConnectionString)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		mySigningKey: mySigningKey,
		articles:     &dbmodel.UserModel{DB: db},
	}

	app.articles.InitialMigration()

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.Handle("/", app.JwtAuthentication(app.homePage))
	mux.Handle("/articles", app.JwtAuthentication(app.ViewArticles))
	mux.Handle("/article/create", app.JwtAuthentication(app.CreateArticle))
	mux.Handle("/article/delete", app.JwtAuthentication(app.DeleteArticle))
	mux.Handle("/article/view", app.JwtAuthentication(app.ViewArticle))
	mux.Handle("/article/update", app.JwtAuthentication(app.UpdateArticle))
	mux.Handle("/articles/delete_all", app.JwtAuthentication(app.DeleteArticles))

	// create a new server
	srv := http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// start the server
	app.infoLog.Println("Starting server on port", httpPort)
	err = srv.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
