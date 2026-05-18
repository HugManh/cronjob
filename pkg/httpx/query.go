package httpx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
	DefaultSort  = "created_at DESC"
)

var allowedSorts = map[string]struct{}{
	"created_at ASC":  {},
	"created_at DESC": {},
	"name ASC":        {},
	"name DESC":       {},
}

type QueryParams struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func ParseQueryParams(c *gin.Context) QueryParams {
	page := parseInt(c.DefaultQuery("page", strconv.Itoa(DefaultPage)), DefaultPage)
	limit := parseInt(c.DefaultQuery("limit", strconv.Itoa(DefaultLimit)), DefaultLimit)
	sort := c.DefaultQuery("sort", DefaultSort)

	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 || limit > MaxLimit {
		limit = DefaultLimit
	}
	if _, ok := allowedSorts[sort]; !ok {
		sort = DefaultSort
	}

	return QueryParams{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
