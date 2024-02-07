package handlers

import (
	"audit-log-service/db"
	"audit-log-service/helpers"
	"audit-log-service/models"
	"fmt"
	"net/http"
)

type jsonResponse struct {
    Error bool `json:"error"`
    Message string `json:"message"`

}

func HandleEvent(w http.ResponseWriter, r *http.Request) {
    var event models.Event
    var resp jsonResponse
    // Parse request body into Event struct
    if err := helpers.DecodeJSON(&event,r); err != nil {
        fmt.Println(err)
        resp = jsonResponse{Error: true,
            Message: "Invalid payload",}
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }
    // Validate and authenticate the request
    // Save event to the database
    db.DBConn.DB.Save(&event)
    // Respond with success or error
    resp = jsonResponse{Error: false,
                         Message: "Event Logged Successfully",
    }
    helpers.WriteJSON(w,http.StatusAccepted,resp)

                    
}

func QueryEvents(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for field=value
    // Validate and authenticate the request
    // Query database for events matching the criteria
    // Serialize results to JSON and respond
}