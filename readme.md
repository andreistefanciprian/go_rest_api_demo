## Description

Been learning how to create a simple WEB API using these Golang packages:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* encoding/json (for transforming data from go struct to json and vice versa)

## Run app

```
# start app and db
docker-compose up --build

# remove containers
docker-compose down

# debug commands
docker network ls
docker run --network go_web_api_demo_demo \
-it --rm --name mysql-client mysql \
mysql -hmy_db -udemouser -pdemopassword \
-e "SELECT * from quickdemo.articles;"
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
POST http://localhost:8080/article/view?id=186

# delete article by id
DEL http://localhost:8080/article/delete?id=127
```