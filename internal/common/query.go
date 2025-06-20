package common

import (
    "github.com/gin-gonic/gin"
    "strconv"
)

// QueryParams chứa các tham số query chung
type QueryParams struct {
    Page  int    `json:"page"`
    Limit int    `json:"limit"`
    Sort  string `json:"sort"`
}

// ParseQueryParams phân tích query parameters từ context
func ParseQueryParams(c *gin.Context) QueryParams {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    sort := c.DefaultQuery("sort", "created_at DESC")

    // Validate
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }
    // Optional: Validate sort để tránh SQL injection
    // Ví dụ: chỉ cho phép sort theo các cột cụ thể
    allowedSorts := map[string]bool{
        "created_at ASC":  true,
        "created_at DESC": true,
        "name ASC":        true,
        "name DESC":       true,
        // Thêm các cột khác nếu cần
    }
    if sort != "" && !allowedSorts[sort] {
        sort = "created_at DESC" // Giá trị mặc định nếu sort không hợp lệ
    }

    return QueryParams{
        Page:  page,
        Limit: limit,
        Sort:  sort,
    }
}