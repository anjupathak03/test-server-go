package tags

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func TestIntegrationCreateTag(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	result, err := testutil.DB.Exec(
"INSERT INTO tags (name, color) VALUES (?, ?)",
"urgent", "#FF0000",
)
	assert.NoError(t, err)

	id, _ := result.LastInsertId()
	assert.NotZero(t, id)

	var name, color string
	err = testutil.DB.QueryRow("SELECT name, color FROM tags WHERE id = ?", id).
		Scan(&name, &color)
	assert.NoError(t, err)
	assert.Equal(t, "urgent", name)
	assert.Equal(t, "#FF0000", color)
}

func TestIntegrationCreateTagDefaultColor(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	_, err := testutil.DB.Exec("INSERT INTO tags (name) VALUES (?)", "low-priority")
	assert.NoError(t, err)

	var color string
	testutil.DB.QueryRow("SELECT color FROM tags WHERE name = ?", "low-priority").Scan(&color)
	assert.Equal(t, "#000000", color)
}

func TestIntegrationCreateDuplicateTag(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	_, err := testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "bug", "#FF0000")
	assert.NoError(t, err)

	_, err = testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "bug", "#00FF00")
	assert.Error(t, err)
}
