package config

// go generate: mockgen -destination=mocks/mock_config.go -source=config/config.go

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var AuthConfiguration AuthConfig


type AuthConfiger interface {
	GenerateToken(id uint) (string, error)
	Validate(token string) ( error)
}

type AuthConfig struct {
	PrivateKey []byte
	PublicKey  []byte
}


type KeyConfig struct {
	PrivateKey    string `yaml:"private_key"`
	PublicKey string `yaml:"public_key"`
}

func LoadConfig(filename string) (AuthConfig, error) {
	var authConfig AuthConfig
	data, err := os.ReadFile(filename)
	if err != nil {
		return authConfig, err
	}
	keyConfig := &KeyConfig{}
	err = yaml.Unmarshal(data, keyConfig)
	if err != nil {
		return authConfig, err
	}
	prvKey, err := os.ReadFile(keyConfig.PrivateKey)
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile(keyConfig.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	authConfig.PrivateKey = prvKey
	authConfig.PublicKey = pubKey

	return authConfig, nil
}

func (config *AuthConfig) GenerateToken(id uint) (string, error) {
	var token string
	var err error
	
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.PrivateKey))
	if err != nil {
		return token, err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})
    
    
	token, err = claims.SignedString(key)

	
	return token, err
}


func (config *AuthConfig) Validate(token string) ( error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(config.PublicKey)
	if err != nil {
		return  errors.Errorf("validate: parse key: %s", err)
	}
 
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
 
		return key, nil
	})
	if err != nil {
		return errors.Errorf("validate: %s", err)
	}
 
	_, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return errors.Errorf("validate: invalid")
	}
 
	return nil
}