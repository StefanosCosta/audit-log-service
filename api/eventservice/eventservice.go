package eventservice

import (
	"audit-log-service/db"
	events "audit-log-service/db/eventsRepository"
	"audit-log-service/helpers"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)


type eventService struct {
	DB *db.DBConnection
	Logger *log.Logger
	EventRepo *events.EventRepository
}

func NewEventService( db *db.DBConnection,
							   logger *log.Logger,
							   userRepo *events.EventRepository,
							   ) (eventService){
	return eventService{DB: db,Logger: logger,EventRepo: userRepo}
}

func (eventSvc *eventService) QueryEvents(queryParams map[string][]string) (helpers.JsonResponse, int, []events.Event){
	var (
		scopes []func(db *gorm.DB) *gorm.DB
		jsonQueries []*datatypes.JSONQueryExpression
		err error
		eventResponse []events.Event
	)

	// Parse query parameters for field=value
	scopes,jsonQueries, err = eventSvc.mapQueryParamsToScopes(queryParams)
    if err != nil {
        return helpers.GetInvalidPayloadResponseWithMessage(err.Error()),http.StatusBadRequest, eventResponse
    }

    // Query database for events matching the criteria
    eventResponse =  eventSvc.EventRepo.Find(jsonQueries,scopes...)
    if len(eventResponse) == 0 {
        return helpers.GetInvalidPayloadResponseWithMessage("No matches found"), http.StatusBadRequest, eventResponse
    }

	return helpers.GetSuccessfulEventSubmissionResponse(),http.StatusAccepted, eventResponse
}

// TODO Add more logic for more complex and flexible querying
// maps query parameters to appropriate scope for querying the db
func (eventSvc *eventService) mapQueryParamsToScopes(queryParams map[string][]string) ([]func(db *gorm.DB) *gorm.DB,[]*datatypes.JSONQueryExpression, error) {
	var (
		scopes []func(db *gorm.DB) *gorm.DB
	    jsonQueries []*datatypes.JSONQueryExpression
	)

	if len(queryParams["timestamp"]) > 0 {
		timestamp, err := time.Parse(time.RFC3339,queryParams["timestamp"][0])
		if err != nil {
			return scopes,jsonQueries, errors.Errorf("Invalid timestamp query parameter %s", err)
		}
		scopes = append(scopes, events.ByTimestampGreaterThan(timestamp))
	}
	delete(queryParams,"timestamp")
	if len(queryParams["eventType"]) > 0 {
		
		scopes = append(scopes, events.ByEventType(queryParams["eventType"][0]))
	}
	delete(queryParams,"eventType")
	if len(queryParams["actorId"]) > 0 {
		
		scopes = append(scopes, events.ByActorID(queryParams["actorId"][0]))
	}
	delete(queryParams,"actorId")
	
	if len(queryParams) > 0 {
		for key, val := range(queryParams) {
			jsonQueries = append(jsonQueries, datatypes.JSONQuery("details").Equals(val[0], key) )
		}
	}

	return scopes,jsonQueries, nil
}

func (eventSvc *eventService) CreateEvent(dbEvent events.Event) (resp string, err error){
	err = eventSvc.EventRepo.Create(&dbEvent)
    if err != nil {
        resp = "Could not create Event. Please try again later"
    } else{
        resp = "Event Logged Successfully"
    }

	return resp, err
}


