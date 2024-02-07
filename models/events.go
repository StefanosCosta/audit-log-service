package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventPayload struct {
    Timestamp time.Time `json:"timestamp,omitempty"`
    Type string    `json:"eventType,omitempty"`
    // Common fields
    ActorID string `json:"actorId,omitempty"`
    // Variant fields
    // Details datatypes.JSON `json:"details,omitempty"`
    Details interface{} `json:"details,omitempty"`
}

type Event struct {
    // Common fields
    gorm.Model
    Timestamp time.Time
    Type string
    ActorID *string 
    // Variant fields
    Details datatypes.JSON
}



