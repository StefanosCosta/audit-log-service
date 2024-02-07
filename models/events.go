package models

import (
	"time"

	"gorm.io/datatypes"
)

// type Customer struct {
// 	CustomerName string
// 	CustomerType string
// 	CustomerInfo datatypes.JSON
//  }

 type Event struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Type string    `json:"eventType"`
    // Common fields
    ActorID string `json:"actorId"`
    // Variant fields
    Details datatypes.JSON `json:"details"`
    // Details interface{} `json:"details"`
}
