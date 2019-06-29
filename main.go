package main

import (
	"errors"
	"fmt"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/injects"
	_ "github.com/xiangrui2019/go-kahla-bot-server/orm"
	"github.com/xiangrui2019/go-kahla-bot-server/routers"
	"github.com/xiangrui2019/go-kahla-bot-server/server"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"os"
	"reflect"
)

func main() {
	app := macaron.New()
	injector := injects.NewInjector(app)

	err := injector.Inject()

	if err != nil {
		log.Fatal(err)
	}

	config := app.GetVal(reflect.TypeOf(injector.Config)).Interface().(*conf.Config)
	pusherserver := server.NewPusherServer(app, injector)
	router := routers.NewRouter(app, injector)
	interrupt := make(chan os.Signal, 1)
	interrupt2 := make(chan struct{})

	err = config.ConfigEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	err = router.ConfigureRouting()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-interrupt
		close(interrupt2)
	}()

	go func() {
		err := errors.New("inital error to run")

		for err != nil  {
			err = pusherserver.Run(interrupt2)

			if err != nil {
				log.Println(err)
			}
		}
	}()

	server := http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: app,
	}

	log.Printf("listening on %s (%s)\n", fmt.Sprintf("%s:%d", config.Host, config.Port), macaron.Env)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
