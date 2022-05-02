package im

import "github.com/iam1912/gemseries/gemim/model"

type broadcast struct {
	users                 map[int]*UserConn
	enteringChannel       chan *UserConn
	leavingChannel        chan *UserConn
	messageChannel        chan *model.Message
	checkUserChannel      chan int
	checkUserCanInChannel chan bool
}

var GlobalBroadcast = &broadcast{
	users:                 make(map[int]*UserConn),
	enteringChannel:       make(chan *UserConn),
	leavingChannel:        make(chan *UserConn),
	messageChannel:        make(chan *model.Message),
	checkUserChannel:      make(chan int),
	checkUserCanInChannel: make(chan bool),
}

func (b *broadcast) Run() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.ID] = user
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			delete(b.users, user.ID)
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.ID == msg.User.ID {
					continue
				}
				user.MessageChannel <- msg
			}
			if msg.Type == model.MsgTypeNormal {
				OfflineProcessor.Save(msg)
			}
		case name := <-b.checkUserChannel:
			if _, ok := b.users[name]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}
	}
}

func (b *broadcast) UserEntering(user *UserConn) {
	b.enteringChannel <- user
}

func (b *broadcast) UserLeaving(user *UserConn) {
	b.leavingChannel <- user
}

func (b *broadcast) Broadcast(msg *model.Message) {
	b.messageChannel <- msg
}

func (b *broadcast) CanEnterRoom(id int) bool {
	b.checkUserChannel <- id
	return <-b.checkUserCanInChannel
}

func (b *broadcast) NotificationEntry(user *model.User, userConn *UserConn) {
	userConn.MessageChannel <- model.NewWelcomeMessage(user)
	msg := model.NewUserEnterMessage(user)
	GlobalBroadcast.Broadcast(msg)
	GlobalBroadcast.UserEntering(userConn)
}

func (b *broadcast) NotificationLeaving(user *model.User, userConn *UserConn) {
	GlobalBroadcast.UserLeaving(userConn)
	msg := model.NewUserLeaveMessage(user)
	GlobalBroadcast.Broadcast(msg)
}
