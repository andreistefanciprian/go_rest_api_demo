package main

import (
	"bytes"
	"encoding/json"
	"fmt" // New import
	"html/template"
	"io/ioutil" // New import
	"log"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Article struct {
	Id      int    `json:"ID"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var mySigningKey = []byte("your-256-bit-secret")
var backendUrl = "http://localhost:8080"

// var mySigningKey = []byte(os.Getenv("JWT_SECRET_KEY"))

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

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

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

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/home.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", allArticles)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
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
		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}
		marshal_struct, _ := json.Marshal(newBook)

		client := &http.Client{}
		endpointUrl := fmt.Sprintf("%s/article/create", backendUrl)
		req, _ := http.NewRequest("POST", endpointUrl, bytes.NewBuffer(marshal_struct))
		req.Header.Set("Token", validToken)
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}
		if res.StatusCode == 200 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/addbook.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error - pars", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error - exec templ", 500)
	}
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
		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}
		marshal_struct, _ := json.Marshal(newBook)

		client := &http.Client{}
		endpointUrl := fmt.Sprintf("%s/article/update?id=%s", backendUrl, id)
		req, _ := http.NewRequest("POST", endpointUrl, bytes.NewBuffer(marshal_struct))
		req.Header.Set("Token", validToken)
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		}
		if res.StatusCode == 200 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/updatebook.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error - pars", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error - exec templ", 500)
	}
}

func handleRequests() {
	fmt.Println("Starting server on port 5002 ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/addbook", addBook)
	mux.HandleFunc("/updatebook", updateBook)
	mux.HandleFunc("/", home)
	http.ListenAndServe(":5002", mux)
}

func main() {
	handleRequests()
}
