package routers

import "gopkg.in/macaron.v1"

type Router struct {
	context *macaron.Macaron
}

func NewRouter(ctx *macaron.Macaron) *Router {
	return &Router{
		context: ctx,
	}
}

func (r *Router) ConfigureRouting() error {
	err := ConfigureMiddlewareRouting(r.context)

	if err != nil {
		return err
	}

	err = ConfigureServiceRouting(r.context)

	if err != nil {
		return err
	}

	return nil
}
