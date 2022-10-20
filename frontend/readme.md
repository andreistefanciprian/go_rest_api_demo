
##### Localhost debug/develop
```
# start only backend with mapped backend port to localhost
docker-compose up --build

# export env vars
export JWT_SECRET_KEY="your-256-bit-secret"
export REST_API_HOST="localhost"
export REST_API_PORT="8080"

# start frontend
go run main.go

# start frontend container
docker image build -t frontend:latest .
docker container run -e JWT_SECRET_KEY=your-256-bit-secret -e REST_API_HOST=backend -e REST_API_PORT=8080 --network go_web_api_demo_demo --name frontend -p 8090:8090 frontend:latest
docker container run \
-e JWT_SECRET_KEY="your-256-bit-secret" \
-e REST_API_HOST="backend" \
-e REST_API_PORT="8080" \
--network go_web_api_demo_demo \
--name frontend \
-p 8090:8090 frontend:latest
```
