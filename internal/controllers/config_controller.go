package controllers

import (
	"fmt"
	"os"
	"td/internal"
)

type ConfigController struct {
	Base Controller
}

func (controller ConfigController) Run() (string, error) {
	config, err := internal.Read(os.DirFS(internal.BasePath()))
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("Current Config:\n  Context: %s", config.Context), nil
	}
}
