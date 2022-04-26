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
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"enter_at"`
	UpdatedAt time.Time
}

var System = &User{}

func FindOrCreateUser(db *gorm.DB, token, nickname, addr string) *User {
	id := 0
	user := &User{}
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

func (user *User) UpdateOnline(db *gorm.DB, isOnline bool) {
	db.Model(&User{}).Where("id = ?", user.ID).Update("Is_online", isOnline)
}
