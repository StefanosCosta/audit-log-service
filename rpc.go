package main

import (
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

    // var event models.EventPayload
    // var resp helpers.JsonResponse
    // Parse request body into Event struct
    // if err := helpers.DecodeEventJSON(&event,r); err != nil {
    //     fmt.Println(err)
    //     // resp = helpers.GetInvalidPayloadResponse()
    //     // err = helpers.WriteJSON(w,http.StatusBadRequest,resp)
    //     return err
		// rpc.Dial()
    // }
	// var resp *string
	fmt.Println("Processing rpc call")

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
    // Respond with success or error
    
    // err =helpers.WriteJSON(w,http.StatusAccepted,resp)
	// if err != nil {
	// 	return err
	// }
	// resp is the message sent back to the RPC caller
	// *resp = "Processed payload via RPC:" + event.Type
	return nil
}
