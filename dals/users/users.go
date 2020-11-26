package model

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func GetDb(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return db
}

func GetUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
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

func GetUser(db *sql.DB, id int) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("SELECT id, name, email FROM users WHERE id=$1")
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

func DeleteUser(db *sql.DB, id int) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("DELETE FROM users WHERE id=$1 RETURNING id, name, email")
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

func CreateUser(db *sql.DB, name string, email string, password string) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email")
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

func UpdateUser(db *sql.DB, id int, name string, email string, password string) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4 RETURNING id, name, email")
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
