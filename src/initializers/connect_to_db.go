package initializers

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectToDb() (*sql.DB, error) {
	dataConnection := os.Getenv("DB")

	psqlInfo := fmt.Sprintf(dataConnection)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error when pinging database: %w", err)
	}

	return db, nil
}
