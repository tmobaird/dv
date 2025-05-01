package main

import (
	"strconv"
)

type Config struct {
	HideCompleted bool `json:"hideCompleted"`
}

func (config *Config) set(key string, value string) {
	switch key {
	case "HideCompleted":
		b, err := strconv.ParseBool(value)
		if err == nil {
			config.HideCompleted = b
		}
	}
}
