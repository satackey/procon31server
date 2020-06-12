package db

import (
	"database/sql"
	"fmt"
)

func GetDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "procon31server:password@tcp(mysql:3306)/procon31")
	if err != nil {
		return nil, fmt.Errorf("データベースに接続できませんでした: %w", err)
	}
	// defer db.Close()

	return db, nil
}
