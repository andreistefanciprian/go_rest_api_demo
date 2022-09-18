## Description

Been learning how to create a simple WEB API using these Golang packages:
* net/http (standard library for using web servers)
* gorm (mysql ORM)
* viper (for reading db credentials from .env file)
* encoding/json (for transforming data from go struct to json and vice versa)

## Create mysql container
```
# create container with mysql db
docker container run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql

# connect to mysql db and create a database and credentials
db_address=`docker container inspect --format '{{ .NetworkSettings.IPAddress }}' mysql`
docker run -it --rm --name mysql-client mysql mysql -h${db_address} -uroot -pmy-secret-pw

CREATE DATABASE quickdemo;
CREATE USER 'demouser'@'%' IDENTIFIED BY 'demopassword';
GRANT ALL PRIVILEGES ON quickdemo.* TO 'demouser'@'%';
FLUSH PRIVILEGES;

# debug commands
docker run -it --rm --name mysql-client mysql mysql -h${db_address} -uroot -pmy-secret-pw -e "SHOW DATABASES;"
docker run -it --rm --name mysql-client mysql mysql -h${db_address} -uroot -pmy-secret-pw -e "DESCRIBE quickdemo.articles;"
docker run -it --rm --name mysql-client mysql mysql -h${db_address} -uroot -pmy-secret-pw -e "SELECT * from quickdemo.articles;"
```

## Run app

```
go mod init
go mod tidy
go run .
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