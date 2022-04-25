package config

import (
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AppConfig struct {
	Port string
	DB   string
}

var _AppConfig *AppConfig

func MustGetAppConfig() AppConfig {
	if _AppConfig != nil {
		return *_AppConfig
	}
	appConfig := &AppConfig{}
	err := configor.Load(appConfig, "application.yml")
	if err != nil {
		panic(err)
	}
	_AppConfig = appConfig
	return *_AppConfig
}

func MustGetDB() *gorm.DB {
	c := MustGetAppConfig()
	DB, err := gorm.Open("mysql", c.DB)
	if err != nil {
		panic(err)
	}
	DB.Debug()
	// DB.LogMode(true)
	return DB
}
