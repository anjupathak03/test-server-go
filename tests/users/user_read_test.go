package users

import (
	"testing"

	"test-server/testutil"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationReadUserByID(t *testing.T) {
	testutil.StartKeploySession("tests/users", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	res, _ := testutil.DB.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		"charlie", "charlie@example.com",
	)
	id, _ := res.LastInsertId()

	var username, email string
	var active bool
	err := testutil.DB.QueryRow(
		"SELECT username, email, active FROM users WHERE id = ?", id,
	).Scan(&username, &email, &active)

	assert.NoError(t, err)
	assert.Equal(t, "charlie", username)
	assert.Equal(t, "charlie@example.com", email)
	assert.True(t, active)
}

func TestIntegrationReadAllUsers(t *testing.T) {
	testutil.StartKeploySession("tests/users", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	testutil.DB.Exec("INSERT INTO users (username, email) VALUES (?, ?)", "u1", "u1@x.com")
	testutil.DB.Exec("INSERT INTO users (username, email) VALUES (?, ?)", "u2", "u2@x.com")

	rows, err := testutil.DB.Query("SELECT id, username FROM users ORDER BY id")
	assert.NoError(t, err)
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		var name string
		assert.NoError(t, rows.Scan(&id, &name))
		ids = append(ids, id)
	}
	assert.Len(t, ids, 2)
}

func TestIntegrationReadNonExistentUser(t *testing.T) {
	testutil.StartKeploySession("tests/users", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	var username string
	err := testutil.DB.QueryRow("SELECT username FROM users WHERE id = ?", 99999).
		Scan(&username)
	assert.Error(t, err)
}
