
##### Localhost debug/develop
```
# start only db with mapped db port to localhost
docker-compose up --build

# export env vars
export JWT_SECRET_KEY=your-256-bit-secret
export MYSQL_PORT=3306
export MYSQL_DATABASE=quickdemo
export MYSQL_USER=demouser
export MYSQL_PASSWORD=demopassword
export MYSQL_HOST=localhost

# start backend
go run .
```