package tags

import (
"testing"

"test-server/testutil"

"github.com/stretchr/testify/assert"
)

func TestIntegrationReadTagByID(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	res, _ := testutil.DB.Exec(
"INSERT INTO tags (name, color) VALUES (?, ?)",
"feature", "#0000FF",
)
	id, _ := res.LastInsertId()

	var name, color string
	err := testutil.DB.QueryRow("SELECT name, color FROM tags WHERE id = ?", id).
		Scan(&name, &color)
	assert.NoError(t, err)
	assert.Equal(t, "feature", name)
	assert.Equal(t, "#0000FF", color)
}

func TestIntegrationReadAllTags(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "t1", "#111111")
	testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "t2", "#222222")
	testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "t3", "#333333")

	rows, err := testutil.DB.Query("SELECT id, name, color FROM tags ORDER BY id")
	assert.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		var id int
		var name, color string
		assert.NoError(t, rows.Scan(&id, &name, &color))
		count++
	}
	assert.Equal(t, 3, count)
}

func TestIntegrationReadTagByName(t *testing.T) {
	testutil.SetupDB(t)
	defer testutil.TeardownDB()
	testutil.CleanTable(t, "tags")

	testutil.DB.Exec("INSERT INTO tags (name, color) VALUES (?, ?)", "docs", "#AABBCC")

	var id int
	var color string
	err := testutil.DB.QueryRow("SELECT id, color FROM tags WHERE name = ?", "docs").
		Scan(&id, &color)
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.Equal(t, "#AABBCC", color)
}
