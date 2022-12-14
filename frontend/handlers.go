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
		app.errorLog.Println("Couldn't reach endpoint", endpointUrl)
		fmt.Fprintf(w, "Error: %s", err.Error())
	}
	if res.StatusCode == 200 {
		app.infoLog.Println("Successfully reached endpoint", endpointUrl)
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

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./templates/base.tmpl",
		"./templates/pages/login.tmpl",
		"./templates/partials/nav.tmpl",
		"./templates/partials/footer.tmpl",
	}

	if r.Method == http.MethodPost {
		user := &auth.User{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		// verify user exists
		registeredUser, err := app.users.ByEmail(user.Email)
		// if user doesn't exist, raise popup warning
		if err != nil {
			user.Errors = make(map[string]string)
			user.Errors["PopUp"] = "Email address is not registered!"
			app.render(w, files, &user)
			return
		}
		// if user exists check if the password hash matches password in db
		if !auth.CheckPasswordHash(user.Password, registeredUser.HashedPassword) {
			user.Errors = make(map[string]string)
			user.Errors["PopUp"] = "Password doesn't match!"
			app.render(w, files, &user)
			return
		} else {
			// redirect use to home page
			app.infoLog.Println(user.Email, "Successful Login.")
			http.Redirect(w, r, "/"+"?login="+registeredUser.FirstName, http.StatusSeeOther)
		}

		// if password hash matched records, generate JWT Token
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
		_, err = app.users.CreateUser(newUser)
		if err != nil {
			newUser.Errors = make(map[string]string)
			newUser.Errors["PopUp"] = "Email address is already registered!"

			app.render(w, files, &newUser)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	app.render(w, files, nil)
}
