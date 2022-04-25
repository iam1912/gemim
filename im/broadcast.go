package im

import (
	"log"

	"github.com/iam1912/gemim/model"
)

type broadcast struct {
	users                 map[string]*UserConn
	enteringChannel       chan *UserConn
	leavingChannel        chan *UserConn
	MessageChannel        chan *model.Message
	checkUserChannel      chan string
	checkUserCanInChannel chan bool
	usersChannel          chan []*model.User
	requestUsersChannel   chan struct{}
}

var GlobalBroadcast = &broadcast{
	users:                 make(map[string]*UserConn),
	enteringChannel:       make(chan *UserConn),
	leavingChannel:        make(chan *UserConn),
	MessageChannel:        make(chan *model.Message, 1024),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
	usersChannel:          make(chan []*model.User),
	requestUsersChannel:   make(chan struct{}),
}

func (b *broadcast) Run() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.user.NickName] = user
		case user := <-b.leavingChannel:
			delete(b.users, user.user.NickName)
			user.CloseMessageChannel()
		case msg := <-b.MessageChannel:
			for _, user := range b.users {
				if user.user.ID == msg.User.ID {
					continue
				}
				user.MessageChannel <- msg
			}
		case name := <-b.checkUserChannel:
			if _, ok := b.users[name]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*model.User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user.user)
			}
		}
	}
}

func (b *broadcast) UserEntering(u *UserConn) {
	b.enteringChannel <- u
}

func (b *broadcast) UserLeaving(u *UserConn) {
	b.leavingChannel <- u
}

func (b *broadcast) Broadcast(msg *model.Message) {
	if len(b.MessageChannel) >= 1024 {
		log.Println("broadcast queue 满了")
	}
	b.MessageChannel <- msg
}

func (b *broadcast) CanEnterRoom(name string) bool {
	b.checkUserChannel <- name
	return <-b.checkUserCanInChannel
}

func (b *broadcast) GetUserList() []*model.User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
