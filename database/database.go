package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var DB *sql.DB

func InitDB(config Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func CreateTodoTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %w", err)
	}

	log.Println("Todos table created or already exists")
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
