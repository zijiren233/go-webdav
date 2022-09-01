package gowebdav

type user struct {
	name, password string
	mode           int
}

const (
	O_RDWR = iota
	O_READONLY
)

func (client *client) AddUser(username, password string, mode int) Client {
	client.userInfo[username] = &user{name: username, password: password, mode: mode}
	return client
}

func (client *client) ChangeUserMode(username string, mode int) Client {
	if v, ok := client.userInfo[username]; ok {
		v.mode = mode
	}
	return client
}

func (client *client) ChangeUserPwd(username, password string) Client {
	if v, ok := client.userInfo[username]; ok {
		v.password = password
	}
	return client
}

func (client *client) SetUserRights(username, password string, mode int) Client {
	if v, ok := client.userInfo[username]; ok {
		v.password = password
		v.mode = mode
	}
	return client
}
