package injects

import (
	"gopkg.in/macaron.v1"
	"log"
	"os"
)

type BasicInject struct {
	context *macaron.Macaron
	logger *log.Logger
}

func NewInjector(ctx *macaron.Macaron) *BasicInject {
	return &BasicInject{
		context: ctx,
		logger: log.New(os.Stdout, "[kahla-bot] ", 0),
	}
}

func (inject *BasicInject) Inject() error {
	inject.context.Map(inject.logger)
	return nil
}