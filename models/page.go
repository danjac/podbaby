package models

import "math"

const defaultPageSize = 10

func NewPaginator(page int64, numRows int64) *Paginator {
	p := &Paginator{
		CurrentPage: page,
		NumRows:     numRows,
		PageSize:    defaultPageSize,
	}
	p.NumPages = int64(math.Ceil(float64(numRows) / float64(defaultPageSize)))
	p.Offset = (page - 1) * defaultPageSize
	return p
}

type Paginator struct {
	NumRows     int64 `json:"numRows"`
	NumPages    int64 `json:"numPages"`
	CurrentPage int64 `json:"page"`
	PageSize    int64 `json:"pageSize"`
	Offset      int64 `json:"-"`
}
