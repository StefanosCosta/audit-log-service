package helpers

import (
	"audit-log-service/models"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// decodeJSON tries to read the body of a request and sets the decoded value to the event pointer passed to it
func DecodeJSON(event *models.Event, r *http.Request) error{
	err := json.NewDecoder(r.Body).Decode(event)
    if err != nil {
        
        return errors.Errorf("Could not decode request body of event: %s", err)
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