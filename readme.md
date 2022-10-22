## Description

Easy REST API backend and Bootstrap frontend written in Golang using standard libraries.
The backend exposes CRUD (Create, Read, Update and Delete) endpoints to allow accessing and manipulating articles in a library stored in mysql db. 

![Alt text](/img.png "Frontend Screenshot")

## API Specification (backend)

The operations that our endpoint will allow include:
* **Create** a new article in response to a valid POST request at /article/create.
* **Fetch** an article in response to a valid GET request at /article/view?id={id}.
* **Fetch** a list of all articles in response to a valid GET request at /articles.
* **Update** an article in response to a valid PUT request at /article/update?id={id}.
* **Delete** an article in response to a valid DELETE request at /article/delete?id={id}.

Used Go libraries:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* encoding/json (for transforming data from go struct to json and vice versa)
* jwt-go (for authenticating requests with JWT tokens)

## Frontend

I'm using https://getbootstrap.com/ to have a nice user interface without going deep on CSS, JS and HTML.

Used Go libraries:
* net/http
* html/template

## Requirements

* Golang (for writing the application code)
* Docker (for packaging and running the application)
* Any browser
* https://jwt.io/ for generating jwt tokens

## Run app with docker-compose

```
# start app
docker-compose up --build

# stop app
docker-compose down

# access frontend at localhost:8090/
```

## (Optional) Test REST API endpoints with curl or Postman

###### Generate jwt token with https://jwt.io/
Make sure you use the same secret key (JWT_SECRET_KEY) to generate a valid JWT token.
```
export JWT_TOKEN="<paste here jwt token>"
```

###### Get all articles
```
curl -X GET http://localhost:8080/articles \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN" 
```

###### Create new article
```
curl -X POST http://localhost:8080/article/create \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN" \
-d '''
{
    "Title": "Book Title",
    "desc": "Book Description",
    "content": "Book Content"
}
'''
```

###### View article by id
```
curl -X GET 'http://localhost:8080/article/view?id=32' \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN"
```

###### Update article by id
```
curl -X POST 'http://localhost:8080/article/update?id=32' \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN"
-d '''
{
    "Title": "Updated Book Title",
    "desc": "Book Description",
    "content": "Book Content"
}
'''
```

###### Delete article by id
```
curl -X POST 'http://localhost:8080/article/delete?id=32' \
-H "Token:$JWT_TOKEN"
```

###### Delete all articles
```
curl -X POST 'http://localhost:8080/articles/delete_all' \
-H "Token:$JWT_TOKEN" 
```