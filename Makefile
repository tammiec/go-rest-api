get:
	go get -t -v
.PHONY: get

run:
	HTTP_HOST=localhost HTTP_PORT=8000 DATABASE_URL=postgresql://tammiechung@localhost:5432/python_project?sslmode=disable go run main.go
.PHONY: run

test:
	go test -coverprofile=cover.out ./...
.PHONY: test
