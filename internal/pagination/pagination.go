package pagination

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	} else if p.PerPage > 100 {
		p.PerPage = 100
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func ParsePagination(r *http.Request) Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	p := Pagination{Page: page, PerPage: perPage}
	p.Normalize()
	return p
}

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
}
