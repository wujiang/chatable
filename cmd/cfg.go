package main

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Postgres        string `json:"postgres"`
	Host            string `json:"host"`
	RedisHost       string `json:"redis_host"`
	SharedQueueKey  string `json:"shared_queue_key"`
	QueueManagerKey string `json:"queue_manager_key"`
	MessageQueueKey string `json:"message_queue_key"`
}

func (cfg *configuration) parse(fn string) error {
	bt, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bt, cfg)
	return err
}
