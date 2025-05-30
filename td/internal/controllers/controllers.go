package controllers

import "github.com/tmobaird/dv/core"

const METADATA_STRING = "(duration=1b,type=blank)"

type Controller struct {
	Args   []string
	Config core.Config
}
