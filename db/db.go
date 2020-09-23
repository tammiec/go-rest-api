package db

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetDbPool(dbUrl string) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}

func GetUsers(dbpool *pgxpool.Pool) string {
	var name string
	err := dbpool.QueryRow(context.Background(), "SELECT name FROM users WHERE id=5").Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return name
}
