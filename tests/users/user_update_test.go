package users

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdateUserEmail(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	res, _ := testutil.DB.Exec(
"INSERT INTO users (username, email) VALUES (?, ?)",
"dave", "dave@old.com",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE users SET email = ? WHERE id = ?", "dave@new.com", id)
	assert.NoError(t, err)

	var email string
	testutil.DB.QueryRow("SELECT email FROM users WHERE id = ?", id).Scan(&email)
	assert.Equal(t, "dave@new.com", email)
}

func TestIntegrationDeactivateUser(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	res, _ := testutil.DB.Exec(
"INSERT INTO users (username, email) VALUES (?, ?)",
"eve", "eve@example.com",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE users SET active = ? WHERE id = ?", false, id)
	assert.NoError(t, err)

	var active bool
	testutil.DB.QueryRow("SELECT active FROM users WHERE id = ?", id).Scan(&active)
	assert.False(t, active)
}

func TestIntegrationDeleteUser(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "users")

	res, _ := testutil.DB.Exec(
"INSERT INTO users (username, email) VALUES (?, ?)",
"frank", "frank@example.com",
)
	id, _ := res.LastInsertId()

	result, err := testutil.DB.Exec("DELETE FROM users WHERE id = ?", id)
	assert.NoError(t, err)

	affected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), affected)

	var count int
	testutil.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	assert.Equal(t, 0, count)
}
