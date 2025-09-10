package main

import (
	crud "coffe/CRUD"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func connect(db_path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться: %w", err)
	}

	return db, nil

}

func main() {
	db, err := connect("./db/coffee.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var choice int = 1
	for choice != 0 {
		fmt.Println(`
			0. Выход 
			1. Добавить клиента
			2. Получить данные всех клиентов
			3. Получить данные клиента по идентификатору
			4. Изменить данные клиента
			5. Удалить данные клиента	

			`)
		fmt.Print("Ваш выбор: ")
		fmt.Scanln(&choice)

		switch choice {
		case 0:
			fmt.Println("Выход из программы.")
		case 1:
			fmt.Println("Добавление клиента.")
			var customer crud.Customer

			fmt.Print("Фамилия: ")
			fmt.Scanln(&customer.LastName)

			fmt.Print("Имя: ")
			fmt.Scanln(&customer.FirstName)

			fmt.Print("Отчество: ")
			fmt.Scanln(&customer.Patronymic)

			fmt.Print("Номер телефона: ")
			fmt.Scanln(&customer.Phone)

			if err := crud.Create(db, customer); err != nil {
				log.Println("Ошибка: ", err)
			} else {
				fmt.Println("Запись успешно добавлена")
			}
		case 2:
			fmt.Println("Список данные всех клиентов.")
			customers, err := crud.RetrieveAll(db)

			if err != nil {
				log.Println(err)
			}

			for _, c := range customers {
				fmt.Printf("%d: %s %s %s | %s\n", c.ID, c.LastName, c.FirstName, c.Patronymic, c.Phone)
			}
		case 3:
			fmt.Println("Данные клиента по идентификатору.")
		case 4:
			fmt.Println("Изменить данные клиента.")
		case 5:
			fmt.Println("Удаление клиента.")
		default:
			fmt.Println("Некорректный ввод!")

		}
	}

}
