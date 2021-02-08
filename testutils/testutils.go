package testutils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func GetDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://root:test@0.0.0.0:6603/test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func RunSQL(db *sql.DB, scriptPath string) {
	file, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		panic(err)
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		fmt.Println("REQUEST", request)
		if len(strings.TrimSpace(request)) > 0 {
			_, err := db.Query("request")
			if err != nil {
				panic(err)
			}
		}
	}
}

func RunSQLFile(db *sql.DB, filename string) {
	dir, _ := os.Getwd()
	RunSQL(db, filepath.Join(dir, `../../`, `sql`, filename))
}
