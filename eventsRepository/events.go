package events

import (
	"log"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Scope func(*gorm.DB) *gorm.DB

type Event struct {
    // Common fields
    gorm.Model
    Timestamp time.Time
    Type string
    ActorID *string 
    // Variant fields
    Details datatypes.JSON
}

type EventRepository struct {
	DbInstance *gorm.DB
	Logger *log.Logger
}


func NewEventRepo(db *gorm.DB, logger *log.Logger) (EventRepository){
	return EventRepository{DbInstance: db, Logger: logger}
}

func (eventRepository *EventRepository) Find(jsonQueries []*datatypes.JSONQueryExpression, scopes ...func(*gorm.DB) *gorm.DB) ([]Event) {
	var events []Event
	eventRepository.DbInstance.Scopes(scopes...).Find(&events, jsonQueries)
	return events
}

func (eventRepository *EventRepository) Create(event *Event) error {
	return eventRepository.DbInstance.Create(event).Error
}
