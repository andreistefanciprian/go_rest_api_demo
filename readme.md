## Description

Been learning how to create a simple REST API server in Golang that exposes endpoints to allow accessing and manipulating articles in a library stored in mysql db. 

We're using these Golang packages:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* encoding/json (for transforming data from go struct to json and vice versa)

## API Specification

The operations that our endpoint will allow include:
**Create** a new article in response to a valid POST request at /article/create.
**Fetch** an article in response to a valid GET request at /article/view?id={id}.
**Fetch** a list of all articles in response to a valid GET request at /articles.
**Update** an article in response to a valid PUT request at /article/update?id={id}.
**Delete** an article in response to a valid DELETE request at /article/delete?id={id}.

## Requirements

* Golang
* Docker
* Postman

## Run app

```
# start app and db
docker-compose up --build

# remove containers
docker-compose down

# debug commands
docker network ls
docker run --network go_web_api_demo_demo -it --rm --name mysql-client mysql \
mysql -hmy_db -udemouser -pdemopassword -e "SELECT * from quickdemo.articles;"
```

## Test app endpoints with Postman

```
# get all articles in db
GET http://localhost:8080/articles

# create new article
POST http://localhost:8080/article/create
POST Body:
{
    "Title": "New Book Title",
    "desc": "Book Description",
    "content": "Book Content"
}

# view article by id
GET http://localhost:8080/article/view?id=186

# update article by id
POST http://localhost:8080/article/update?id=186

# delete article by id
DEL http://localhost:8080/article/delete?id=127
```