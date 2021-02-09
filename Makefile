get:
	go get -t -v
	go get github.com/joho/godotenv/cmd/godotenv # This lets you run godotenv from the cli
.PHONY: get

build-mocks:
	PATH=$(PATH):$(HOME)/go/bin mockgen -source=handlers/users/users.go -destination=generatedmocks/handlers/users/users.go
	PATH=$(PATH):$(HOME)/go/bin mockgen -source=services/users/users.go -destination=generatedmocks/services/users/users.go
	PATH=$(PATH):$(HOME)/go/bin mockgen -source=clients/render/render.go -destination=generatedmocks/clients/render/render.go
	PATH=$(PATH):$(HOME)/go/bin mockgen -source=dals/users/users.go -destination=generatedmocks/dals/users/users.go
.PHONY: build-mocks

run:
	PATH=$(PATH):$(HOME)/go/bin godotenv -f .env go run .
.PHONY: run

test:
	go test -coverprofile=cover.out ./...
.PHONY: test

docker-deps-start:
	COMPOSE_PROJECT_NAME=test docker-compose up -d
.PHONY: docker-deps-start
