package main

import (
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	"github.com/xiangrui2019/go-kahla-bot-server/routers"
	"github.com/xiangrui2019/go-kahla-bot-server/server"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"os"
)

func main() {
	app := macaron.New()
	router := routers.NewRouter(app)
	injector := injects.NewInjector(app)
	config, err := conf.LoadConfigFromFile("./config.toml")
	pusherserver := server.NewPusherServer()
	interrupt := make(chan os.Signal, 1)
	interrupt2 := make(chan struct{})

	if err != nil {
		log.Fatal(err)
	}

	err = config.ConfigEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	err = router.ConfigureRouting()

	if err != nil {
		log.Fatal(err)
	}

	err = injector.Inject()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-interrupt
		close(interrupt2)
	}()

	go func() {
		err := pusherserver.Run(interrupt2)

		if err != nil {
			log.Fatal(err)
		}
	}()

	httpserver := http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: app,
	}

	log.Printf("listening on %s (%s)\n", fmt.Sprintf("%s:%d", config.Host, config.Port), macaron.Env)

	err = httpserver.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
