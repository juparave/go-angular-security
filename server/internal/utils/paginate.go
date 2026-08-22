package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var safeSortColRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`)

// PaginateParam represents parameters for pagination, sorting, and filtering
type PaginateParam struct {
	PageNum   int    `json:"pageNum" query:"page"`
	PageSize  int    `json:"pageSize" query:"page_size"`
	SortCol   string `json:"sortCol" query:"sort_col"`
	SortOrd   string `json:"sortOrd" query:"sort_ord"`
	Filter    string `json:"filter" query:"filter"`
	DateRange string `json:"dateRange" query:"date_range"`
}

// PaginateResponse represents the generic pagination response payload
// 100% compatible with Angular tables and standard frontend paginators
type PaginateResponse[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	LastPage int   `json:"last_page"`
}

// ParsePaginateParams extracts and sanitizes pagination parameters from a Fiber context.
// Supports both Angular Material (page_index, page_size) and standard (page, limit) query conventions.
func ParsePaginateParams(c *fiber.Ctx, defaultPageSize ...int) PaginateParam {
	defSize := 25
	if len(defaultPageSize) > 0 && defaultPageSize[0] > 0 {
		defSize = defaultPageSize[0]
	}

	// Support page_index, page, pageNum
	pageStr := c.Query("page_index", c.Query("page", c.Query("pageNum", "0")))
	pageNum, err := strconv.Atoi(pageStr)
	if err != nil || pageNum < 0 {
		pageNum = 0
	}

	// Support page_size, pageSize, limit
	sizeStr := c.Query("page_size", c.Query("pageSize", c.Query("limit", strconv.Itoa(defSize))))
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = defSize
	}
	if pageSize > 100 {
		pageSize = 100 // Cap max page size
	}

	// Support sort_col, sortCol, sort
	sortCol := strings.TrimSpace(c.Query("sort_col", c.Query("sortCol", c.Query("sort", ""))))
	if !safeSortColRegex.MatchString(sortCol) {
		sortCol = "" // Reject invalid characters for SQL safety
	}

	// Support sort_ord, sortOrd, order
	sortOrd := strings.ToUpper(strings.TrimSpace(c.Query("sort_ord", c.Query("sortOrd", c.Query("order", "ASC")))))
	if sortOrd != "DESC" {
		sortOrd = "ASC"
	}

	return PaginateParam{
		PageNum:   pageNum,
		PageSize:  pageSize,
		SortCol:   sortCol,
		SortOrd:   sortOrd,
		Filter:    strings.TrimSpace(c.Query("filter", c.Query("q", ""))),
		DateRange: strings.TrimSpace(c.Query("date_range", "")),
	}
}

// Paginate executes a type-safe paginated query against GORM.
// It computes the total count safely using an isolated session without limit/offset/order,
// applies limits, offsets, and ordering, and returns a strongly-typed PaginateResponse[T].
func Paginate[T any](db *gorm.DB, params PaginateParam) (PaginateResponse[T], error) {
	var items []T
	var total int64

	// 1. Calculate total record count using an isolated session without limit/offset/order
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Model(new(T)).Count(&total).Error; err != nil {
		return PaginateResponse[T]{
			Data:     []T{},
			Total:    0,
			Page:     params.PageNum,
			PageSize: params.PageSize,
			LastPage: 0,
		}, fmt.Errorf("failed to count records: %w", err)
	}

	// Default & limit boundaries
	limit := params.PageSize
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	offset := params.PageNum * limit

	// 2. Apply sorting if a valid column was provided
	queryDB := db
	if params.SortCol != "" && safeSortColRegex.MatchString(params.SortCol) {
		order := "ASC"
		if strings.EqualFold(params.SortOrd, "desc") {
			order = "DESC"
		}
		queryDB = queryDB.Order(fmt.Sprintf("%s %s", params.SortCol, order))
	}

	// 3. Query items
	if err := queryDB.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return PaginateResponse[T]{
			Data:     []T{},
			Total:    total,
			Page:     params.PageNum,
			PageSize: limit,
			LastPage: calculateLastPage(total, limit),
		}, fmt.Errorf("failed to fetch paginated records: %w", err)
	}

	// Ensure non-nil slice in JSON output
	if items == nil {
		items = []T{}
	}

	return PaginateResponse[T]{
		Data:     items,
		Total:    total,
		Page:     params.PageNum,
		PageSize: limit,
		LastPage: calculateLastPage(total, limit),
	}, nil
}

func calculateLastPage(total int64, limit int) int {
	if total == 0 || limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}
