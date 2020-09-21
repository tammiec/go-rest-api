package dbint

import (
	"database/sql"
	
	// "github.com/jackc/pgx/v4"
)


// GetDB returns a DB connection
func GetConn(driverName string, driverString string) (*sql.DB, error) {
	return sql.Open(driverName, driverString)
}

// func select(db *sql.DB, tableName string, whereCondition string) ([]interface{}, error) {
// 	return []interface{}, error
// }
