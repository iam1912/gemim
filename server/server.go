package server

import (
	"net/http"

	"github.com/iam1912/gemim/handler/home"
	"github.com/iam1912/gemim/handler/ws"
	"github.com/iam1912/gemim/im"
	"github.com/jinzhu/gorm"
)

func RegisterHandles(db *gorm.DB) {
	go im.GlobalBroadcast.Run()

	homeHandler := home.NewHandler(db)
	wsHandler := ws.NewHandler(db)

	http.HandleFunc("/", homeHandler.Home)
	http.HandleFunc("/user_list", homeHandler.UserList)
	http.HandleFunc("/ws", wsHandler.Ws)
}
