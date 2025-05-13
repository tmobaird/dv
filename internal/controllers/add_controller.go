package controllers

import (
	"fmt"
	"os"
	"td/internal"
	"td/internal/models"
)

type AddController struct {
	Base Controller
}

func (controller AddController) Run() (string, error) {
	file, err := os.OpenFile(internal.TodoFilePath(controller.Base.Config.Context), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	todo := models.Todo{Name: controller.Base.Args[0], Complete: false, Metadata: models.DefaultMetadata()}
	toWrite := todo.ToMd()
	_, err = file.Write([]byte(toWrite))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" added to list.", todo.Name), nil
}
