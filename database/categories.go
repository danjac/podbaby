package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type CategoryReader interface {
	SelectAll(DataHandler) ([]models.Category, error)
	SelectByChannelID(DataHandler, int64) ([]models.Category, error)
}

type CategoryDB interface {
	CategoryReader
}

type CategorySqlDB struct {
	CategoryReader
}

func newCategoryDB() CategoryDB {
	return &CategorySqlDB{
		CategoryReader: &CategorySqlReader{},
	}
}

type CategorySqlReader struct{}

func (r *CategorySqlReader) SelectAll(dh DataHandler) ([]models.Category, error) {
	q := `SELECT id, name, parent_id 
FROM categories 
WHERE id IN (SELECT category_id FROM channels_categories)
ORDER BY name`
	var categories []models.Category
	err := sqlx.Select(dh, &categories, q)
	return categories, dbErr(err, q)
}

func (r *CategorySqlReader) SelectByChannelID(dh DataHandler, channelID int64) ([]models.Category, error) {
	q := `SELECT c.id, c.name, c.parent_id FROM categories c 
JOIN channels_categories cc ON cc.category_id=c.id
WHERE cc.channel_id=$1
GROUP BY c.id
ORDER BY c.name`
	var categories []models.Category
	err := sqlx.Select(dh, &categories, q, channelID)
	return categories, dbErr(err, q)
}
