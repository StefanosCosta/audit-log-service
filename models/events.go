package models

import (
	"time"

	"gorm.io/datatypes"
)

 type Event struct {
    ID        string    `json:"id,omitempty"`
    Timestamp time.Time `json:"timestamp,omitempty"`
    Type string    `json:"eventType,omitempty"`
    // Common fields
    ActorID string `json:"actorId,omitempty"`
    // Variant fields
    Details datatypes.JSON `json:"details,omitempty"`
    // Details interface{} `json:"details"`
}
