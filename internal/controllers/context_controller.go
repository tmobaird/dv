package controllers

import (
	"fmt"
	"td/internal"
)

type ContextController struct {
	Base Controller
}

func (controller ContextController) Run() (string, error) {
	if len(controller.Base.Args) > 0 {
		controller.Base.Config.Context = controller.Base.Args[0]
		err := internal.PersistConfig(controller.Base.Config)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Updated context to %s", controller.Base.Config.Context), nil
		}
	} else {
		return controller.Base.Config.Context, nil
	}
}
