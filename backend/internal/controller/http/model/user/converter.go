package user

import "github.com/kotopesp/sos-kotopes/internal/core"

func (u *User) ToCoreUser() core.User {
	return core.User{
		Username:     u.Username,
		PasswordHash: u.Password,
		Description:  u.Description,
		Photo:        u.Photo,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
	}
}
