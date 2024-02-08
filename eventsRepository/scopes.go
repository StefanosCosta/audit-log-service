package events

import (
	"time"

	"gorm.io/gorm"
)

func ByEventType( eventType string) func (db *gorm.DB) *gorm.DB{
	return func (db *gorm.DB) *gorm.DB {
		return db.Where("type like (?)", eventType)
	}
}

func ByActorID(actorId string) func (db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Where("actor_id = (?)", actorId)
	}
}

func ByTimestampGreaterThan(timestamp time.Time) func (db *gorm.DB) *gorm.DB  {
	return func (db *gorm.DB) *gorm.DB {
		return db.Scopes().Where("timestamp > (?)", timestamp) 
	}
}

// func ByDetails(db *gorm.DB, ) *gorm.DB {
// 	datatypes.JSONQuery("attributes").Equals("jinzhu", "name")
// }