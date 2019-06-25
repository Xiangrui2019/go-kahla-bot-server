package server

import (
	"errors"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/enums"
	"github.com/xiangrui2019/go-kahla-bot-server/kahla"
)

type PusherEventServer struct {
	client *kahla.Client
	config *conf.Config
}

func NewPusherServer() *PusherEventServer {
	var err error

	server := new(PusherEventServer)

	server.config, err = conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		return nil
	}

	server.client = kahla.NewClient(server.config.BotConfig.KahlaServer, "https://oss.cdn.aiursoft.com")

	return server
}

func (server *PusherEventServer) login() error {
	response, _, err := server.client.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{})

	if err != nil {
		return err
	}

	if response.Code == enums.ResponseCodeOK {
		return errors.New(response.Message)
	}

	return nil
}

func (server *PusherEventServer) initpusher() (*string, error) {
	response, _, err := server.client.Auth.InitPusher()

	if err != nil {
		return nil, err
	}

	if response.Code == enums.ResponseCodeOK {
		return nil, errors.New(response.Message)
	}

	return &response.ServerPath, nil
}

func (server *PusherEventServer) EventHandler(i interface{}) {

}

func (server *PusherEventServer) Run(interrupt chan struct{}) error {
	return nil
}