package controllers

import (
	"fmt"
	"os"

	"github.com/tmobaird/dv/core"
	"github.com/tmobaird/dv/td/internal/models"
)

type AddController struct {
	Base           Controller
	MetadataString string
}

func (controller AddController) Run() (string, error) {
	file, err := os.OpenFile(core.TodoFilePath(controller.Base.Config.Context), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	todo := models.Todo{Name: controller.Base.Args[0], Complete: false, Metadata: models.ParseTodoMetadata(controller.MetadataString)}
	toWrite := todo.ToMd()
	_, err = file.Write([]byte(toWrite))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" added to list.", todo.Name), nil
}
