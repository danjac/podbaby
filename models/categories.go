package models

import "database/sql"

// Category is a single feed category e.g. "Arts"
type Category struct {
	ID       int           `db:"id" json:"id"`
	Name     string        `db:"name" json:"name"`
	ParentID sql.NullInt64 `db:"parent_id" json:"parentId"`
}
