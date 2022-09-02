package gowebdav

import (
	"sync"
	"syscall"
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

type mode = int

const (
	O_RDWR   mode = syscall.O_RDONLY
	O_WRONLY mode = syscall.O_WRONLY
	O_RDONLY mode = syscall.O_RDWR
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
