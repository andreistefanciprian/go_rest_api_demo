## Description

Been learning how to create a simple REST API server in Golang that exposes endpoints to allow accessing and manipulating articles in a library stored in mysql db. 

We're using these Golang packages:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* encoding/json (for transforming data from go struct to json and vice versa)

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

## Run app with docker-compose

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

## Test app endpoints with curl

```
# get all articles in db
curl -X GET -H 'Content-Type: application/json' http://localhost:8080/articles

# create new article
curl -X POST http://localhost:8080/article/create \
-H 'Content-Type: application/json' \
-d '''
{
    "Title": "Book Title",
    "desc": "Book Description",
    "content": "Book Content"
}
'''

# view article by id
curl -X GET -H 'Content-Type: application/json' 'http://localhost:8080/article/view?id=32'

# update article by id
curl -X POST 'http://localhost:8080/article/update?id=32' \
-H 'Content-Type: application/json' \
-d '''
{
    "Title": "Updated Book Title",
    "desc": "Book Description",
    "content": "Book Content"
}
'''

# delete article by id
curl -X DELETE 'http://localhost:8080/article/delete?id=32'

# delete all articles
curl -X DELETE 'http://localhost:8080/articles/delete_all_'
```