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
	// 创建新的馬卡龙实例
	app := macaron.New()
	// 创建基本服务注入
	injector := injects.NewInjector(app)

	// 注入服务
	err := injector.Inject()
	// 检查错误,如果出错,直接报错退出
	if err != nil {
		// 报错退出
		log.Fatal(err)
	}

	// 通过反射从注入中获取对应的服务
	config := app.GetVal(reflect.TypeOf(injector.Config)).Interface().(*conf.Config)
	// 创建推送监听器
	pusherserver := server.NewPusherServer(app, injector)
	// 创建API路由器
	router := routers.NewRouter(app, injector)
	// 创建两个Channel
	interrupt := make(chan os.Signal, 1)
	interrupt2 := make(chan struct{})

	// 配置馬卡龙运行环境 dev test prod
	err = config.ConfigEnvironment()

	// 检查错误,如果出错,直接报错退出
	if err != nil {
		// 报错退出
		log.Fatal(err)
	}

	// 配置路由器
	err = router.ConfigureRouting()

	// 检查错误,如果出错,直接报错退出
	if err != nil {
		// 报错退出
		log.Fatal(err)
	}

	go func() {
		<-interrupt
		close(interrupt2)
	}()

	// 运行监听器
	go func() {
		err := errors.New("inital error to run")

		for err != nil  {
			err = pusherserver.Run(interrupt2)

			if err != nil {
				log.Println(err)
			}
		}
	}()

	// 给馬卡龙实例创建Go HTTP服务器
	httpServer := http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: app,
	}

	// 打印执行日志
	log.Printf("listening on %s (%s)\n", fmt.Sprintf("%s:%d", config.Host, config.Port), macaron.Env)

	// 运行馬卡龙实例对应的HTTP服务器
	err = httpServer.ListenAndServe()

	// 检查错误,如果出错,直接报错退出
	if err != nil {
		// 报错退出
		log.Fatal(err)
	}
}