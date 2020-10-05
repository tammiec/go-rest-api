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

func GetUsers(db *sql.DB) []*User {
	users := make([]*User, 0)
	rows, err := db.Query("SELECT id, name, email, password FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Row error: %v\n", err)
			os.Exit(1)
		}
		users = append(users, user)
	}
	return users
}

func GetUser(db *sql.DB, id int) *User {
	user := &User{}
	stmt, err := db.Prepare("SELECT id, name, email, password FROM users WHERE id=$1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Row error: %v\n", err)
		os.Exit(1)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Row error: %v\n", err)
		os.Exit(1)
	}
	return user
}

func DeleteUser(db *sql.DB, id int) *User {
	user := &User{}
	stmt, err := db.Prepare("DELETE FROM users WHERE id=$1 RETURNING id, name, email, password")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Row error: %v\n", err)
		os.Exit(1)
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	return user
}
