package handlers

import (
	"audit-log-service/config"
	"audit-log-service/db"
	eventsRepository "audit-log-service/eventsRepository"
	"audit-log-service/helpers"
	"audit-log-service/models"
	usersRepository "audit-log-service/usersRepository"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Register(w http.ResponseWriter, r *http.Request) {
    var err error
    var userPayload models.UserPayload

    if err := helpers.DecodeUserJSON(&userPayload,r); err != nil {
        fmt.Println(err)
        resp := helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }
	
    userRepo := usersRepository.NewUserRepo(db.DBConn.DB,log.Default())
    _, err = userRepo.FindSingleUser(usersRepository.ByEmailEquals(userPayload.Email))
    if err != nil {
        if errors.Is(errors.Cause(err), gorm.ErrRecordNotFound) {
            password, err  := bcrypt.GenerateFromPassword([]byte(userPayload.Password), 14)
            if err != nil {
                resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
                helpers.WriteJSON(w,http.StatusInternalServerError,resp)
                return
            }
            user := usersRepository.User{
		
                Email:    userPayload.Email,
                Password: password,
                // Role:     "User",
            }
			if err := userRepo.Create(&user); err != nil {
                resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
                helpers.WriteJSON(w,http.StatusInternalServerError,resp)
                return
            }
            resp := helpers.GetResponseWithMessage("Registration Successful",false)
            helpers.WriteJSON(w,http.StatusAccepted,resp)
		} else {
            resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
            helpers.WriteJSON(w,http.StatusInternalServerError,resp)
		}
    } else {
        resp := helpers.GetInvalidPayloadResponseWithMessage("Email or Password Already Exist")
		helpers.WriteJSON(w,http.StatusConflict,resp)
        return
	}
    
}

func Login(w http.ResponseWriter, r *http.Request) {
    var user []usersRepository.User
    var err error
    var userPayload models.UserPayload

    if err := helpers.DecodeUserJSON(&userPayload,r); err != nil {
        fmt.Println(err)
        resp := helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

	userRepo := usersRepository.NewUserRepo(db.DBConn.DB,log.Default())
    user, err = userRepo.Find(usersRepository.ByEmailEquals(userPayload.Email))
    if err != nil {
        return
    }

	if errors.Is(errors.Cause(err), gorm.ErrRecordNotFound) {
		resp := helpers.GetResponseWithMessage("Email or Password does not exist", false)
		helpers.WriteJSON(w,http.StatusNotFound,resp)
        return
	}

	if err := bcrypt.CompareHashAndPassword(user[0].Password, []byte(userPayload.Password)); err != nil {
		resp := helpers.GetResponseWithMessage("Email or Password incorrect", false)
		helpers.WriteJSON(w,http.StatusBadRequest,resp)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user[0].ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})
    
    key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.AuthConfig.PrivateKey))
	token, err := claims.SignedString(key)

	if err != nil {
		resp := helpers.GetResponseWithMessage("Login Failed", false)
		helpers.WriteJSON(w,http.StatusInternalServerError,resp)
		return 
	}

    cookie := &http.Cookie{
        Name:     "jwt-token",
        Value:    token,
        Path:     "/",
        Secure:   true,
        HttpOnly: true,
        Expires: time.Now().Add(time.Hour * 24),
    }
    http.SetCookie(w, cookie)


	helpers.WriteJSON(w,http.StatusAccepted,"Login Successful")
}

func HandleEvent(w http.ResponseWriter, r *http.Request) {
    // Validate and authenticate the request
    // reqToken := r.Header.Get("Authorization")
    // fmt.Println(reqToken)
    // splitToken := strings.Split(reqToken, "Bearer ")
    // reqToken = splitToken[1]
    
    var event models.EventPayload
    var resp helpers.JsonResponse
    // Parse request body into Event struct
    if err := helpers.DecodeEventJSON(&event,r); err != nil {
        fmt.Println(err)
        resp = helpers.GetInvalidPayloadResponse()
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

    // Save event to the database
    dbEvent := helpers.MapEventPayloadToDb(event)
    eventRepo := eventsRepository.NewEventRepo(db.DBConn.DB,log.Default())
    err :=  eventRepo.Create(&dbEvent)
    if err != nil {
        resp = helpers.GetInvalidPayloadResponseWithMessage("Could not create Event. Please try again later")
    } else{
        resp = helpers.GetSuccessfulEventSubmissionResponse()
    }
    // Respond with success or error
    
    helpers.WriteJSON(w,http.StatusAccepted,resp)                    
}



func QueryEvents(w http.ResponseWriter, r *http.Request) {

    var (
        resp helpers.JsonResponse
        scopes []func(db *gorm.DB) *gorm.DB
        jsonQueries []*datatypes.JSONQueryExpression

        err error
        events []models.EventPayload
    )

    reqToken := r.Header.Get("Authorization")
    fmt.Println(reqToken)

    queryParams := r.URL.Query()
    scopes,jsonQueries, err = helpers.MapQueryParamsToScopes(queryParams)
    if err != nil {
        resp = helpers.GetInvalidPayloadResponseWithMessage(err.Error())
        helpers.WriteJSON(w,http.StatusBadRequest,resp)
        return
    }

    // Parse query parameters for field=value
    eventRepo := eventsRepository.NewEventRepo(db.DBConn.DB,log.Default())
    eventResponse :=  eventRepo.Find(jsonQueries,scopes...)
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