package authenticationservice_test

import (
	authenticationservice "audit-log-service/api/authservice"
	"audit-log-service/api/testutils"
	"audit-log-service/db"
	users "audit-log-service/db/usersRepository"
	mock_config "audit-log-service/mocks"
	"log"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser(t *testing.T) {
	assert.Nil(t, testutils.InitTestDb())
	defer testutils.RemoveTestDb()
	
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockAuthProvider := mock_config.NewMockAuthConfiger(mockCtrl)
	
	userRepo := users.NewUserRepo(db.DBConn.DB,log.Default())
	password := "12346789SW!"
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	assert.Nil(t,err)
	email := "mrpotatohead@gmail.com"

	newUser := users.User{Email: email, Password: pass,}
	userRepo.Create(&newUser)

	tokenString := "Token string"

	mockAuthProvider.EXPECT().GenerateToken(newUser.ID).Return(tokenString,nil)

	authSvc := authenticationservice.NewAuthenticationService(db.DBConn, log.Default(), &userRepo,mockAuthProvider)

	resp, status, token := authSvc.LoginUser(email, password)
	assert.Equal(t, http.StatusAccepted,status)
	assert.Equal(t, tokenString, token)
	assert.Equal(t,resp.Error, false)
	assert.Equal(t, resp.Message, "Login Successful")

}