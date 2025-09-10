package crud

import (
	"database/sql"
	"fmt"
)

// Клиент кофейни
type Customer struct {
	ID         int
	LastName   string
	FirstName  string
	Patronymic string
	Phone      string
}

// Добавление нового клиента
func Create(db *sql.DB, c Customer) error {
	query := `
        INSERT INTO customers (last_name, first_name, patronymic, phone)
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, c.LastName, c.FirstName, c.Patronymic, c.Phone)
	if err != nil {
		return fmt.Errorf("ошибка добавления клиента: %w", err)
	}	
	return nil
}

// Получение всех клиентов
func RetrieveAll(db *sql.DB) ([]Customer, error) {
	query := `
        SELECT id, last_name, first_name, patronymic, phone FROM customers
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выборки клиентов: %w", err)
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var c Customer
		var patronymic sql.NullString

		if err := rows.Scan(&c.ID, &c.LastName, &c.FirstName, &patronymic, &c.Phone); err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	return customers, nil
}

// Получение одного клиента по ID
func Retrieve(db *sql.DB, id int) (*Customer, error) {
	query := `
        SELECT id, last_name, first_name, patronymic, phone FROM customers WHERE id = ?
    `
	var c Customer
	var patronymic sql.NullString

	err := db.QueryRow(query, id).Scan(&c.ID, &c.LastName, &c.FirstName, &patronymic, &c.Phone)
	if err == sql.ErrNoRows {
		return nil, nil // Клиент не найден
	} else if err != nil {
		return nil, fmt.Errorf("ошибка поиска клиента: %w", err)
	}

	return &c, nil
}

// Редактирование данных клиента
func Update(db *sql.DB, c Customer) error {
	query := `
        UPDATE customers 
        SET last_name = ?, first_name = ?, patronymic = ?, phone = ?
        WHERE id = ?
    `
	result, err := db.Exec(query, c.LastName, c.FirstName, c.Patronymic, c.Phone, c.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления клиента: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("клиент с ID %d не найден", c.ID)
	}

	return nil
}

// DELETE — Удаление клиента по ID
func Delete(db *sql.DB, id int) error {
	query := `DELETE FROM customers WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления клиента: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("клиент с ID %d не найден", id)
	}

	return nil
}
