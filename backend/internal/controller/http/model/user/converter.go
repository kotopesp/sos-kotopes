package user

import "gitflic.ru/spbu-se/sos-kotopes/internal/core"

func (u *User) ToCoreUser() *core.User {
	if u == nil {
		return nil
	}
	return &core.User{
		Username: u.Username,
		Password: u.Password,
	}
}
