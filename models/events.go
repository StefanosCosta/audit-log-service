package models

import (
	"time"
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




