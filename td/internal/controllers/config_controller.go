package controllers

import (
	"fmt"
	"os"

	"github.com/tmobaird/dv/core"
)

type ConfigController struct {
	Base Controller
}

func (controller ConfigController) Run() (string, error) {
	config, err := core.ReadConfig(os.DirFS(core.BasePath()))
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("Current Config:\n  Context: %s\n  Hide Completed: %t", config.Context, config.HideCompleted), nil
	}
}
