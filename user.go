package gowebdav

import (
	"sync"
	"syscall"
)

type userfunc interface {
	// Change the username of the current user
	ChangeName(username string)
	// Change the current user's password
	ChangePwd(password string)
	// Change the permissions of the current user
	ChangeMode(mode mode)
	// Change all the information of the current user
	ReSetInfo(username, password string, mode mode)
	// Check if the password is correct
	ComparePassword(password string) bool

	Name() string
	Mode() mode
}

type user struct {
	name, password string
	mode           mode

	lock *sync.RWMutex
}

type mode int

const (
	O_RDWR   mode = syscall.O_RDONLY
	O_WRONLY mode = syscall.O_WRONLY
	O_RDONLY mode = syscall.O_RDWR
)

func (u *user) ComparePassword(password string) bool {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.password == password
}

func (u *user) Name() string {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.name
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

func (u *user) ReSetInfo(username, password string, mode mode) {
	u.lock.Lock()
	u.name = username
	u.password = password
	u.mode = mode
	u.lock.Unlock()
}
