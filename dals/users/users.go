package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	model "github.com/tammiec/go-rest-api/models/user"
)

type Deps struct{}

type Config struct {
	Url string
}

type Users interface {
	GetName() string
	Ping(ctx context.Context) error
	GetUsers() ([]*model.UserResponse, error)
	GetUser(id int) (*model.UserResponse, error)
	CreateUser(name string, email string, password string) (*model.UserResponse, error)
	UpdateUser(id int, name string, email string, password string) (*model.UserResponse, error)
	DeleteUser(id int) (*model.UserResponse, error)
}

type impl struct {
	db *sql.DB
}

func New(deps *Deps, config *Config) Users {
	db, err := sql.Open("postgres", config.Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &impl{db: db}
}

func (impl *impl) GetName() string {
	return "Users"
}

func (impl *impl) Ping(ctx context.Context) error {
	return impl.db.PingContext(ctx)
}

func (impl *impl) GetUsers() ([]*model.UserResponse, error) {
	rows, err := impl.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.UserResponse, 0)
	for rows.Next() {
		user := &model.UserResponse{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) < 1 {
		return nil, errors.New("no users found")
	}

	return users, err
}

func (impl *impl) GetUser(id int) (*model.UserResponse, error) {
	user := &model.UserResponse{}
	stmt, err := impl.db.Prepare("SELECT id, name, email FROM users WHERE id=$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (impl *impl) DeleteUser(id int) (*model.UserResponse, error) {
	user := &model.UserResponse{}
	stmt, err := impl.db.Prepare("DELETE FROM users WHERE id=$1 RETURNING id, name, email")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (impl *impl) CreateUser(name string, email string, password string) (*model.UserResponse, error) {
	user := &model.UserResponse{}
	stmt, err := impl.db.Prepare("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name, email, password).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (impl *impl) UpdateUser(id int, name string, email string, password string) (*model.UserResponse, error) {
	user := &model.UserResponse{}
	stmt, err := impl.db.Prepare("UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4 RETURNING id, name, email")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name, email, password, id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, err
}
