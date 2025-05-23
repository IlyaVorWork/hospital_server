package postgres

import "database/sql"
import _ "github.com/lib/pq"

func OpenPostgresDB(connURL string) (*sql.DB, func(db *sql.DB), error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	return db, closePostgresDB, nil
}

func closePostgresDB(db *sql.DB) {
	_ = db.Close()
}
