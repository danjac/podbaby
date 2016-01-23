package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type CategoryReader interface {
	SelectAll() ([]models.Category, error)
	SelectByChannelID(int64) ([]models.Category, error)
}

type CategoryDB struct {
	CategoryReader
}

func newCategoryDB(db *sqlx.DB) *CategoryDB {
	return &CategoryDB{
		CategoryReader: &CategoryDBReader{db},
	}
}

type CategoryDBReader struct {
	*sqlx.DB
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
