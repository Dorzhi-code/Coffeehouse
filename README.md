# Golang
Предметная область - кофейня. База данных состоит 4 сущностей: Customers (Клиенты), Employees (Сотрудники), Products (Продукты), Orders (Заказы)

Установка драйвера go get github.com/mattn/go-sqlite3 (требуется С- компилятор (gcc))
Запуск программы - go run main.go
Создание базы - (cmd) sqlite3.exe coffee.db < create_tables.sql
Наполнение базы - (cmd) sqlite3.exe coffee.db < fill_tables.sql