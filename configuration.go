package main

import (
	"encoding/json"
	"github.com/Bluecodelf/rets"
	"io/ioutil"
)

type Configuration struct {
	SQL  rets.DatabaseConfiguration `json:"sql"`
	Port string                     `json:"port"`
}

func ReadConfiguration(path string) (configuration *Configuration, err error) {
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	configuration = new(Configuration)
	err = json.Unmarshal(data, configuration)
	return
}
