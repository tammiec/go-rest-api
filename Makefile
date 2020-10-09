run:
	DATABASE_URL=postgresql://tammiechung@localhost:5432/python_project?sslmode=disable go run main.go
.PHONY: run

test:
	go test -coverprofile=cover.out ./...
.PHONY: test
