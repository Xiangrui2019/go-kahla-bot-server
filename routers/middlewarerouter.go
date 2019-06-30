package routers

import (
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/middlewares"
	"gopkg.in/macaron.v1"
)

var config *conf.Config

// 配置中间件路由
func ConfigureMiddlewareRouting(context *macaron.Macaron) error {
	var err error
	config, err = conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		return err
	}

	context.Use(macaron.Logger())
	context.Use(macaron.Recovery())
	context.Use(macaron.Renderer())
	context.Use(middlewares.CorsMiddleware(config.CorsOriginURL))

	return nil
}
