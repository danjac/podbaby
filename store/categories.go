package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type CategoryReader interface {
	SelectAll(DataHandler, *[]models.Category) error
	SelectByChannelID(DataHandler, *[]models.Category, int) error
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

func (r *categorySqlReader) SelectAll(dh DataHandler, categories *[]models.Category) error {
	q := `
    SELECT id, name, parent_id 
    FROM categories 
    WHERE id IN (SELECT category_id FROM channels_categories)
    ORDER BY name`
	return sqlx.Select(dh, categories, q)
}

func (r *categorySqlReader) SelectByChannelID(dh DataHandler, categories *[]models.Category, channelID int) error {
	q := `
    SELECT c.id, c.name, c.parent_id FROM categories c 
    JOIN channels_categories cc ON cc.category_id=c.id
    WHERE cc.channel_id=$1
    GROUP BY c.id
    ORDER BY c.name`
	return sqlx.Select(dh, categories, q, channelID)
}
