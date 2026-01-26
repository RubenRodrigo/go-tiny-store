package pagination

import (
	"math"
	"net/http"
	"strconv"
)

// Params
type Params struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p *Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Params) Limit() int {
	return p.PageSize
}

// Info
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type Result[T any] struct {
	Meta Meta `json:"meta"`
	Data []T  `json:"data"`
}

func ParseParams(r *http.Request) Params {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}

}

func BuildMeta(params Params, totalItems int64) Meta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PageSize)))

	return Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func BuildResult[T any](params Params, totalItems int64, data []T) Result[T] {
	return Result[T]{
		Meta: BuildMeta(params, totalItems),
		Data: data,
	}
}
