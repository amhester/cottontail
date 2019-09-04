package main

import (
	"encoding/json"
	"io/ioutil"
)

// QueueConfigs Provides a set of alert configurations for different sets of queues and their alert triggers
type QueueConfigs struct {
	Default *QueueConfig            `json:"default"`
	Queues  map[string]*QueueConfig `json:"queues"`
}

// QueueConfig provides a set of alert configurations for different types of triggers
type QueueConfig struct {
	Threshold *AlertConfig `json:"threshold"`
	Deviation *AlertConfig `json:"deviation"`
	Consumer  *AlertConfig `json:"consumer"`
}

// AlertConfig defines whether a certain type of alert trigger is enabled and what the triggering value is
type AlertConfig struct {
	Enabled bool `json:"enabled"`
	Value   int  `json:"value"`
}

// GetAlertConfig parses a json file specifed via the path argument and returns a set of configurations for queue based alert triggers
func GetAlertConfig(path string) (QueueConfigs, error) {
	var conf QueueConfigs
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
