package database

import (
	"fmt"
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CategoryReader interface {
	SelectAll() ([]models.Category, error)
	SelectByChannelID(int64) ([]models.Category, error)
}

type CategoryWriter interface {
	Create(*models.Channel) error
}

type CategoryDB struct {
	CategoryReader
	CategoryWriter
}

func newCategoryDB(db sqlx.Ext) *CategoryDB {
	return &CategoryDB{
		CategoryReader: &CategoryDBReader{db},
		CategoryWriter: &CategoryDBWriter{db},
	}
}

type CategoryDBReader struct {
	sqlx.Ext
}

func (db *CategoryDBReader) SelectAll() ([]models.Category, error) {
	q := `SELECT id, name, parent_id 
FROM categories 
WHERE id IN (SELECT category_id FROM channels_categories)
ORDER BY name`
	var categories []models.Category
	err := sqlx.Select(db, &categories, q)
	return categories, dbErr(err, q)
}

func (db *CategoryDBReader) SelectByChannelID(channelID int64) ([]models.Category, error) {
	q := `SELECT c.id, c.name, c.parent_id FROM categories c 
JOIN channels_categories cc ON cc.category_id=c.id
WHERE cc.channel_id=$1
GROUP BY c.id
ORDER BY c.name`
	var categories []models.Category
	err := sqlx.Select(db, &categories, q, channelID)
	return categories, dbErr(err, q)
}

type CategoryDBWriter struct {
	sqlx.Ext
}

func (db *CategoryDBWriter) Create(channel *models.Channel) error {
	if len(channel.Categories) == 0 {
		return nil
	}
	args := []interface{}{
		channel.ID,
	}

	params := make([]string, 0, len(channel.Categories))
	for i, category := range channel.Categories {
		params = append(params, fmt.Sprintf("$%v", i+2))
		args = append(args, category)
	}

	q := fmt.Sprintf("SELECT add_categories($1, ARRAY[%s])", strings.Join(params, ", "))
	_, err := db.Exec(q, args...)
	return dbErr(err, q)
}
