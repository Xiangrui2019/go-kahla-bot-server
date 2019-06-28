package injects

import (
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"gopkg.in/macaron.v1"
	"log"
	"os"
)

type BasicInject struct {
	Context *macaron.Macaron
	Logger *log.Logger
	Client *kahla.Client
}

func NewInjector(ctx *macaron.Macaron) *BasicInject {
	c, _ := conf.LoadConfigFromFile("./config.toml")


	return &BasicInject{
		Context: ctx,
		Logger: log.New(os.Stdout, "[kahla-bot] ", 0),
		Client: kahla.NewClient(c.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com"),
	}
}

func (inject *BasicInject) Inject() error {
	inject.Context.Map(inject.Logger)
	inject.Context.Map(inject.Client)
	return nil
}