package tags

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdateTagColor(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	res, _ := testutil.DB.Exec(
"INSERT INTO tags (name, color) VALUES (?, ?)",
"wip", "#FFFFFF",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE tags SET color = ? WHERE id = ?", "#00FF00", id)
	assert.NoError(t, err)

	var color string
	testutil.DB.QueryRow("SELECT color FROM tags WHERE id = ?", id).Scan(&color)
	assert.Equal(t, "#00FF00", color)
}

func TestIntegrationUpdateTagName(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	res, _ := testutil.DB.Exec(
"INSERT INTO tags (name, color) VALUES (?, ?)",
"old-name", "#AAAAAA",
)
	id, _ := res.LastInsertId()

	_, err := testutil.DB.Exec("UPDATE tags SET name = ? WHERE id = ?", "new-name", id)
	assert.NoError(t, err)

	var name string
	testutil.DB.QueryRow("SELECT name FROM tags WHERE id = ?", id).Scan(&name)
	assert.Equal(t, "new-name", name)
}

func TestIntegrationDeleteTag(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	res, _ := testutil.DB.Exec(
"INSERT INTO tags (name, color) VALUES (?, ?)",
"temp-tag", "#999999",
)
	id, _ := res.LastInsertId()

	result, err := testutil.DB.Exec("DELETE FROM tags WHERE id = ?", id)
	assert.NoError(t, err)

	affected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), affected)

	var count int
	testutil.DB.QueryRow("SELECT COUNT(*) FROM tags WHERE id = ?", id).Scan(&count)
	assert.Equal(t, 0, count)
}
