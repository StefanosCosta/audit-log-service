package handlers

import (
	"audit-log-service/db"
	eventsRepository "audit-log-service/eventsRepository"
	"audit-log-service/helpers"
	"audit-log-service/models"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
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

    var (
        resp helpers.JsonResponse
        scopes []func(db *gorm.DB) *gorm.DB
        err error
        events []models.EventPayload
    )

    queryParams := r.URL.Query()
    scopes, err = helpers.MapQueryParamsToScopes(queryParams)
    if err != nil {
        resp = helpers.GetInvalidPayloadResponseWithMessage(err.Error())
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

    // Parse query parameters for field=value
    eventRepo := eventsRepository.NewEventRepo(db.DBConn.DB,log.Default())
    eventResponse :=  eventRepo.Find(scopes...)
    if len(eventResponse) == 0 {
        helpers.WriteJSON(w,http.StatusBadRequest,helpers.GetInvalidPayloadResponseWithMessage("No matches found"))
        return
    }
    // Validate and authenticate the request
    // Query database for events matching the criteria
    // Serialize results to JSON and respond
    for _,event := range (eventResponse) {
        events = append(events, helpers.MapDbPayloadToEvent(event))
    }

    helpers.WriteJSON(w,http.StatusBadRequest,events)
}