package users

import (
	"testing"

	"test-server/testutil"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateUser(t *testing.T) {
	testutil.StartKeploySession("tests/users", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	result, err := testutil.DB.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		"alice", "alice@example.com",
	)
	assert.NoError(t, err)

	id, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.NotZero(t, id)

	var username, email string
	err = testutil.DB.QueryRow("SELECT username, email FROM users WHERE id = ?", id).
		Scan(&username, &email)
	assert.NoError(t, err)
	assert.Equal(t, "alice", username)
	assert.Equal(t, "alice@example.com", email)
}

func TestIntegrationCreateUserDuplicate(t *testing.T) {
	testutil.StartKeploySession("tests/users", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	_, err := testutil.DB.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		"bob", "bob@example.com",
	)
	assert.NoError(t, err)

	_, err = testutil.DB.Exec(
		"INSERT INTO users (username, email) VALUES (?, ?)",
		"bob", "bob2@example.com",
	)
	assert.Error(t, err)
}

func TestIntegrationCreateMultipleUsers(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	names := []string{"user1", "user2", "user3"}
	for _, name := range names {
		_, err := testutil.DB.Exec(
			"INSERT INTO users (username, email) VALUES (?, ?)",
			name, name+"@example.com",
		)
		assert.NoError(t, err)
	}

	var count int
	err := testutil.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}
