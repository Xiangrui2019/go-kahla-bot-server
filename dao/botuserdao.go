package dao

import (
	"errors"

	"github.com/xiangrui2019/go-kahla-bot-server/models"
	"github.com/xiangrui2019/go-kahla-bot-server/orm"
)

// 创建Bot用户
func CreateBotUser(user *models.BotUser) error {
	_, err := orm.X.Insert(user)
	return err
}

// 删除Bot用户
func DeleteBotUser(id int64) error {
	_, err := orm.X.Delete(&models.BotUser{Id: id})
	return err
}

// 通过卡拉UserID删除Bot用户
func DeleteBotUserByKahlaId(id string) error {
	_, err := orm.X.Delete(&models.BotUser{KahlaUserId: id})
	return err
}

// 通过ID查找Bot用户
func GetBotUserById(id int64) (*models.BotUser, error) {
	user := &models.BotUser{}

	has, err := orm.X.Id(id).Get(user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// 通过Token查找Bot用户
func GetBotUserByToken(token string) (*models.BotUser, error) {
	user := &models.BotUser{
		Token: token,
	}

	has, err := orm.X.Get(user)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// 获取所有的BotUser
func GetAllBotUser() (as []*models.BotUser, err error) {
	err = orm.X.Find(&as)
	return as, err
}
