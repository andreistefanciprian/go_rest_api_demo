## Description

Been learning how to create a simple REST API server in Golang that exposes endpoints to allow accessing and manipulating articles in a library stored in mysql db. 

We're using these Golang packages:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* encoding/json (for transforming data from go struct to json and vice versa)
* jwt-go (for authenticating requests with JWT tokens)

## API Specification

The operations that our endpoint will allow include:
* **Create** a new article in response to a valid POST request at /article/create.
* **Fetch** an article in response to a valid GET request at /article/view?id={id}.
* **Fetch** a list of all articles in response to a valid GET request at /articles.
* **Update** an article in response to a valid PUT request at /article/update?id={id}.
* **Delete** an article in response to a valid DELETE request at /article/delete?id={id}.

## Requirements

* Golang (for writing the application code)
* Docker (for packaging and running the application)
* curl (for testing)
* https://jwt.io/ for generating jwt tokens

## Run app with docker-compose

```
# start app
docker-compose up --build

# remove containers
docker-compose down

# debug commands
docker network ls
docker run --network go_web_api_demo_demo -it --rm --name mysql-client mysql \
mysql -hmy_db -udemouser -pdemopassword -e "SELECT * from quickdemo.articles;"
```

## Test REST API endpoints with curl

```
# generate jwt token from https://jwt.io/
# make sure you use the same secret key (JWT_SECRET_KEY) to generate a valid JWT token
JWT_TOKEN="<paste here jwt token>"

# get all articles
curl -X GET http://localhost:8080/articles \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN" 

# create new article
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

# view article by id
curl -X GET 'http://localhost:8080/article/view?id=32' \
-H 'Content-Type: application/json' \
-H "Token:$JWT_TOKEN"


# update article by id
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

# delete article by id
curl -X DELETE -H "Token:$JWT_TOKEN" 'http://localhost:8080/article/delete?id=32'

# delete all articles
curl -X DELETE -H "Token:$JWT_TOKEN" 'http://localhost:8080/articles/delete_all_'
```