package ws

import (
	"log"
	"net/http"
	"net/url"

	"github.com/iam1912/gemseries/gemim/handler/helpers"
	"github.com/iam1912/gemseries/gemim/im"
	"github.com/iam1912/gemseries/gemim/model"
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
	token := url.QueryEscape(r.FormValue("token"))
	token, _ = url.PathUnescape(token)
	nickname := r.FormValue("nickname")
	ctx := r.Context()

	if l := len(nickname); l < 2 || l > 20 {
		helpers.WirteAndClose(model.NewErrorMessage("非法昵称"), conn, ctx, "nickname illegal")
		return
	}

	user := model.FindOrCreateUser(h.DB, token, nickname, r.RemoteAddr)
	if !im.GlobalBroadcast.CanEnterRoom(user.ID) {
		helpers.WirteAndClose(model.NewErrorMessage("当前用户已存在"), conn, ctx, "user is exist")
		return
	}

	userConn := im.NewUserConn(conn, user.ID)
	go userConn.Write(ctx)

	im.GlobalBroadcast.NotificationEntry(user, userConn)
	user.UpdateOnline(h.DB, true)

	err = userConn.Read(ctx, user)

	im.GlobalBroadcast.NotificationLeaving(user, userConn)
	user.UpdateOnline(h.DB, false)

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
