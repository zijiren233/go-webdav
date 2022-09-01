package gowebdav

import (
	"sync"
)

type userfunc interface {
	AddUser(username, password string, mode int) *user
	DelUser(username, password string, mode int)
	ChangeUserMode(username string, mode int)
	ChangeUserPwd(username, password string)
	SetUserRights(username, password string, mode int)
}

type users struct {
	usermap map[string]*user

	lock *sync.RWMutex
}

type user struct {
	name, password string
	mode           int

	lock *sync.RWMutex
}

const (
	O_RDWR = iota
	O_READONLY
)

func (u *users) AddUser(username, password string, mode int) *user {
	u.lock.RLock()
	_, ok := u.usermap[username]
	u.lock.RUnlock()
	if ok {
		return nil
	}
	newuser := user{name: username, password: password, mode: mode, lock: &sync.RWMutex{}}
	u.lock.Lock()
	u.usermap[username] = &newuser
	u.lock.Unlock()
	return &newuser
}

func (u *users) DelUser(username, password string, mode int) {
	u.lock.RLock()
	_, ok := u.usermap[username]
	u.lock.RUnlock()
	if !ok {
		return
	}
	u.lock.Lock()
	delete(u.usermap, username)
	u.lock.Unlock()
}

func (u *users) ChangeUserMode(username string, mode int) {
	u.lock.RLock()
	v, ok := u.usermap[username]
	u.lock.RUnlock()
	if ok {
		v.lock.Lock()
		v.mode = mode
		v.lock.RUnlock()
	}
}

func (u *users) ChangeUserPwd(username, password string) {
	u.lock.RLock()
	v, ok := u.usermap[username]
	u.lock.RUnlock()
	if ok {
		v.lock.Lock()
		v.password = password
		v.lock.RUnlock()
	}
}

func (u *users) SetUserRights(username, password string, mode int) {
	u.lock.RLock()
	v, ok := u.usermap[username]
	u.lock.RUnlock()
	if ok {
		v.lock.Lock()
		v.password = password
		v.mode = mode
		v.lock.RUnlock()
	}
}
