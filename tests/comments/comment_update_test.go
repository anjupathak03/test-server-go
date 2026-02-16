package comments

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdateCommentBody(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	res, _ := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
todoID, "original text", "editor",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE comments SET body = ? WHERE id = ?", "updated text", id)
	assert.NoError(t, err)

	var body string
	testutil.DB.QueryRow("SELECT body FROM comments WHERE id = ?", id).Scan(&body)
	assert.Equal(t, "updated text", body)
}

func TestIntegrationUpdateCommentAuthor(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	res, _ := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
todoID, "some note", "old_author",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE comments SET author = ? WHERE id = ?", "new_author", id)
	assert.NoError(t, err)

	var author string
	testutil.DB.QueryRow("SELECT author FROM comments WHERE id = ?", id).Scan(&author)
	assert.Equal(t, "new_author", author)
}

func TestIntegrationDeleteComment(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "comments")
	testutil.CleanTable(t, "todos")

	todoID := seedTodo(t)

	res, _ := testutil.DB.Exec(
"INSERT INTO comments (todo_id, body, author) VALUES (?, ?, ?)",
todoID, "to be deleted", "cleaner",
)
	id, _ := res.LastInsertId()

	result, err := testutil.DB.Exec("DELETE FROM comments WHERE id = ?", id)
	assert.NoError(t, err)

	affected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), affected)

	var count int
	testutil.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE id = ?", id).Scan(&count)
	assert.Equal(t, 0, count)
}
