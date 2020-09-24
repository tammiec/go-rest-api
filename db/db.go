package db

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id int
	Name string
	Email string
	Password string
}

func GetDbPool(dbUrl string) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func GetUsers(dbpool *pgxpool.Pool) []*User {
	users := make([]*User, 0)
	rows, err := dbpool.Query(context.Background(), "SELECT id, name, email, password FROM users")
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

func GetUser(dbpool *pgxpool.Pool, id int) *User {
	user := &User{}
	err := dbpool.QueryRow(context.Background(), fmt.Sprintf("SELECT id, name, email, password FROM users WHERE id=%v", id)).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Row error: %v\n", err)
		os.Exit(1)
	}
	return user
}
