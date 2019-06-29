package routers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"gopkg.in/macaron.v1"
)

type Router struct {
	context *macaron.Macaron
	inject *injects.BasicInject
}

func NewRouter(ctx *macaron.Macaron, injector *injects.BasicInject) *Router {
	return &Router{
		context: ctx,
		inject: injector,
	}
}

// 配置全部路由
func (r *Router) ConfigureRouting() error {
	err := ConfigureMiddlewareRouting(r.context)

	if err != nil {
		return err
	}

	err = ConfigureServiceRouting(r.context, r.inject)

	if err != nil {
		return err
	}

	return nil
}
