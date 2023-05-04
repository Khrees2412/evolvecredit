package repositories

import "gorm.io/gorm"

func paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
