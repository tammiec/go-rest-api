# Basic Golang REST API

This is a very basic REST API written in Golang and runs on a local PostgreSQL database. 

### Setup
```sh
# Clone this repository on your local machine
git clone git@github.com:tammiec/go-rest-api.git
cd go-rest-api
# Install Dependencies
make get
```

### Run the app
```
make run
```

### Run the tests
```
make test
```

## Endpoints

### Check to see if server is ready
* `GET /readiness`
#### Request
```
curl -i http://localhost:8000/readiness
```

#### Response
```sh
HTTP/1.1 200 OK
Date: Tue, 10 Nov 2020 21:47:22 GMT
Content-Length: 5
Content-Type: text/plain; charset=utf-8

Ready
```

### Get a list of all users
* `GET /users`
#### Request
```
curl -i http://localhost:8000/users
```

#### Response
```sh
HTTP/1.1 200 OK
Date: Tue, 10 Nov 2020 21:49:01 GMT
Content-Length: 377
Content-Type: text/plain; charset=utf-8

[
    {"Id":1,"Name":"Peter","Email":"peter@mail.com"},
    {"Id":2,"Name":"John","Email":"john@mail.com"}
]
```

### Create a new user
* `POST /users`
#### Request
Required params:
* name `string`
* email `string`
* password `string`
```
curl -i --data "name=John&email=j@mail.com&password=password" http://localhost:8000/users
```

#### Response
```sh
HTTP/1.1 200 OK
Date: Tue, 10 Nov 2020 21:58:02 GMT
Content-Length: 44
Content-Type: text/plain; charset=utf-8

{"Id":3,"Name":"John","Email":"j@mail.com"}
```

### Get a specific user
* `GET /user/{userID}`
#### Request
```
curl -i http://localhost:8000/users/1
```

#### Response
```sh
HTTP/1.1 200 OK
Date: Tue, 10 Nov 2020 21:52:13 GMT
Content-Length: 49
Content-Type: text/plain; charset=utf-8

{"Id":1,"Name":"Peter","Email":"peter@mail.com"}
```

### TODO
- deal with missing params on create and update endpoints
- migrate tests and use generated mocks
- set up dockerfile and docker container for PSQL and schema to test DAL
- handle errors
