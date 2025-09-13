package databse

import (
	"database/sql"
	"fmt"
)

// Подключение к БД
func Connect(db_path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться: %w", err)
	}

	return db, nil
}