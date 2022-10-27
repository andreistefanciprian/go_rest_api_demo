package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	auth "github.com/andreistefanciprian/go_rest_api_demo/frontend/authentication"
	jwt "github.com/golang-jwt/jwt/v4"
)

func (app *application) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    "Dow Jones",
		"expires": time.Now().Add(time.Minute * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(app.mySigningKey)
	if err != nil {
		app.errorLog.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (app *application) render(w http.ResponseWriter, files []string, data interface{}) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error - pars", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error - exec templ", http.StatusInternalServerError)
	}
}

func (app *application) sendApiRequest(w http.ResponseWriter, r *http.Request, JwtToken string, endpointUrl string, book *Article) {
	marshal_struct, _ := json.Marshal(book)
	client := &http.Client{}
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {

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
	validToken, err := app.GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	if r.Method == http.MethodGet {
		client := &http.Client{}
		endpointUrl := fmt.Sprintf("%s/articles", app.backendUrl)
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

		app.render(w, files, allArticles)
	}

	if r.Method == http.MethodPost && r.FormValue("add") == "Add" {
		newBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		endpointUrl := fmt.Sprintf("%s/article/create", app.backendUrl)

		app.sendApiRequest(w, r, validToken, endpointUrl, newBook)
	}

	if r.Method == http.MethodPost && r.FormValue("update") == "Update" {
		id := r.PostFormValue("id")
		endpointUrl := fmt.Sprintf("%s/article/update?id=%s", app.backendUrl, id)

		updatedBook := &Article{
			Title:   r.PostFormValue("title"),
			Desc:    r.PostFormValue("description"),
			Content: r.PostFormValue("content"),
		}

		app.sendApiRequest(w, r, validToken, endpointUrl, updatedBook)

	}

	if r.Method == http.MethodPost && r.FormValue("delete") == "Delete" {
		id := r.PostFormValue("id")
		endpointUrl := fmt.Sprintf("%s/article/delete?id=%s", app.backendUrl, id)

		app.sendApiRequest(w, r, validToken, endpointUrl, nil)
	}
}

var dbCon = auth.UserGorm{}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		rawPassword := r.PostFormValue("password")
		passHash, err := auth.HashPassword(rawPassword)
		if err != nil {
			app.errorLog.Println("Password couldn't be hashed.")
		}
		user := &auth.User{
			Email:          r.PostFormValue("email"),
			HashedPassword: passHash,
		}

		marshal_struct, _ := json.Marshal(user)
		app.infoLog.Println(string(marshal_struct))

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/login.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	app.render(w, files, nil)
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
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
			app.errorLog.Println("Password couldn't be hashed.")
		}
		newUser.HashedPassword = passHash
		if dbCon.Connect(auth.DbConnectionString) {
			_, err := dbCon.CreateUser(newUser)
			if err != nil {
				newUser.Errors = make(map[string]string)
				newUser.Errors["Email"] = "Email address is already registered!"

				app.render(w, files, &newUser)
				return
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}

	app.render(w, files, nil)
}
