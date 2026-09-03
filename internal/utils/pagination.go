package utils

import (
	"net/http"
	"strconv"
)

const (
	DefaultPaginationLimit = 10
	DefaultPaginationPage  = 1
	MaxPaginationLimit     = 100
)

type Pagination struct {
	Limit int
	Page  int
}

type PaginationMeta struct {
	Limit      int   `json:"limit"`
	Page       int   `json:"page"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

type QueryParamReader interface {
	QueryParam(name string) string
}

func ParsePagination(queryParams QueryParamReader) (Pagination, error) {
	pagination := Pagination{
		Limit: DefaultPaginationLimit,
		Page:  DefaultPaginationPage,
	}

	if value := queryParams.QueryParam("limit"); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil || limit < 1 || limit > MaxPaginationLimit {
			return Pagination{}, APIError(http.StatusBadRequest, "limit must be between 1 and 100")
		}
		pagination.Limit = limit
	}

	if value := queryParams.QueryParam("page"); value != "" {
		page, err := strconv.Atoi(value)
		if err != nil || page < 1 {
			return Pagination{}, APIError(http.StatusBadRequest, "page must be a positive integer")
		}
		pagination.Page = page
	}

	return pagination, nil
}

func (p Pagination) Offset() int {
	if p.Page <= 1 || p.Limit <= 0 {
		return 0
	}

	return (p.Page - 1) * p.Limit
}

func (p Pagination) Metadata(total int64) PaginationMeta {
	totalPages := int64(0)
	if total > 0 && p.Limit > 0 {
		totalPages = (total + int64(p.Limit) - 1) / int64(p.Limit)
	}

	return PaginationMeta{
		Limit:      p.Limit,
		Page:       p.Page,
		Total:      total,
		TotalPages: totalPages,
	}
}
