package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/andreistefanciprian/go_rest_api_demo/frontend/authentication"
)

type Article struct {
	Id      int    `json:"ID"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	mySigningKey []byte
	backendUrl   string
}

func main() {
	httpPort := ":8090"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	backendUrl := fmt.Sprintf("http://%s:%s", os.Getenv("REST_API_HOST"), os.Getenv("REST_API_PORT"))

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		backendUrl:   backendUrl,
		mySigningKey: mySigningKey,
	}

	// connect to db + migrate db
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	auth.DbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	var db = &auth.UserGorm{}

	db.Connect(auth.DbConnectionString)
	db.InitialMigration()

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/login", app.login)
	mux.HandleFunc("/register", app.register)
	mux.HandleFunc("/", app.home)

	// create a new server
	srv := http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// start the server
	app.infoLog.Println("Starting server on port", httpPort)
	err := srv.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}

}
