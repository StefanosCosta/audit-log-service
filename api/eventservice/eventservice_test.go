package eventservice_test

import (
	"audit-log-service/api/eventservice"
	"audit-log-service/api/testutils"
	"audit-log-service/db"
	events "audit-log-service/db/eventsRepository"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)



type eventScenario struct {
	Event events.Event
	ErrorExpected bool
	Response string
}

type queryScenario struct {
	Description string
	QueryParams map[string][]string
	ExpectedStatus int
	ExpectedEvents []events.Event
	ErrorExpected bool
	Response string
}

func TestAddandQueryEvent(t *testing.T) {
	assert.Nil(t, testutils.InitTestDb())
	defer testutils.RemoveTestDb()
	var dbEventScenarios []eventScenario
	actorId := "1234"
	actorId2 := "12345"

	allEvents := []events.Event{
		{
			Timestamp: time.Now(),
			Type: "Customer Signup",
			ActorID: &actorId,
			Details: datatypes.JSON([]byte(`{"name": "mr", "surname" : "potato"}`)),
		},
		{
			Timestamp: time.Now(),
			Type: "Customer billing",
			ActorID: &actorId,
			Details: datatypes.JSON([]byte(`{"billingAmount": "12345"}`)),
		},
		{
			Timestamp: time.Now(),
			Type: "Pdf Downloaded",
			ActorID: &actorId,
			Details: datatypes.JSON([]byte(`{"pdf": "shopping_cart_receipt"}`)),
		},
		{
			Timestamp: time.Now(),
			Type: "Account Deactivation",
			ActorID: &actorId2,
			Details: datatypes.JSON([]byte(`{"reason": "unsatisfactory_service"}`)),
		},
	}

	dbEventScenarios = []eventScenario{
		{
			Event: allEvents[0],
			ErrorExpected: false,
			Response: "Event Logged Successfully",
		},
		{
			Event: allEvents[1],
			ErrorExpected: false,
			Response: "Event Logged Successfully",
		},
		{
			Event: allEvents[2],
			ErrorExpected: false,
			Response: "Event Logged Successfully",
		},
		{
			Event: allEvents[3],
			ErrorExpected: false,
			Response: "Event Logged Successfully",
		},
	}
		

	eventRepo := events.NewEventRepo(db.DBConn.DB,log.Default())
	eventSvc := eventservice.NewEventService(db.DBConn,log.Default(),&eventRepo)

	for _, eventScenario := range dbEventScenarios {
		resp, err :=eventSvc.CreateEvent(eventScenario.Event)
		if !eventScenario.ErrorExpected {
			assert.Nil(t,err)
		} else {
			assert.NotNil(t,err)
		}
		assert.Equal(t, eventScenario.Response, resp)
	}

	var queryScenarios []queryScenario
	
	
	queryScenarios = []queryScenario{
		{
			Description: "Looking for customer signup event with name json equal to mr - exists",
			QueryParams: map[string][]string{
				"timestamp" : {time.Now().Add(-time.Hour).Format(time.RFC3339)},
				"eventType" : {"Customer Signup"},
				"name" : {"mr"},
		}, ExpectedStatus: http.StatusAccepted,
		   ErrorExpected: false,
		   ExpectedEvents: []events.Event{dbEventScenarios[0].Event},
		},
		{
			Description: "Looking for event with billing amount in json - exists",
			QueryParams: map[string][]string{
				"billingAmount" : {"12345"},
		}, ExpectedStatus: http.StatusAccepted,
		   ErrorExpected: false,
		   ExpectedEvents: []events.Event{dbEventScenarios[1].Event},
		},
		{
			Description: "Looking for event with actorId - does not exist",
			QueryParams: map[string][]string{
				"actorId" : {"62"},
		}, ExpectedStatus: http.StatusBadRequest,
		   ErrorExpected: true,
		   ExpectedEvents: []events.Event{},
		},
		{
			Description: "Looking for events after 1 minute ago - multiple exist",
			QueryParams: map[string][]string{
				"timestamp" : {time.Now().Add(-time.Hour).Format(time.RFC3339)},
		}, ExpectedStatus: http.StatusAccepted,
		   ErrorExpected: false,
		   ExpectedEvents: allEvents,
		},
	}

	

	for _,query := range(queryScenarios) {
		resp, status, dbEvents := eventSvc.QueryEvents(query.QueryParams)
		assert.Equal(t,query.ExpectedStatus,status)
		assert.Equal(t,query.ErrorExpected, resp.Error)
		assert.Equal(t,len(query.ExpectedEvents),len(dbEvents))
		if len(query.ExpectedEvents) > 0 {
			for _, event := range(query.ExpectedEvents) {
				exists := false
				for _, returnedEvent := range(dbEvents) {
					if returnedEvent.Timestamp.Equal(event.Timestamp.UTC())  && *returnedEvent.ActorID == *event.ActorID && returnedEvent.Type == event.Type && reflect.DeepEqual(returnedEvent.Details,event.Details){
						exists = true
						break
					}
				}
				assert.True(t,exists)
			}
		}
	}
}