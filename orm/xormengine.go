package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/xiangrui2019/go-kahla-bot-server/conf"
	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"log"
)

var X *xorm.Engine

func init() {
	config, err := conf.LoadConfigFromFile("./config.toml")

	if err != nil {
		log.Fatal(err)
	}

	X, err = xorm.NewEngine("mysql", config.MySqlDSN)

	if err != nil {
		log.Fatal(err)
	}

	if err := X.Sync(new(models.BotUser)); err != nil {
		log.Fatal(err)
	}
}