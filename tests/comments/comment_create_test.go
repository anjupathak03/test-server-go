package comments

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func seedTodo(t *testing.T) int64 {
	t.Helper()
	res, err := testutil.DB.Exec(
"INSERT INTO todos (title, description) VALUES (?, ?)",
"Seed Todo", "Created for comment tests",
)
	assert.NoError(t, err)
	id, _ := res.LastInsertId()
	return id
}

func TestIntegrationCreateComment(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	result, err := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
todoID, "Looks good!", "reviewer1",
)
	assert.NoError(t, err)

	id, _ := result.LastInsertId()
	assert.NotZero(t, id)

	var body, author string
	err = testutil.DB.QueryRow("SELECT body, author FROM comments WHERE id = ?", id).
		Scan(&body, &author)
	assert.NoError(t, err)
	assert.Equal(t, "Looks good!", body)
	assert.Equal(t, "reviewer1", author)
}

func TestIntegrationCreateMultipleComments(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	bodies := []string{"First", "Second", "Third"}
	for _, b := range bodies {
		_, err := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
todoID, b, "author",
)
		assert.NoError(t, err)
	}

	var count int
	err := testutil.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE todo_id = ?", todoID).
		Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestIntegrationCreateCommentInvalidTodo(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	_, err := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
99999, "orphan comment", "nobody",
)
	assert.Error(t, err)
}
