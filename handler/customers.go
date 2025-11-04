package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Клиент кофейни
type Customer struct {
	ID         int    `json:"id"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	Patronymic string `json:"patronymic"`
	Phone      string `json:"phone"`
}

type CustomerHandler struct {
	DB *sql.DB
}

// Добавление нового клиента
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := `
        INSERT INTO customers (last_name, first_name, patronymic, phone)
        VALUES (?, ?, ?, ?)
    `
	result, err := h.DB.Exec(query, c.LastName, c.FirstName, c.Patronymic, c.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var id, _ = result.LastInsertId()
	c.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)

}

// Получение всех клиентов
func (h *CustomerHandler) RetrieveAll(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	sortByAttribute := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	filterAttribute := r.URL.Query().Get("attribute")
	filterValue := r.URL.Query().Get("value")

	query := `
        SELECT id, last_name, first_name, patronymic, phone FROM customers
    `
	var args []interface{}
	var conditions []string

	if filterAttribute != "" && filterValue != "" {
		conditions = append(conditions, filterAttribute+"= ?")
		args = append(args, filterValue)
	}

	if search != "" {
		like := "%" + search + "%"
		conditions = append(conditions, "lower(last_name) LIKE ?  OR lower(first_name) LIKE ?  OR lower(patronymic) LIKE ?")

		args = append(args, like, like, like)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	if sortByAttribute != "" {
		if order != "desk" {
			order = "asc"
		}
		query += " ORDER BY " + sortByAttribute + " " + order
	}
	print("query: ", query)
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := make([]Customer, 0)

	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.LastName, &c.FirstName, &c.Patronymic, &c.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)

}

// Получение одного клиента по ID
func (h *CustomerHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := `
        SELECT id, last_name, first_name, patronymic, phone FROM customers WHERE id = ?
    `
	var c Customer

	err := h.DB.QueryRow(query, id).Scan(&c.ID, &c.LastName, &c.FirstName, &c.Patronymic, &c.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// Редактирование данных клиента
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `
        UPDATE customers 
        SET last_name = ?, first_name = ?, patronymic = ?, phone = ?
        WHERE id = ?
    `
	result, err := h.DB.Exec(query, c.LastName, c.FirstName, c.Patronymic, c.Phone, id)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}
	c.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// Удаление клиента по ID
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	str_id := mux.Vars(r)["id"]
	id, err := strconv.Atoi(str_id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM customers WHERE id = ?`
	result, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
