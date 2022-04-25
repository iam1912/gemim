package ws

import (
	"log"
	"net/http"

	"github.com/iam1912/gemim/handler/helpers"
	"github.com/iam1912/gemim/im"
	"github.com/iam1912/gemim/model"
	"github.com/jinzhu/gorm"
	"nhooyr.io/websocket"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		DB: db,
	}
}

func (h Handler) Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := r.FormValue("token")
	nickname := r.FormValue("nickname")
	ctx := r.Context()

	if l := len(nickname); l < 2 || l > 20 {
		// log.Println("nickname illegal: ", nickname)
		helpers.WirteAndClose(model.NewErrorMessage("非法昵称"), conn, ctx, "nickname illegal")
		return
	}

	user := model.FindAndCreateUser(h.DB, token, nickname, r.RemoteAddr)
	if !im.GlobalBroadcast.CanEnterRoom(user.NickName) {
		// log.Println("用户已经存在：", nickname)
		helpers.WirteAndClose(model.NewErrorMessage("当前用户已存在"), conn, ctx, "user is exist")
		return
	}

	userConn := im.NewUserConn(conn, user)
	go userConn.Write(ctx)

	userConn.MessageChannel <- model.NewWelcomeMessage(user)
	msg := model.NewUserEnterMessage(user)
	im.GlobalBroadcast.Broadcast(msg)

	im.GlobalBroadcast.UserEntering(userConn)
	// log.Println("user:", nickname, "joins chat")

	err = userConn.Read(ctx, user)

	im.GlobalBroadcast.UserLeaving(userConn)
	msg = model.NewUserLeaveMessage(user)
	im.GlobalBroadcast.Broadcast(msg)
	// log.Println("user:", nickname, "leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
