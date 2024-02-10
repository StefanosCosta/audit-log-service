package handlers

import (
	authenticationservice "audit-log-service/api/authservice"
	"audit-log-service/api/eventservice"
	"audit-log-service/config"
	"audit-log-service/db"
	eventsRepository "audit-log-service/db/eventsRepository"
	usersRepository "audit-log-service/db/usersRepository"
	"audit-log-service/helpers"
	"audit-log-service/models"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
    var userPayload models.UserPayload

    if err := helpers.DecodeUserJSON(&userPayload,r); err != nil {
        fmt.Println(err)
        resp := helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }
	
    userRepo := usersRepository.NewUserRepo(db.DBConn.DB,log.Default())
    authService := authenticationservice.NewAuthenticationService(db.DBConn,log.Default(),&userRepo,&config.AuthConfiguration)
    resp, status := authService.RegisterUser(userPayload.Email,userPayload.Password)
    helpers.WriteJSON(w,status,resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
    var userPayload models.UserPayload

    if err := helpers.DecodeUserJSON(&userPayload,r); err != nil {
        fmt.Println(err)
        resp := helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

	userRepo := usersRepository.NewUserRepo(db.DBConn.DB,log.Default())
    authService := authenticationservice.NewAuthenticationService(db.DBConn,log.Default(),&userRepo,&config.AuthConfiguration)
    
    resp, status, token := authService.LoginUser(userPayload.Email,userPayload.Password)
    if !resp.Error {
        cookie := &http.Cookie{
            Name:     "jwt-token",
            Value:    token,
            Path:     "/",
            Secure:   true,
            HttpOnly: true,
            Expires: time.Now().Add(time.Hour * 24),
        }
        http.SetCookie(w, cookie)
    }
    
	helpers.WriteJSON(w,status,resp)
}

func QueryEvents(w http.ResponseWriter, r *http.Request) {

    var (
        resp helpers.JsonResponse
        dbEvents []eventsRepository.Event
        eventsResponse []models.EventPayload
        status int
    )
    // Validate and authenticate the request
    reqToken := r.Header.Get("Authorization")

    if err := config.AuthConfiguration.Validate(reqToken); err != nil {
        resp = helpers.GetInvalidPayloadResponseWithMessage("Invalid access token")
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

    queryParams := r.URL.Query()
    eventRepo := eventsRepository.NewEventRepo(db.DBConn.DB,log.Default())
    eventService := eventservice.NewEventService(db.DBConn,log.Default(),&eventRepo)
    resp, status, dbEvents = eventService.QueryEvents(queryParams)

    if resp.Error {
        helpers.WriteJSON(w,status,resp)
        return
    }
    
    // Serialize results to JSON and respond
    eventsResponse = helpers.MapDbPayloadsToEvents(dbEvents)

    helpers.WriteJSON(w,status,eventsResponse)
    
}