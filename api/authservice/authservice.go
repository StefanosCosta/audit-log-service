package authenticationservice

import (
	"audit-log-service/config"
	"audit-log-service/db"
	users "audit-log-service/db/usersRepository"
	usersRepository "audit-log-service/db/usersRepository"
	"audit-log-service/helpers"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type authenticationService struct {
	DB *db.DBConnection
	Logger *log.Logger
	UserRepo *users.UserRepository
	AuthProvider config.AuthConfiger
}

func NewAuthenticationService( db *db.DBConnection,
							   logger *log.Logger,
							   userRepo *users.UserRepository,
							   authProvider config.AuthConfiger) (authenticationService){
	return authenticationService{DB: db,Logger: logger,UserRepo: userRepo,AuthProvider: authProvider}
}

func (authService *authenticationService) RegisterUser(email string, password string) (helpers.JsonResponse,int){
	var err error
	_, err = authService.UserRepo.FindSingleUser(usersRepository.ByEmailEquals(email))
    if err != nil {
        if errors.Is(errors.Cause(err), gorm.ErrRecordNotFound) {
            password, err  := bcrypt.GenerateFromPassword([]byte(password), 14)
            if err != nil {
                resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
                return resp, http.StatusInternalServerError
            }
            user := usersRepository.User{
                Email:    email,
                Password: password,
            }
			if err := authService.UserRepo.Create(&user); err != nil {
                resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
                
                return resp, http.StatusInternalServerError
            }
            resp := helpers.GetResponseWithMessage("Registration Successful",false)
            return resp, http.StatusAccepted
		} else {
            resp := helpers.GetInvalidPayloadResponseWithMessage("Failed to register")
            return resp, http.StatusInternalServerError
		}
    } else {
        resp := helpers.GetInvalidPayloadResponseWithMessage("Email or Password Already Exist")
        return resp, http.StatusConflict
	}
}


func (authService *authenticationService) LoginUser(email string, password string) (helpers.JsonResponse,int, string){
	var resp helpers.JsonResponse
	var token string
	user, err := authService.UserRepo.Find(usersRepository.ByEmailEquals(email))
    

	if errors.Is(errors.Cause(err), gorm.ErrRecordNotFound) {
        return helpers.GetResponseWithMessage("Email or Password does not exist", false), http.StatusNotFound, token
	} else if err != nil{
		return resp, http.StatusInternalServerError, token
	}

	if err := bcrypt.CompareHashAndPassword(user[0].Password, []byte(password)); err != nil {
		return helpers.GetResponseWithMessage("Email or Password incorrect", false), http.StatusBadRequest, token
	}

    token, err = authService.AuthProvider.GenerateToken(user[0].ID)

    if err != nil {
		return helpers.GetResponseWithMessage("Login Failed", true), http.StatusInternalServerError, token
	}
	return helpers.GetResponseWithMessage("Login Successful", false), http.StatusAccepted,token
}



