package main

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/xconnio/wampproto-go/util"
	"github.com/xconnio/xconn-go"
	"gorm.io/gorm"
)

type UserManager struct {
	onlineUserByID map[int]*User

	sync.Mutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		onlineUserByID: make(map[int]*User),
	}
}

func (m *UserManager) statusOnlineHandler(db *gorm.DB) func(event *xconn.Event) {
	return func(event *xconn.Event) {
		userID, ok := util.AsInt64(event.Arguments[0])
		if !ok {
			fmt.Println("first argument must be a valid int")
			return
		}

		user, err := GetUserByID(db, uint(userID))
		if err != nil {
			fmt.Println(err)
			return
		}

		m.Lock()
		m.onlineUserByID[int(userID)] = &user
		m.Unlock()
	}
}

func (m *UserManager) statusOfflineHandler(event *xconn.Event) {
	userID, ok := util.AsInt64(event.Arguments[0])
	if !ok {
		fmt.Println("first argument must be a valid int")
		return
	}

	m.Lock()
	delete(m.onlineUserByID, int(userID))
	m.Unlock()
}

func (m *UserManager) onlineUser(_ context.Context, _ *xconn.Invocation) *xconn.Result {
	users := slices.Collect(maps.Values(m.onlineUserByID))
	return &xconn.Result{Arguments: []any{users}}
}
