package model

import (
	"time"

	"github.com/iam1912/gemim/utils"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID        int       `json:"id"`
	Token     string    `json:"token" gorm:"unique"`
	NickName  string    `json:"nickname"`
	Addr      string    `json:"addr"`
	CreatedAt time.Time `json:"enter_at"`
}

var System = &User{}

func FindAndCreateUser(db *gorm.DB, token, nickname, addr string) *User {
	var (
		user = &User{}
		id   int
	)
	if token != "" {
		if err := db.Where("token = ? AND nick_name = ?", token, nickname).First(&user).Error; err == nil {
			return user
		}
	}
	db.Model(&User{}).Count(&id)
	user = &User{
		ID:       id + 1,
		Token:    utils.GenToken(id, nickname, "123123"),
		NickName: nickname,
		Addr:     addr,
	}
	db.Save(&user)
	return user
}
