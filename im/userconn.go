package im

import (
	"context"
	"errors"
	"io"

	"github.com/iam1912/gemseries/gemim/model"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type UserConn struct {
	ID             int
	MessageChannel chan *model.Message
	conn           *websocket.Conn
}

func NewUserConn(c *websocket.Conn, id int) *UserConn {
	return &UserConn{
		ID:             id,
		MessageChannel: make(chan *model.Message, 32),
		conn:           c,
	}
}

func (c *UserConn) Write(ctx context.Context) {
	for msg := range c.MessageChannel {
		wsjson.Write(ctx, c.conn, msg)
	}
}

func (c *UserConn) Read(ctx context.Context, user *model.User) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, c.conn, &receiveMsg)
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		sendMsg := model.NewMessage(user, receiveMsg["content"], receiveMsg["send_time"])
		sendMsg.SplicingAt()
		GlobalBroadcast.Broadcast(sendMsg)
	}
}

func (c *UserConn) CloseMessageChannel() {
	close(c.MessageChannel)
}
