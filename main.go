package main

import (
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/routers"
	"gopkg.in/macaron.v1"
)

func main() {
	app := macaron.New()
	router := routers.NewRouter(app)
	injector := injects.NewInjector(app)
	config, err := conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		panic(err)
	}

	err = config.ConfigEnvironment()

	if err != nil {
		panic(err)
	}

	err = router.ConfigureRouting()

	if err != nil {
		panic(err)
	}

	err = injector.Inject()

	if err != nil {
		panic(err)
	}

	app.Run(config.Host, config.Port)
}
