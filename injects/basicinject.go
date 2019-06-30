package injects

import (
	"log"
	"os"

	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
	"gopkg.in/macaron.v1"
)

type BasicInject struct {
	Context *macaron.Macaron
	Logger  *log.Logger
	Client  *kahla.Client
	Config  *conf.Config
}

func NewInjector(ctx *macaron.Macaron) *BasicInject {
	c, _ := conf.LoadConfigFromFile("./config.toml")

	return &BasicInject{
		Context: ctx,
		Logger:  log.New(os.Stdout, "[kahla-bot] ", 0),
		Client:  kahla.NewClient(c.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com"),
		Config:  c,
	}
}

// 注册所有需要注入的对象
func (inject *BasicInject) Inject() error {
	inject.Context.Map(inject.Logger)
	inject.Context.Map(inject.Client)
	inject.Context.Map(inject.Config)
	return nil
}
