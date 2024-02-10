package users

import "gorm.io/gorm"


func ByEmailEquals(email string) func (db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Where("email = (?)", email)
	}
}

