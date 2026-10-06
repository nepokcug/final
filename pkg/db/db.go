package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// schema — SQL-код для создания таблицы и индекса
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_date ON scheduler(date);
`

// DB — глобальная переменная для доступа к БД
var DB *sql.DB

// Init открывает БД и создаёт таблицу, если её нет
func Init(dbFile string) error {
	// 1. Проверяем dbFile
	if err := validateDBFile(dbFile); err != nil {
		return err
	}

	// 2. Создаём папку, если её нет
	dir := filepath.Dir(dbFile)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// 3. Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// 4. Открываем БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// 5. Проверяем соединение
	if err := DB.Ping(); err != nil {
		return err
	}

	// 6. Создаём таблицу и индекс, если файла не было
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateDBFile проверяет корректность пути к файлу БД
func validateDBFile(dbFile string) error {
	// Пустая строка
	if dbFile == "" {
		return errors.New("dbFile не может быть пустым")
	}

	// Точка или две точки (папки)
	if dbFile == "." || dbFile == ".." {
		return errors.New("dbFile не может быть папкой")
	}

	// Заканчивается на / или \ (папка)
	if strings.HasSuffix(dbFile, "/") || strings.HasSuffix(dbFile, "\\") {
		return errors.New("dbFile не может заканчиваться на разделитель")
	}

	// Не заканчивается на .db
	if !strings.HasSuffix(dbFile, ".db") {
		return errors.New("dbFile должен заканчиваться на .db")
	}

	return nil
}
