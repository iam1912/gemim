package im

import (
	"container/ring"

	"github.com/iam1912/gemim/model"
)

type offlineProcessor struct {
	n          int
	recentRing *ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	return &offlineProcessor{
		n:          10,
		recentRing: ring.New(10),
	}
}

func (o *offlineProcessor) Save(msg *model.Message) {
	if msg.Type != model.MsgTypeNormal {
		return
	}
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()
}

func (o *offlineProcessor) Send(user *UserConn) {
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*model.Message)
		}
	})
}
