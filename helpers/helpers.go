package helpers

import (
	events "audit-log-service/eventsRepository"
	"audit-log-service/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JsonResponse struct {
    Error bool `json:"error"`
    Message string `json:"message"`
}

func GetInvalidPayloadResponse()(JsonResponse) {
    return JsonResponse{Error: true, Message: "Invalid payload",}
}

func GetInvalidPayloadResponseWithMessage(message string)(JsonResponse) {
    return JsonResponse{Error: true, Message:message,}
}

func GetSuccessfulEventSubmissionResponse()(JsonResponse) {
    return JsonResponse{Error: false, Message: "Event Logged Successfully",
}
}

func MapQueryParamsToScopes(queryParams map[string][]string) ([]func(db *gorm.DB) *gorm.DB, error) {
	var scopes []func(db *gorm.DB) *gorm.DB

	if len(queryParams["timestamp"]) >0 {
		timestamp, err := time.Parse("2020-06-30T18:00:00.000Z",queryParams["timestamp"][0])
		if err != nil {
			return scopes, errors.Errorf("Invalid timestamp query parameter %s", err)
		}
		scopes = append(scopes, events.ByTimestampGreaterThan(timestamp))
	}
	if len(queryParams["eventType"]) >0 {
		
		scopes = append(scopes, events.ByEventType(queryParams["eventType"][0]))
	}

	if len(queryParams["actorId"]) >0 {
		
		scopes = append(scopes, events.ByActorID(queryParams["actorId"][0]))
	}


	// for key, val := range(commonFields) {
	// 	if len(queryParams[key]) > 0 {
	// 		scopes = append(scopes, val)
	// 	}
	// }

	return scopes, nil
}


func MapEventPayloadToDb(event models.EventPayload) (events.Event){
	// var details string
	
	dbEvent := events.Event{
		Timestamp: event.Timestamp,
		Type: event.Type,
		ActorID: &event.ActorID,
	}

	if event.Details != nil {
		m, err := json.Marshal(event.Details)
		if err == nil {
			
			fmt.Println(m)
			dbEvent.Details = datatypes.JSON([]byte(m))
		}
	}
	
	return dbEvent
}

func MapDbPayloadToEvent(event events.Event) (models.EventPayload){
	// var details string
	eventPayload := models.EventPayload{
		Timestamp: event.Timestamp,
		Type: event.Type,
		ActorID: *event.ActorID,
		Details: event.Details,
	}
	
	return eventPayload
}

// decodeJSON tries to read the body of a request and sets the decoded value to the event pointer passed to it
func DecodeJSON(event *models.EventPayload, r *http.Request) error{
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&event)
    if err != nil {
        
        return errors.Errorf("Could not decode request body of event: %s", err)
    }
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

// writeJSON takes a response status code, other data and writes a json response to the client
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}