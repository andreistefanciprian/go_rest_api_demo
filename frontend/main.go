package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/andreistefanciprian/go_rest_api_demo/frontend/authentication"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	users        *auth.UserModel
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
	httpPort := ":8090"
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	backendUrl := fmt.Sprintf("http://%s:%s", os.Getenv("REST_API_HOST"), os.Getenv("REST_API_PORT"))

	// connect to db + migrate db
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := connectDB(dbConnectionString)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		backendUrl:   backendUrl,
		mySigningKey: mySigningKey,
		users:        &auth.UserModel{DB: db},
	}

	app.users.InitialMigration()

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
	err = srv.ListenAndServe()
	if err != nil {
		app.errorLog.Fatal(err)
	}

}
