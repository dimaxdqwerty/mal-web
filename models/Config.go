package models

import (
	"github.com/BurntSushi/toml"
	"log"
)

type config struct {
	MalClientID string
}

var conf config

func GetMalClientID() string {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	return conf.MalClientID
}
