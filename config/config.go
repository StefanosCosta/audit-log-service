package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var AuthConfig authConfig

type authConfig struct {
	PrivateKey []byte
	PublicKey  []byte
}

type KeyConfig struct {
	PrivateKey    string `yaml:"private_key"`
	PublicKey string `yaml:"public_key"`
}

func (config *authConfig) LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	keyConfig := &KeyConfig{}
	err = yaml.Unmarshal(data, keyConfig)
	if err != nil {
		return err
	}
	prvKey, err := os.ReadFile(keyConfig.PrivateKey)
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile(keyConfig.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	AuthConfig.PrivateKey = prvKey
	AuthConfig.PublicKey = pubKey

	return nil
}