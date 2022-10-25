package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	auth "github.com/andreistefanciprian/go_rest_api_demo/frontend/authentication"
	jwt "github.com/golang-jwt/jwt/v4"

	_ "github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/opentelemetry-go-contrib/launcher"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Article struct {
	Id      int    `json:"ID"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var httpPort = ":8090"
var mySigningKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var backendUrl = fmt.Sprintf("http://%s:%s", os.Getenv("REST_API_HOST"), os.Getenv("REST_API_PORT"))

func GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    "Dow Jones",
		"expires": time.Now().Add(time.Minute * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func render(w http.ResponseWriter, files []string, data interface{}) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error - pars", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error - exec templ", http.StatusInternalServerError)
	}
}

func sendApiRequest(w http.ResponseWriter, r *http.Request, JwtToken string, endpointUrl string, book *Article) {
	marshal_struct, _ := json.Marshal(book)
	client := &http.Client{}
	// endpointUrl := fmt.Sprintf("%s%s", backendUrl, endpooint)
	req, _ := http.NewRequest("POST", endpointUrl, bytes.NewBuffer(marshal_struct))
	req.Header.Set("Token", JwtToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}
	if res.StatusCode == 200 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// html go templated files
	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/home.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}
	// generate JWT token
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	if r.Method == http.MethodGet {
		client := &http.Client{}
		endpointUrl := fmt.Sprintf("%s/articles", backendUrl)
		req, _ := http.NewRequest("GET", endpointUrl, nil)
		req.Header.Set("Token", validToken)
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		var allArticles []Article

		jsonErr := json.Unmarshal([]byte(string(body)), &allArticles)

		if jsonErr != nil {
			fmt.Println(err)
			fmt.Println("ERROR Unmarshaling of JSON failed.")
		}

		render(w, files, allArticles)
	}

	if r.Method == http.MethodPost && r.FormValue("add") == "Add" {
		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		endpointUrl := fmt.Sprintf("%s/article/create", backendUrl)

		sendApiRequest(w, r, validToken, endpointUrl, newBook)
	}

	if r.Method == http.MethodPost && r.FormValue("update") == "Update" {
		id := r.PostFormValue("id")
		endpointUrl := fmt.Sprintf("%s/article/update?id=%s", backendUrl, id)

		updatedBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		sendApiRequest(w, r, validToken, endpointUrl, updatedBook)

	}

	if r.Method == http.MethodPost && r.FormValue("delete") == "Delete" {
		id := r.PostFormValue("id")
		endpointUrl := fmt.Sprintf("%s/article/delete?id=%s", backendUrl, id)

		sendApiRequest(w, r, validToken, endpointUrl, nil)
	}
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addbook" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		validToken, err := GenerateJWT()
		if err != nil {
			fmt.Println("Failed to generate token")
		}

		endpointUrl := fmt.Sprintf("%s/article/create", backendUrl)

		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		sendApiRequest(w, r, validToken, endpointUrl, newBook)
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/addbook.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	render(w, files, nil)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/updatebook" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		validToken, err := GenerateJWT()
		if err != nil {
			fmt.Println("Failed to generate token")
		}
		id := r.PostFormValue("id")
		endpointUrl := fmt.Sprintf("%s/article/update?id=%s", backendUrl, id)

		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		sendApiRequest(w, r, validToken, endpointUrl, newBook)
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/updatebook.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	render(w, files, nil)
}

var dbCon = auth.UserGorm{}

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		rawPassword := r.PostFormValue("password")
		passHash, err := auth.HashPassword(rawPassword)
		if err != nil {
			errorLog.Println("Password couldn't be hashed.")
		}
		user := &auth.User{
			Email:          r.PostFormValue("email"),
			HashedPassword: passHash,
		}

		marshal_struct, _ := json.Marshal(user)
		infoLog.Println(string(marshal_struct))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/login.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	render(w, files, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/register.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	if r.Method == http.MethodPost {
		newUser := &auth.User{
			FirstName: r.PostFormValue("firstname"),
			LastName:  r.PostFormValue("lastname"),
			Email:     r.PostFormValue("email"),
			Password:  r.PostFormValue("password"),
		}
		passHash, err := auth.HashPassword(newUser.Password)
		if err != nil {
			errorLog.Println("Password couldn't be hashed.")
		}
		newUser.HashedPassword = passHash
		if dbCon.Connect(auth.DbConnectionString) {
			_, err := dbCon.CreateUser(newUser)
			if err != nil {
				newUser.Errors = make(map[string]string)
				newUser.Errors["Email"] = "Email address is already registered!"

				render(w, files, &newUser)
				return
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}

	render(w, files, nil)
}

func main() {
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

	// use honeycomb distro to setup OpenTelemetry SDK
	otelShutdown, errr := launcher.ConfigureOpenTelemetry()
	if errr != nil {
		log.Fatalf("error setting up OTel SDK - %e", errr)
	}
	defer otelShutdown()

	// create a new serve mux and register the handlers
	mux := http.NewServeMux()

	loginHandler := http.HandlerFunc(login)
	wrappedLoginHandler := otelhttp.NewHandler(loginHandler, "login")
	mux.Handle("/login", wrappedLoginHandler)

	registerHandler := http.HandlerFunc(register)
	wrappedRegisterHandler := otelhttp.NewHandler(registerHandler, "register")
	mux.Handle("/register", wrappedRegisterHandler)

	homeHandler := http.HandlerFunc(home)
	wrappedHomeHandler := otelhttp.NewHandler(homeHandler, "home")
	mux.Handle("/", wrappedHomeHandler)

	// mux.HandleFunc("/login", login)
	// mux.HandleFunc("/register", register)
	mux.HandleFunc("/addbook", addBook)
	mux.HandleFunc("/updatebook", updateBook)
	// mux.HandleFunc("/", home)

	// create a new server
	srv := http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// start the server
	infoLog.Println("Starting server on port", httpPort)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}
