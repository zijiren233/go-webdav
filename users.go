package gowebdav

import "sync"

type usersfunc interface {
	// Get the number of users
	UserNum() int
	// Find a user by username, and return false if the lookup fails
	FindUser(username string) (userfunc, bool)
	// Add a new user, if the user exists, return the existing user
	AddUser(username, password string, mode mode) userfunc
	// Delete users
	DelUser(username string)
}

type users struct {
	usermap map[string]*user

	lock *sync.RWMutex
}

func newusers() usersfunc {
	return &users{usermap: make(map[string]*user), lock: &sync.RWMutex{}}
}

func (u *users) UserNum() int {
	u.lock.RLock()
	num := len(u.usermap)
	u.lock.RUnlock()
	return num
}

func (u *users) FindUser(username string) (userfunc, bool) {
	u.lock.RLock()
	v, ok := u.usermap[username]
	u.lock.RUnlock()
	return v, ok
}

func (u *users) AddUser(username, password string, mode mode) userfunc {
	u.lock.RLock()
	v, ok := u.usermap[username]
	u.lock.RUnlock()
	if ok {
		return v
	}
	newuser := user{name: username, password: password, mode: mode, lock: &sync.RWMutex{}}
	u.lock.Lock()
	u.usermap[username] = &newuser
	u.lock.Unlock()
	return &newuser
}

func (u *users) DelUser(username string) {
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
