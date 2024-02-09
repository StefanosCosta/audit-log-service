package main

import (
	"audit-log-service/config"
	"audit-log-service/db"
	eventsRepository "audit-log-service/eventsRepository"
	"audit-log-service/helpers"
	"audit-log-service/models"
	"fmt"
	"log"
	"net/http"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// LogInfo writes our payload to mongo
func (rpcServer *RPCServer) LogInfo(r *http.Request, event *models.EventPayload, resp *string) error {

	fmt.Println("Processing rpc call")
    reqToken := r.Header.Get("Authorization")
    fmt.Println(reqToken)

    if err := config.AuthConfig.Validate(reqToken); err != nil {
        *resp = "invalid access token"
        return err
    }
    // Save event to the database
    dbEvent := helpers.MapEventPayloadToDb(*event)
    eventRepo := eventsRepository.NewEventRepo(db.DBConn.DB,log.Default())
    err :=  eventRepo.Create(&dbEvent)
    if err != nil {
        *resp = "Could not create Event. Please try again later"
        return err
    } else{
        *resp = "Event Logged Successfully"
    }

	return nil
}
