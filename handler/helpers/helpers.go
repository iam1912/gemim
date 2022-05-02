package helpers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/iam1912/gemseries/gemim/model"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func String(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}

func RenderSuccessJSON(w http.ResponseWriter, data interface{}) {
	result, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func RenderFailureJSON(w http.ResponseWriter, data interface{}) {
	result, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func WirteAndClose(msg *model.Message, c *websocket.Conn, ctx context.Context, closeMsg string) {
	wsjson.Write(ctx, c, msg)
	c.Close(websocket.StatusUnsupportedData, closeMsg)
}
