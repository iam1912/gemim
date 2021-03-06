package main

import (
	"log"
	"net/http"

	"github.com/iam1912/gemseries/gemim/config"
	"github.com/iam1912/gemseries/gemim/model"
	"github.com/iam1912/gemseries/gemim/server"
)

func main() {
	c := config.MustGetAppConfig()
	db := config.MustGetDB()
	db.AutoMigrate(&model.User{}, &model.Message{})
	defer db.Close()

	server.RegisterHandles(db)

	log.Println("localhost:8022")
	http.ListenAndServe(c.Port, nil)
}
