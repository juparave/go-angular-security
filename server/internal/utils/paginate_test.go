package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestItem struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100"`
	Category string `gorm:"size:50"`
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&TestItem{})
	assert.NoError(t, err)

	// Seed 15 items
	for i := 1; i <= 15; i++ {
		db.Create(&TestItem{
			Name:     string(rune('A' + i - 1)),
			Category: "category_1",
		})
	}

	return db
}

func TestPaginate(t *testing.T) {
	db := setupTestDB(t)

	t.Run("first page with limit 5", func(t *testing.T) {
		params := PaginateParam{
			PageNum:  0,
			PageSize: 5,
		}

		res, err := Paginate[TestItem](db.Model(&TestItem{}), params)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), res.Total)
		assert.Equal(t, 5, len(res.Data))
		assert.Equal(t, 0, res.Page)
		assert.Equal(t, 5, res.PageSize)
		assert.Equal(t, 3, res.LastPage)
	})

	t.Run("second page with limit 5", func(t *testing.T) {
		params := PaginateParam{
			PageNum:  1,
			PageSize: 5,
		}

		res, err := Paginate[TestItem](db.Model(&TestItem{}), params)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), res.Total)
		assert.Equal(t, 5, len(res.Data))
		assert.Equal(t, 1, res.Page)
	})

	t.Run("with sorting and filtering", func(t *testing.T) {
		params := PaginateParam{
			PageNum:  0,
			PageSize: 5,
			SortCol:  "name",
			SortOrd:  "DESC",
		}

		res, err := Paginate[TestItem](db.Model(&TestItem{}), params)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(res.Data))
		assert.Equal(t, "O", res.Data[0].Name)
	})
}

func TestParsePaginateParams(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		params := ParsePaginateParams(c, 20)
		return c.JSON(params)
	})

	req := httptest.NewRequest("GET", "/test?page_index=2&page_size=10&sort_col=name&sort_ord=desc&filter=cenote", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
