package home

import (
	"net/http"
	"os"
	"text/template"

	"github.com/iam1912/gemim/handler/helpers"
	"github.com/iam1912/gemim/im"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		DB: db,
	}
}

func (h Handler) Home(w http.ResponseWriter, req *http.Request) {
	path, _ := os.Getwd()
	t, err := template.ParseFiles(path + "/template/home.html")
	if err != nil {
		helpers.String(w, "模板解析错误")
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		helpers.String(w, "模板解析错误")
		return
	}
}

func (h Handler) UserList(w http.ResponseWriter, req *http.Request) {
	userList := im.GlobalBroadcast.GetUserList()
	if len(userList) == 0 {
		helpers.RenderFailureJSON(w, "[]")
	}
	helpers.RenderSuccessJSON(w, userList)
}
