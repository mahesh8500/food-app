To run it locally:  
go mod tidy

Run postgres with docker:  
docker run --name postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=food-app -p 5432:5432 -d postgres:16

Run code:   
go run cmd/server/main.go

Run tests:  
go test ./handlers -v

Run via docker:     
go build -o food-app .     
docker build -t food-app:latest .
docker run --rm -p 8080:8080 \
-e PG_CONN_STR="postgres://user:secret@host.docker.internal:5432/food-app?sslmode=disable" \
food-app:latest


