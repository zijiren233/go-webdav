package gowebdav

import (
	"sync"
)

type userfunc interface {
	ChangeName(username string)
	ChangeMode(mode mode)
	ChangePwd(password string)
	SetInfo(username, password string, mode mode)
	Mode() mode

	comparePassword(password string) bool
}

type user struct {
	name, password string
	mode           mode

	lock *sync.RWMutex
}

type mode = uint

const (
	O_RDWR mode = iota
	O_READONLY
)

func (u *user) comparePassword(password string) bool {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.password == password
}

func (u *user) Mode() mode {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.mode
}

func (u *user) ChangeName(username string) {
	u.lock.Lock()
	u.name = username
	u.lock.Unlock()
}

func (u *user) ChangeMode(mode mode) {
	u.lock.Lock()
	u.mode = mode
	u.lock.Unlock()
}

func (u *user) ChangePwd(password string) {
	u.lock.Lock()
	u.password = password
	u.lock.Unlock()
}

func (u *user) SetInfo(username, password string, mode mode) {
	u.lock.Lock()
	u.name = username
	u.password = password
	u.mode = mode
	u.lock.Unlock()
}
