package model

import (
	"os"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

type User struct {
	Id int
	Name string
	Email string
	Password string
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
	users := make([]*User, 0)
	rows, err := db.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

func GetUser(db *sql.DB, id int) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("SELECT id, name, email, password FROM users WHERE id=$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, err
}

func DeleteUser(db *sql.DB, id int) (*User, error) {
	user := &User{}
	stmt, err := db.Prepare("DELETE FROM users WHERE id=$1 RETURNING id, name, email, password")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, err
}
