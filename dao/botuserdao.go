package dao

import (
	"errors"
	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"github.com/xiangrui2019/go-kahla-bot-server/orm"
)

func CreateBotUser(user *models.BotUser) error {
	_, err := orm.X.Insert(user)
	return err
}

func DeleteBotUser(id int64) error {
	_, err := orm.X.Delete(&models.BotUser{Id: id})
	return err
}

func DeleteBotUserByKahlaId(id string) error {
	_, err := orm.X.Delete(&models.BotUser{KahlaUserId: id})
	return err
}

func GetBotUserById(id int64) (*models.BotUser, error) {
	user := &models.BotUser{}

	has, err := orm.X.Id(id).Get(user)
	
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("User not found")
	}

	return user, nil
}

func GetBotUserByToken(token string) (*models.BotUser, error) {
	user := &models.BotUser{
		Token: token,
	}

	has, err := orm.X.Get(user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("User not found")
	}

	return user, nil
}
