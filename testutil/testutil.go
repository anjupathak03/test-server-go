package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the shared test database connection.
var DB *sql.DB

// SetupDB initialises the test database and creates all required tables.
func SetupDB(t *testing.T) {
	t.Helper()
	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "3306")
	user := envOr("DB_USER", "root")
	pass := envOr("DB_PASSWORD", "password")
	name := envOr("DB_NAME", "todo_test_db")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err = DB.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	mustExec(t, DB, `CREATE TABLE IF NOT EXISTS todos (
id INT AUTO_INCREMENT PRIMARY KEY,
title VARCHAR(255) NOT NULL,
description TEXT,
completed BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)

	mustExec(t, DB, `CREATE TABLE IF NOT EXISTS users (
id INT AUTO_INCREMENT PRIMARY KEY,
username VARCHAR(100) NOT NULL UNIQUE,
email VARCHAR(255) NOT NULL UNIQUE,
active BOOLEAN DEFAULT TRUE,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)

	mustExec(t, DB, `CREATE TABLE IF NOT EXISTS tags (
id INT AUTO_INCREMENT PRIMARY KEY,
name VARCHAR(100) NOT NULL UNIQUE,
color VARCHAR(7) DEFAULT '#000000',
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)

	mustExec(t, DB, `CREATE TABLE IF NOT EXISTS comments (
id INT AUTO_INCREMENT PRIMARY KEY,
todo_id INT NOT NULL,
body TEXT NOT NULL,
author VARCHAR(100) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE
)`)
}

// CleanTable truncates the given table.
func CleanTable(t *testing.T, table string) {
	t.Helper()
	DB.Exec("SET FOREIGN_KEY_CHECKS=0")
	DB.Exec("TRUNCATE TABLE " + table)
	DB.Exec("SET FOREIGN_KEY_CHECKS=1")
}

// TeardownDB closes the shared connection.
func TeardownDB() {
	if DB != nil {
		DB.Close()
		DB = nil
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustExec(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec failed: %v\nquery: %s", err, query)
	}
}

// StartKeploySession notifies the Keploy agent to begin a new sandbox session.
// The call is skipped silently when KEPLOY_AGENT_URI (or KEPLOY_MCP_URI) is unset.
func StartKeploySession(location, name string) {
	if location == "" || name == "" {
		return
	}
	agentURI := os.Getenv("KEPLOY_AGENT_URI")
	if agentURI == "" {
		agentURI = os.Getenv("KEPLOY_MCP_URI")
	}
	if agentURI == "" {
		log.Println("Keploy session skipped: KEPLOY_AGENT_URI not set")
		return
	}
	endpoint := strings.TrimRight(agentURI, "/") + "/sandbox/scope"
	payload, _ := json.Marshal(map[string]string{"location": location, "name": name})
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("Keploy session start failed (non-fatal): %v", err)
		return
	}
	resp.Body.Close()
}
