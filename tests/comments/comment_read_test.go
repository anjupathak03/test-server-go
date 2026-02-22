package comments

import (
	"testing"

	"test-server/testutil"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationReadCommentByID(t *testing.T) {
	testutil.StartKeploySession("tests/comments", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	res, _ := testutil.DB.Exec(
		"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
		todoID, "Needs review", "alice",
	)
	id, _ := res.LastInsertId()

	var body, author string
	var tid int
	err := testutil.DB.QueryRow(
		"SELECT todo_id, body, author FROM comments WHERE id = ?", id,
	).Scan(&tid, &body, &author)
	assert.NoError(t, err)
	assert.Equal(t, int(todoID), tid)
	assert.Equal(t, "Needs review", body)
	assert.Equal(t, "alice", author)
}

func TestIntegrationReadCommentsByTodo(t *testing.T) {
	testutil.StartKeploySession("tests/comments", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	testutil.DB.Exec("INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)", todoID, "c1", "a1")
	testutil.DB.Exec("INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)", todoID, "c2", "a2")

	rows, err := testutil.DB.Query(
		"SELECT id, body, author FROM comments WHERE todo_id = ? ORDER BY id", todoID,
	)
	assert.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		var id int
		var body, author string
		assert.NoError(t, rows.Scan(&id, &body, &author))
		count++
	}
	assert.Equal(t, 2, count)
}

func TestIntegrationReadNonExistentComment(t *testing.T) {
	testutil.StartKeploySession("tests/comments", t.Name())
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")

	var body string
	err := testutil.DB.QueryRow("SELECT body FROM comments WHERE id = ?", 99999).Scan(&body)
	assert.Error(t, err)
}
