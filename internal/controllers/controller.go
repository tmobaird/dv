package controllers

import "td/internal"

const METADATA_STRING = "(duration=1b,type=blank)"

type Controller struct {
	Args   []string
	Config internal.Config
}
