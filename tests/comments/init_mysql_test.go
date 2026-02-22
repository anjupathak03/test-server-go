package comments

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// init makes an outgoing MySQL call to verify the database is reachable and
// that the required tables (todos, comments) already exist before any comment
// test is allowed to run.
func init() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "password"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "todo_test_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	log.Printf("[comments/init_mysql] pinging MySQL at %s:%s/%s", host, port, name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("[comments/init_mysql] WARNING: could not open MySQL connection: %v", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Printf("[comments/init_mysql] WARNING: MySQL ping failed: %v", err)
		return
	}

	log.Printf("[comments/init_mysql] MySQL is reachable")

	// Verify required tables exist.
	for _, table := range []string{"todos", "comments"} {
		var tname string
		row := db.QueryRow(
			"SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?",
			name, table,
		)
		if err := row.Scan(&tname); err != nil {
			log.Printf("[comments/init_mysql] WARNING: table %q not found in schema %q: %v", table, name, err)
		} else {
			log.Printf("[comments/init_mysql] table %q confirmed", tname)
		}
	}
}
