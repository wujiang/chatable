package main

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Postgres string `json:"postgres"`
	Host     string `json:"host"`
}

func (cfg *configuration) parse(fn string) error {
	bt, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bt, cfg)
	return err
}
