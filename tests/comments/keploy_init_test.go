package comments

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"test-server/testutil"
)

// init is the canonical package-level init for the comments test package.
// It calls StartKeploySession first, then performs connectivity pre-checks.
func init() {
	testutil.StartKeploySession("tests/comments", "keploy")

	// MySQL connectivity check (consolidated from init_mysql_test.go).
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
	} else {
		defer db.Close()
		if err = db.Ping(); err != nil {
			log.Printf("[comments/init_mysql] WARNING: MySQL ping failed: %v", err)
		} else {
			log.Printf("[comments/init_mysql] MySQL is reachable")
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
	}

	// HTTP connectivity check (consolidated from init_http_test.go).
	const url = "https://jsonplaceholder.typicode.com/comments/1"

	httpClient := &http.Client{Timeout: 5 * time.Second}

	log.Printf("[comments/init_http] fetching sample comment from %s", url)

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("[comments/init_http] WARNING: HTTP request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[comments/init_http] JSONPlaceholder responded with HTTP %d", resp.StatusCode)
}
