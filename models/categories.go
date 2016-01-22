package models

import "database/sql"

type Category struct {
	ID       int64         `db:"id" json:"id"`
	Name     string        `db:"name" json:"name"`
	ParentID sql.NullInt64 `db:"parent_id" json:"parentId"`
}
