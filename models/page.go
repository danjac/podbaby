package models

import "math"

const defaultPageSize = 10

// NewPaginator creates a new Paginator instance with precalculated fields
func NewPaginator(page int, numRows int) *Paginator {
	p := &Paginator{
		CurrentPage: page,
		NumRows:     numRows,
		PageSize:    defaultPageSize,
	}
	p.NumPages = int(math.Ceil(float64(numRows) / float64(defaultPageSize)))
	p.Offset = (page - 1) * defaultPageSize
	return p
}

// Paginator contains row/page/offset etc for a page instance
type Paginator struct {
	NumRows     int `json:"numRows"`
	NumPages    int `json:"numPages"`
	CurrentPage int `json:"page"`
	PageSize    int `json:"pageSize"`
	Offset      int `json:"-"`
}
