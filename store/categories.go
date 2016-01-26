package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type CategoryReader interface {
	SelectAll(DataHandler) ([]models.Category, error)
	SelectByChannelID(DataHandler, int64) ([]models.Category, error)
}

type CategoryStore interface {
	CategoryReader
}

type categorySqlStore struct {
	CategoryReader
}

func newCategoryStore() CategoryStore {
	return &categorySqlStore{
		CategoryReader: &categorySqlReader{},
	}
}

type categorySqlReader struct{}

func (r *categorySqlReader) SelectAll(dh DataHandler) ([]models.Category, error) {
	q := `
    SELECT id, name, parent_id 
    FROM categories 
    WHERE id IN (SELECT category_id FROM channels_categories)
    ORDER BY name`
	var categories []models.Category
	err := sqlx.Select(dh, &categories, q)
	return categories, err
}

func (r *categorySqlReader) SelectByChannelID(dh DataHandler, channelID int64) ([]models.Category, error) {
	q := `
    SELECT c.id, c.name, c.parent_id FROM categories c 
    JOIN channels_categories cc ON cc.category_id=c.id
    WHERE cc.channel_id=$1
    GROUP BY c.id
    ORDER BY c.name`
	var categories []models.Category
	err := sqlx.Select(dh, &categories, q, channelID)
	return categories, err
}
