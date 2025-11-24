package internal

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLite driver
)

const DBPath = "db/chefops.db"

// OpenDB opens the SQLite database with foreign keys enabled.
func OpenDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", DBPath)
	return sql.Open("sqlite", dsn)
}
