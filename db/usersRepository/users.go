package users

import (
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password []byte
	Role     string
}


type UserRepository struct {
	DbInstance *gorm.DB
	Logger *log.Logger
}


func NewUserRepo(db *gorm.DB, logger *log.Logger) (UserRepository){
	return UserRepository{DbInstance: db, Logger: logger}
}

func (userRepository *UserRepository) Find(scopes ...func(*gorm.DB) *gorm.DB) ([]User, error) {
	var users []User
	return users, userRepository.DbInstance.Scopes(scopes...).Find(&users).Error
}

func (userRepository *UserRepository) FindSingleUser(scopes ...func(*gorm.DB) *gorm.DB) (User, error) {
	var user User
	return user, userRepository.DbInstance.Scopes(scopes...).First(&user).Error
}

func (userRepository *UserRepository) Create(user *User) error {
	return userRepository.DbInstance.Create(user).Error
}