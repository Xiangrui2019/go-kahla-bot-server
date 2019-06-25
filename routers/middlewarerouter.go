package routers

import (
	"gopkg.in/macaron.v1"
)

func ConfigureMiddlewareRouting(context *macaron.Macaron) error {
	context.Use(macaron.Logger())
	context.Use(macaron.Recovery())

	return nil
}