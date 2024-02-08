package handlers

import (
	"audit-log-service/db"
	events "audit-log-service/eventsRepository"
	"audit-log-service/helpers"
	"audit-log-service/models"
	"fmt"
	"log"
	"net/http"
)



func HandleEvent(w http.ResponseWriter, r *http.Request) {
    var event models.EventPayload
    var resp helpers.JsonResponse
    // Parse request body into Event struct
    if err := helpers.DecodeJSON(&event,r); err != nil {
        fmt.Println(err)
        resp = helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }
    // Validate and authenticate the request

    // Save event to the database
    dbEvent := helpers.MapEventPayloadToDb(event)
    db.DBConn.DB.Create(&dbEvent)
    // Respond with success or error
    resp = helpers.GetSuccessfulEventSubmissionResponse()
    helpers.WriteJSON(w,http.StatusAccepted,resp)                    
}



func QueryEvents(w http.ResponseWriter, r *http.Request) {

    var event models.EventPayload
    var resp helpers.JsonResponse
    // Parse request body into Event struct
    if err := helpers.DecodeJSON(&event,r); err != nil {
        fmt.Println(err)
        resp = helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }
    // Parse query parameters for field=value
    eventRepo := events.NewEventRepo(db.DBConn.DB,log.Default())
    eventRepo.Find(events.ByEventType(event.Type),
    events.ByActorID(event.ActorID),
    events.ByTimestampGreaterThan(event.Timestamp))
    // Validate and authenticate the request
    // Query database for events matching the criteria
    // Serialize results to JSON and respond
}