
##### Localhost debug/develop
```
# start only backend with mapped backend port to localhost
docker-compose up --build

# export env vars
export JWT_SECRET_KEY="your-256-bit-secret"
export REST_API_HOST="localhost"
export REST_API_PORT="8080"
export MYSQL_PORT=3306
export MYSQL_DATABASE=quickdemo
export MYSQL_USER=demouser
export MYSQL_PASSWORD=demopassword
export MYSQL_HOST=localhost

# start frontend
go run main.go

# start frontend container
docker image build -t frontend:latest .

docker container run \
-e JWT_SECRET_KEY="your-256-bit-secret" \
-e REST_API_HOST="backend" \
-e REST_API_PORT="8080" \
-e MYSQL_PORT=3306 \
-e MYSQL_DATABASE=quickdemo \
-e MYSQL_USER=demouser \
-e MYSQL_PASSWORD=demopassword \
-e MYSQL_HOST=my_db \
--network go_rest_api_demo_demo \
--name frontend \
-p 8090:8090 frontend:latest
```
