# Coffehouse
REST API для управления базой данных кофейни.
## Структура БД
![alt text](database/schema.svg)

## Установка драйверов
```
go get github.com/mattn/go-sqlite3 (требуется С- компилятор (gcc))
go get github.com/gorilla/mux 
```
Запуск программы - ``` go run main.go```<br>
Создание базы - (cmd) ```sqlite3 database/coffee.db < database/create_tables.sql```<br>
Наполнение базы - (cmd) ```sqlite3 database/coffee.db < database/fill_tables.sql```

## Структура проекта
```
Coffehouse  # Корневая папка проекта
├── database # Все, что связано с БД
│   ├── coffee.db # Файл SQLite базы данных
│   ├── create_tables.sql # SQL-скрипт для создания БД
│   ├── fill_tables.sql # SQL-скрипт для наполнения БД
│   ├── schema.svg # Схема БД
│   └── sqlite.go # Подключение к БД
├── go.mod
├── go.sum 
├── handler 
│   └── customers.go # CRUD-обработчик
├── main.go # Точка входа
├── README.md #
└── requests #
    └── Queries.postman_collection.json # Набор запросов для тестов
```