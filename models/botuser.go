package models

type BotUser struct {
	Id int64
	Token string `xorm:"unique"`
	Nickname string
	KahlaUserId string `xorm:"unique"`
	Version int `xorm:"version"`
}