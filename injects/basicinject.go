package injects

import (
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"gopkg.in/macaron.v1"
	"log"
	"os"
)

type BasicInject struct {
	context *macaron.Macaron
	logger *log.Logger
	client *kahla.Client
}

func NewInjector(ctx *macaron.Macaron) *BasicInject {
	c, _ := conf.LoadConfigFromFile("./config.toml")


	return &BasicInject{
		context: ctx,
		logger: log.New(os.Stdout, "[kahla-bot] ", 0),
		client: kahla.NewClient(c.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com"),
	}
}

func (inject *BasicInject) Inject() error {
	inject.context.Map(inject.logger)
	inject.context.Map(inject.client)
	return nil
}