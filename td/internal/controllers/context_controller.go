package controllers

import (
	"fmt"

	"github.com/tmobaird/dv/core"
)


type ContextController struct {
	Base            Controller
	ChangeToDefault bool
}

func (controller ContextController) Run() (string, error) {
	if controller.ChangeToDefault {
		return controller.updateConfig("main")
	} else if len(controller.Base.Args) > 0 {
		return controller.updateConfig(controller.Base.Args[0])
	} else {
		return controller.Base.Config.Context, nil
	}
}

func (controller ContextController) updateConfig(context string) (string, error) {
	controller.Base.Config.Context = context
	err := core.PersistConfig(controller.Base.Config)
	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("Updated context to %s", context), nil
	}
}
