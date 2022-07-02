package models

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	MalClientID string
}

var conf Config

func GetMalClientID() string {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	return conf.MalClientID
}
