package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

// CategoryReader handles reads from category data store
type CategoryReader interface {
	SelectAll(DataHandler, *[]models.Category) error
	SelectByChannelID(DataHandler, *[]models.Category, int) error
}

// CategoryStore manages interactions with category data store
type CategoryStore interface {
	CategoryReader
}

type categorySQLStore struct {
	CategoryReader
}

func newCategoryStore() CategoryStore {
	return &categorySQLStore{
		CategoryReader: &categorySQLReader{},
	}
}

type categorySQLReader struct{}

func (r *categorySQLReader) SelectAll(dh DataHandler, categories *[]models.Category) error {
	q := `
    SELECT id, name, parent_id 
    FROM categories 
    WHERE id IN (SELECT category_id FROM channels_categories)
    ORDER BY name`
	return handleError(sqlx.Select(dh, categories, q), q)
}

func (r *categorySQLReader) SelectByChannelID(dh DataHandler, categories *[]models.Category, channelID int) error {
	q := `
    SELECT c.id, c.name, c.parent_id FROM categories c 
    JOIN channels_categories cc ON cc.category_id=c.id
    WHERE cc.channel_id=$1
    GROUP BY c.id
    ORDER BY c.name`
	return handleError(sqlx.Select(dh, categories, q, channelID), q)
}
