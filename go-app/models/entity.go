package models

import "gorm.io/gorm"

// Use Entity for paginate multiple models
type Entity interface {
	Count(db *gorm.DB) int64
	Take(db *gorm.DB, limit int, offset int) interface{}
}
